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
}

func NewServer(db *store.Store, router chi.Router, log *zap.Logger, notifier notify.Notifier) *Server {
	server := &Server{
		db:       db,
		context:  context.Background(),
		log:      log,
		notifier: notifier,
	}

	router.HandleFunc("/login", server.loginHandler)
	router.HandleFunc("/refresh", server.refreshHandler)

	server.log.Info("router configured")
	return server
}

func (s Server) loginHandler(writer http.ResponseWriter, request *http.Request) {
	s.log.Info("serving /login")
	srv := auth.NewService("my_secret", s.db, s.notifier, pkg.NormalClock{})
	userID := request.FormValue("user_id")
	secret := request.FormValue("secret")
	ip := request.RemoteAddr
	if userID == "" {
		http.Error(writer, "empty user", http.StatusUnauthorized)
		s.log.Error("empty user id")
		return
	}
	var authTokens auth.AuthTokens
	authTokens, err := srv.Authorize(s.context, secret, userID, ip)
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
	srv := auth.NewService("my_secret", s.db, s.notifier, pkg.NormalClock{})

	refresh := request.FormValue("refresh_token")
	userID := request.FormValue("user_id")
	newAccess, err := srv.Refresh(s.context, refresh, request.RemoteAddr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to refresh access token", zap.Error(err))
		return
	}
	if err := json.NewEncoder(writer).Encode(newAccess); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		s.log.Error("failed to encode auth tokens", zap.String("user_id", userID))
		return
	}

	s.log.Info("/refresh serving done")
}
