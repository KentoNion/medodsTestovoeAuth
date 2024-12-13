package server

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"medodsTestovoe/auth"
	"medodsTestovoe/auth/pkg"
	notify "medodsTestovoe/gates/notifier"
	store "medodsTestovoe/gates/postgres"
	"net/http"
)

type Server struct {
	db       *store.Store
	context  context.Context
	log      *zap.Logger
	notifier notify.Notifier
	srv      *auth.Service
}

func NewServer(db *store.Store, router chi.Router, log *zap.Logger, notifier notify.Notifier) *Server {
	server := &Server{
		db:       db,
		context:  context.Background(),
		log:      log,
		notifier: notifier,
		srv:      auth.NewService("my_secret", db, notifier, pkg.NormalClock{}),
	}

	router.Method(http.MethodPost, "/login", http.HandlerFunc(server.loginHandler))
	router.Method(http.MethodPost, "/refresh", http.HandlerFunc(server.refreshHandler))

	server.log.Info("router configured")
	return server
}

func (s Server) loginHandler(writer http.ResponseWriter, request *http.Request) {
	s.log.Info("serving /login")
	userID := request.FormValue("user_id")
	secret := request.FormValue("secret")
	ip := request.RemoteAddr
	if userID == "" {
		http.Error(writer, "empty user", http.StatusUnauthorized)
		s.log.Error("empty user id")
		return
	}
	var authTokens auth.AuthTokens
	authTokens, err := s.srv.Authorize(s.context, secret, userID, ip)
	if err == auth.ErrWrongToken {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		zap.Error(err)
		return
	}

	if err := json.NewEncoder(writer).Encode(authTokens); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to encode auth tokens", zap.String("user_id", userID))
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
	newTokens, err := s.srv.Refresh(s.context, refresh, request.RemoteAddr)
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
