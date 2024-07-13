package main

import (
	"JustTesting/internal/config"
	"JustTesting/internal/http-server/handlers/alias/get"
	"JustTesting/internal/http-server/handlers/delete"
	"JustTesting/internal/http-server/handlers/deleteRange"
	"JustTesting/internal/http-server/handlers/redirect"
	"JustTesting/internal/http-server/handlers/url/save"
	"JustTesting/internal/http-server/middleware/logger"
	"JustTesting/internal/lib/logger/sl"
	"JustTesting/internal/storage/postgreSql"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting server", slog.String("env", cfg.Env))
	log.Debug("Debug logging enabled")

	//connString := "postgres://postgres:angelo4ek@localhost:5432/test1?sslmode=disable"
	connStr := "host=localhost port=5432 user=postgres password=angelo4ek dbname=test1 sslmode=disable"
	storage, err := postgreSql.NewPG(connStr)
	if err != nil {
		fmt.Errorf("Failed to initialize storage: %v", sl.Err(err))
		os.Exit(1)
	}
	fmt.Println(storage)

	router := chi.NewRouter()
	//router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Route("/url", func(r chi.Router) {
		//r.Use(middleware.BasicAuth("url-shortener", map[string]string{
		//	cfg.HTTPServer.User: cfg.HTTPServer.Password,
		//}))
		r.Post("/", save.New(log, storage))
		r.Delete("/deleteRange", deleteRange.New(log, storage))
		r.Delete("/delete", delete.New(log, storage))
	})
	//router.Use(middleware.Recoverer)
	//router.Use(middleware.URLFormat)

	router.Get("/getalias", get.New(log, storage))

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeOut,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
	}
	log.Error("failed to start server", sl.Err(srv.Shutdown(context.Background())))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:

		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:

		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
