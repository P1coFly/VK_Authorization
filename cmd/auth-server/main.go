package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/P1coFly/VK_Authorization/docs"
	"github.com/P1coFly/VK_Authorization/internal/config"
	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/authorize"
	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/feed"
	"github.com/P1coFly/VK_Authorization/internal/http-server/handlers/register"
	"github.com/P1coFly/VK_Authorization/internal/storage/postgresql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// При изменение анотации swagger, перед запуском контейнера, надо сгенерировать документацию
// swag init -g .\cmd\auth-server\main.go

// @title Authorizathion App API
// @version 1.0
// @description API Server for Registration and Authorizathion

// @host localhost:cfg.Port
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Need a Bearer Token, like: Bearer `<`my_token`>`
func main() {
	// читаем конфиг
	cfg := config.MustLoad()
	// инициализируем логер
	log := setupLogger(cfg.Env)

	log.Info("starting api-servies", "env", cfg.Env)
	log.Debug("cfg data", "data", cfg)

	// читаем ключ для генирации jwt
	keyJWT := os.Getenv("KEY_JWT")
	if keyJWT == "" {
		log.Error("CONFIG_PATH is not set")
		os.Exit(1)
	}

	// инициализируем storage
	storage, err := postgresql.New(cfg.Host_db, cfg.Port_db, cfg.User_db, cfg.Password_db, cfg.Name_db)
	if err != nil {
		log.Error("failed to connect storage", "error", err)
		os.Exit(1)
	}
	log.Info("connect to db is successful", "host", cfg.Host_db)

	// инициализируем router
	router := chi.NewRouter()

	// добавляем endpoints
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/register", register.New(log, storage))
	router.Post("/authorize", authorize.New(log, storage, []byte(keyJWT)))

	router.Get("/feed", feed.New(log, []byte(keyJWT)))

	//Для доступа к swagger надо пройти по URI /swagger/
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //По URI /swagger/doc.json будет ледать спецификация в формате JSON
	))

	log.Info("starting server", slog.String("port", cfg.Port))

	// инициализируем server и запускаем
	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", "error", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
