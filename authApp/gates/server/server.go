package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"medodsTestovoe/auth"
	"medodsTestovoe/auth/pkg"
	"net/http"
)

// Интерфейсы здесь необходимы для реализации моков в тестах
type db interface {
	Save(ctx context.Context, token pkg.Refresh, userID string, ip string) error
	Get(ctx context.Context, userID string, token pkg.Refresh) (bool, string, error)
	Delete(ctx context.Context, userID string) error
	CheckUserExist(ctx context.Context, userID string) (bool, error)
}
type notifier interface {
	NotifyNewLogin(ctx context.Context, userID string, oldIP string, newIP string) error
}

type Server struct {
	db       db
	context  context.Context
	log      *zap.Logger
	notifier notifier
	srv      *auth.Service
}

func NewServer(db db, router chi.Router, log *zap.Logger, notifier notifier) *Server {
	server := &Server{ //формируем структуру сервера
		db:       db,
		context:  context.Background(),
		log:      log,
		notifier: notifier,
		srv:      auth.NewService("my_secret", db, notifier, pkg.NormalClock{}),
	}
	//роутим эндпоинты
	router.Method(http.MethodPost, "/login", http.HandlerFunc(server.loginHandler))
	router.Method(http.MethodPost, "/refresh", http.HandlerFunc(server.refreshHandler))
	server.log.Info("router configured")
	return server
}

func (s Server) loginHandler(writer http.ResponseWriter, request *http.Request) {
	s.log.Info("serving /login")
	userID := request.FormValue("user_id")
	secret := request.FormValue("secret")
	if userID == "" { //защита от пустого юзера
		http.Error(writer, "empty user", http.StatusUnauthorized)
		s.log.Error("empty user id")
		return
	}
	if secret == "" { //защита от пустого секрета (см пункт 10 особенностей в readme файле)
		http.Error(writer, "empty secret", http.StatusUnauthorized)
		s.log.Error("empty secret")
		return
	}
	var authTokens auth.AuthTokens                                                    //создаём пустые токены которые будет заполнять данными и отдавать в ответе
	authTokens, err := s.srv.Authorize(s.context, secret, userID, request.RemoteAddr) //здесь основная логика authorize
	if err == auth.ErrWrongToken {
		http.Error(writer, err.Error(), http.StatusUnauthorized) //если токен кривой, делаем статус анотхарайзд
		zap.Error(err)
		return
	}
	if err == auth.ErrUserAlreadyExists {
		http.Error(writer, err.Error(), http.StatusConflict)
		s.log.Error("user already exists") //если такой юзер уже есть, тогда статус конфликт
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError) //если что-то другое пошло не так, то пишем что что-то на сервере пошло не так
		zap.Error(err)
		return
	}

	if err := json.NewEncoder(writer).Encode(authTokens); err != nil { //если всё ок то формируем тело ответа в Джейсона
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to encode auth tokens", zap.String("user_id", userID)) //ошибка если не получилось сформировать тело ответа
		return
	}
	s.log.Info("/login serving done")
	return
}

func (s Server) refreshHandler(writer http.ResponseWriter, request *http.Request) {
	s.log.Info("serving /refresh")

	refreshStr := request.FormValue("refresh_token")
	refresh := pkg.Refresh(refreshStr)
	userID := request.FormValue("user_id")
	newTokens, err := s.srv.Refresh(s.context, userID, refresh, request.RemoteAddr)
	if userID == "" {
		http.Error(writer, "empty user", http.StatusUnauthorized)
		s.log.Error("empty user id")
		return
	}
	if refresh == "" {
		http.Error(writer, "empty refresh token", http.StatusUnauthorized)
		s.log.Error("empty refresh token")
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to refresh access token", zap.Error(err))
		return
	}
	if err := json.NewEncoder(writer).Encode(newTokens); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to encode auth tokens", zap.String("user_id", userID))
		return
	}

	s.log.Info("/refresh serving done")
}
