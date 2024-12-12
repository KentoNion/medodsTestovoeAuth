package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	notify "medodsTestovoe/gates/notifier"
	store "medodsTestovoe/gates/postgres"
	"medodsTestovoe/gates/server"
	"medodsTestovoe/internal/config"
	"net/http"
)

func main() {
	//считываем файл конфига
	Cfg, err := config.MustLoad() // ипортируем конфиг
	if err != nil {
		panic(err)
	}

	//регестрируем логгер
	log, err := zap.NewDevelopment() //регестрируем логгер
	if err != nil {
		panic(err)
	}

	//регестрируем нотифаер
	notifier := notify.InitNotifier()

	//инициализируем бд
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s%s sslmode=%s", Cfg.DB.DbUser, Cfg.DB.DbPassword, Cfg.DB.DbName, Cfg.DB.DbHost, Cfg.DB.DbPort, Cfg.DB.DbSSLMode)
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db := store.NewDB(conn)

	//накатываем миграцию
	err = goose.Up(conn.DB, "./gates\\postgres\\migrations")
	if err != nil {
		panic(err)
	}

	//регестрируем роутер
	router := chi.NewRouter()

	//запускаем сервер
	_ = server.NewServer(db, router, log, notifier)
	err = http.ListenAndServe(Cfg.Server.ServerHost+":"+Cfg.Server.ServerPort, router)
	if err != nil {
		log.Error("server error", zap.Error(err))
	}
}
