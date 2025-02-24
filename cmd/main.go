package main

import (
	"log/slog"
	"net/http"

	"db_practice/config"
	"db_practice/internal/database"
	"db_practice/internal/handler"
	"db_practice/internal/models"
	"db_practice/internal/repository"
	"db_practice/internal/services"
)

func main() {
	cfg := config.GetConfig("config.yaml")

	db := database.ConnectDB(cfg.DB.Connection)
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		slog.Error("Failed to migrate database: ", slog.Any("error", err))
	}

	orderRepo := repository.NewOrderRepository(db)

	orderChannel := make(chan models.Order)
	go func() {
		if err := services.ParseOrdersFromFile(cfg.File.Path, orderChannel); err != nil {
			slog.Error("Failed to parse file: ", slog.Any("error", err))
		}
		close(orderChannel)
	}()

	go func() {
		for order := range orderChannel {
			if err := orderRepo.SaveOrder(&order); err != nil {
				slog.Error("Failed to save order: ", slog.Any("error", err))
			}
		}
	}()
	service := services.NewService(orderRepo)
	httpServer := handler.NewHTTPServer(service)
	router := httpServer.Routes()

	slog.Info("Starting server on ", slog.String("port ", cfg.Server.Port))
	err := http.ListenAndServe(cfg.Server.Port, router) // was (fmt.Sprintf(":%s", cfg.Server.Port), router)
	if err != nil {
		slog.Error("Can't start service:", slog.Any("error", err))
	}
}
