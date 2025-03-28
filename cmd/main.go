package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"db_practice/config"
	"db_practice/internal/cache"
	"db_practice/internal/database"
	"db_practice/internal/handler"
	"db_practice/internal/models"
	"db_practice/internal/repository"
	"db_practice/internal/services"
	"db_practice/internal/websocket"
)

func main() {
	cfg := config.GetConfig("config.yaml")

	db := database.ConnectDB(cfg.DB.Connection)
	defer db.Close()

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

	cache := cache.NewCache(orderRepo)
	service := services.NewService(cache)

	httpServer := handler.NewHTTPServer(service)
	router := httpServer.Routes()

	wsServer := websocket.NewWSServer(service)
	wsRouter := wsServer.WSRoute()

	slog.Info("Starting server on", slog.String("port", cfg.Server.Port))

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), router)
		if err != nil {
			slog.Error("Can't start service:", slog.Any("error", err))
		}
	}()

	slog.Info("Starting server on", slog.String("port", cfg.Server.Ws))
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Ws), wsRouter)
		if err != nil {
			slog.Error("Can't start service:", slog.Any("error", err))
		}
	}()

	<-ctx.Done()
	slog.Info("Got signal, exit program")
}
