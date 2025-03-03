package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

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

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt) // graceful shutdown
	defer cancelFunc()

	orderRepo := repository.NewOrderRepository(db)
	orderChannel := make(chan models.Order)

	go func() {
		if err := services.ParseOrdersFromFile(ctx, cfg.File.Path, orderChannel); err != nil {
			slog.Error("Failed to parse file: ", slog.Any("error", err))
		}
		close(orderChannel)
	}()

	go func() {
		for order := range orderChannel {
			if err := orderRepo.SaveOrder(ctx, &order); err != nil {
				slog.Error("Failed to save order: ", slog.Any("error", err))
			}
		}
	}()
	service := services.NewService(orderRepo)
	httpServer := handler.NewHTTPServer(service)
	router := httpServer.Routes()

	slog.Info("Starting server on ", slog.String("port ", cfg.Server.Port))
	go func() {
		err := http.ListenAndServe(cfg.Server.Port, router) // was (fmt.Sprintf(":%s", cfg.Server.Port), router)
		if err != nil {
			slog.Error("Can't start service:", slog.Any("error", err))
		}
	}()
	<-ctx.Done()
	slog.Info("Got signal, exit program")
}
