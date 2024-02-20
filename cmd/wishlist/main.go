package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"
	"wish_list/internal/config"
	"wish_list/internal/http-server/handlers/sharelist"
	"wish_list/internal/http-server/handlers/wishlist"
	"wish_list/internal/http-server/handlers/wishlist/item"
	"wish_list/internal/http-server/middleware/logger"
	"wish_list/internal/storage/postgres"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)

	log.Info("App started", slog.String("env", cfg.Env))
	log.Debug("Debugging started")

	storage, err := postgres.New(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Error("failed to init storage: ", err)
		os.Exit(1)
	}

	log.Info("storage successfully initialized")

	router := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(corsHandler.Handler)

	router.Post("/api/wishlist/create", wishlist.Create(log, storage))             //создание вишлиста в личном кабинете
	router.Post("/api/items/add", item.Create(log, storage))                       // добавление подарка в вишлист из ЛК
	router.Get("/api/sharelist/{alias}", sharelist.GetList(log, storage))          // получение вишлиста по алиасу
	router.Get("/api/wishlist/getforuser", wishlist.GetAllLists(log, storage))     // получение списка вишлистов пользователя в ЛК
	router.Post("/api/wishlist/delete", wishlist.Delete(log, storage))             // удаление конкретного вишлиста в ЛК
	router.Get("/api/wishlist/{wishlistId}/items", item.GetByWishId(log, storage)) // получение списка подарков конкретного вишлиста в ЛК
	router.Post("/api/item/delete", item.Delete(log, storage))                     // удаление подарка из вишлиста в ЛК

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
