package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"db_practice/internal/models"
	"db_practice/internal/services"
)

type HTTPServer struct {
	Service services.Service
}

func NewHTTPServer(service services.Service) *HTTPServer {
	return &HTTPServer{Service: service}
}

func (s *HTTPServer) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.Service.SaveOrder(&order); err != nil {
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	slog.Info(
		"Added order:",
		slog.Int("ShopID=", order.Payment.ShopID),
		slog.String("Address=", order.Payment.Address),
		slog.Float64("TotalAmount=", order.Payment.TotalAmount),
		slog.Int("Items=", len(order.Payment.Items)),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order added successfully"))
}

func (s *HTTPServer) GetOrdersByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	const inputLayout = "2006-01-02T15:04:05.000" // input format
	const dbLayout = "2006-01-02 15:04:05.000"    // DB format

	startTime, err := time.Parse(inputLayout, start)
	if err != nil {
		errors.Wrap(err, "Can't parse input start time format")
		return
	}

	endTime, err := time.Parse(inputLayout, end)
	if err != nil {
		errors.Wrap(err, "Can't parse input end time format")
		return
	}

	slog.Info(
		"Period:",
		slog.Time("Parsed start: ", startTime),
		slog.Time("Parsed end: ", endTime),
	)

	startTimeFormatted, err := time.Parse(dbLayout, startTime.Format(dbLayout))

	if err != nil {
		errors.Wrap(err, "failed to normalize start time")
		return
	}

	endTimeFormatted, err := time.Parse(dbLayout, endTime.Format(dbLayout))
	if err != nil {
		errors.Wrap(err, "failed to normalize end time")
		return
	}

	orders, err := s.Service.GetOrdersByPeriod(startTimeFormatted, endTimeFormatted)
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		slog.Error("Failed to retrieve orders", slog.Any("error", err))
		return
	}
	if len(orders) == 0 {
		slog.Info("No orders found for the specified period")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (s *HTTPServer) GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	shops, err := s.Service.GetShops()
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
}

func (s *HTTPServer) GetRevenueByShopHandler(w http.ResponseWriter, r *http.Request) {
	revenue, err := s.Service.GetRevenueByShop()
	if err != nil {
		http.Error(w, "Failed to retrieve revenue data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenue)
}

func (s *HTTPServer) GetAverageCheckByShopHandler(w http.ResponseWriter, r *http.Request) {
	averageCheck, err := s.Service.GetAverageCheckByShop()
	if err != nil {
		http.Error(w, "Failed to retrieve average check data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(averageCheck)
}

func (s *HTTPServer) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/add-order", s.AddOrderHandler)
	router.Get("/orders-by-period", s.GetOrdersByPeriodHandler)
	router.Get("/shops", s.GetShopsHandler)
	router.Get("/revenue-by-shop", s.GetRevenueByShopHandler)
	router.Get("/average-check-by-shop", s.GetAverageCheckByShopHandler)
	return router
}
