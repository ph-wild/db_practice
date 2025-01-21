package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"db_practice/internal/models"
	"db_practice/internal/repositories"

	"github.com/go-chi/chi/v5"
)

type HTTPServer struct {
	OrderRepo *repositories.OrderRepository
}

func (s *HTTPServer) AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.OrderRepo.SaveOrder(&order); err != nil {
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	// Логирование добавленного заказа в консоль
	fmt.Printf("Added order: ShopID=%d, Address=%s, TotalAmount=%.2f, Items=%d\n",
		order.Payment.ShopID,
		order.Payment.Address,
		order.Payment.TotalAmount,
		len(order.Payment.Items),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order added successfully"))
}

func (s *HTTPServer) GetOrdersByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	orders, err := s.OrderRepo.GetOrdersByPeriod(start, end)
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
func (s *HTTPServer) GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	shops, err := s.OrderRepo.GetShops()
	if err != nil {
		http.Error(w, "Failed to retrieve shops", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shops)
}

func (s *HTTPServer) GetRevenueByShopHandler(w http.ResponseWriter, r *http.Request) {
	revenue, err := s.OrderRepo.GetRevenueByShop()
	if err != nil {
		http.Error(w, "Failed to retrieve revenue data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenue)
}

func (s *HTTPServer) GetAverageCheckByShopHandler(w http.ResponseWriter, r *http.Request) {
	averageCheck, err := s.OrderRepo.GetAverageCheckByShop()
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
