package main

import (
	"db_practice/internal/database"
	"db_practice/internal/models"
	"db_practice/internal/repositories"
	"db_practice/internal/services"
	"log"
	"net/http"

	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	// Загрузка конфигурации
	file, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	defer file.Close()

	var cfg struct {
		DB     struct{ Connection string }
		Server struct{ Port string }
		File   struct{ Path string }
	}
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// Подключение к БД
	db := database.ConnectDB(cfg.DB.Connection)
	defer db.Close()

	// Миграция
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Инициализация репозитория
	orderRepo := &repositories.OrderRepository{DB: db}

	// Парсинг файла
	ch := make(chan models.Order)
	go func() {
		if err := services.ParseOrdersFromFile(cfg.File.Path, ch); err != nil {
			log.Fatalf("Failed to parse file: %v", err)
		}
	}()

	// Сохранение данных из канала
	for order := range ch {
		if err := orderRepo.SaveOrder(&order); err != nil {
			log.Printf("Failed to save order: %v", err)
		}
	}

	// Запуск HTTP-сервера
	httpServer := &services.HTTPServer{OrderRepo: orderRepo}
	router := httpServer.Routes()

	log.Printf("Starting server on port %s...", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}
