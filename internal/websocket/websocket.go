package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"db_practice/internal/models"
	"db_practice/internal/services"
)

type WSServer struct {
	Service services.ServiceInterface
}

func NewWSServer(service services.ServiceInterface) *WSServer {
	return &WSServer{Service: service}
}

func (s *WSServer) AddOrderHandler(conn *websocket.Conn, data json.RawMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		//errors.Wrap(err, "Error reading AddOrderHandler")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Invalid request payload"})
		return
	}

	if err := s.Service.SaveOrder(ctx, &order); err != nil {
		//errors.Wrap(err, "AddOrderHandler saving error")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Failed to save order"})
		return
	}

	slog.Info(
		"Added order:",
		slog.Int("ShopID=", order.Payment.ShopID),
		slog.String("Address=", order.Payment.Address),
		slog.Float64("TotalAmount=", order.Payment.TotalAmount),
		slog.Int("Items=", len(order.Payment.Items)),
	)
	conn.WriteJSON(map[string]string{"message": "Order added successfully"})
}

func (s *WSServer) GetOrdersByPeriodHandler(conn *websocket.Conn, data json.RawMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	var req struct {
		Start string `json:"start"` // FIXME time.Time
		End   string `json:"end"`   // FIXME time.Time
	}

	if err := json.Unmarshal(data, &req); err != nil {
		//errors.Wrap(err, "GetOrdersByPeriodHandler error data")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Invalid request format"})
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05.000", req.Start)
	if err != nil {
		//errors.Wrap(err, "Cannot parse start time")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Invalid start time format"})
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04:05.000", req.End)
	if err != nil {
		//errors.Wrap(err, "Cannot parse end time")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Invalid end time format"})
		return
	}

	orders, err := s.Service.GetOrdersByPeriod(ctx, startTime, endTime)
	if err != nil {
		//errors.Wrap(err, "GetOrdersByPeriod Error")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Can't get Orders"})
		return
	}
	conn.WriteJSON(orders)
}

func (s *WSServer) GetShopsHandler(conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second) // TODO: seconds to config
	defer cancel()
	shops, err := s.Service.GetShops(ctx)
	if err != nil {
		//errors.Wrap(err, "GetShops failed")
		slog.Error("%+v\n", err)
		conn.WriteJSON(map[string]string{"error": "Failed to retrieve shops"})
		return
	}
	conn.WriteJSON(shops)
}

func (s *WSServer) GetRevenueByShopHandler(conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second) // TODO: to config
	defer cancel()
	revenue, err := s.Service.GetRevenueByShop(ctx)
	if err != nil {
		//errors.Wrap(err, "GetRevenueByShopHandler failed")
		slog.Error("%+v\n", err)
		conn.WriteJSON("Failed to retrieve revenue data")
		return
	}
	conn.WriteJSON(revenue)
}

func (s *WSServer) GetAverageCheckByShopHandler(conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second) // TODO: to config
	defer cancel()
	averageCheck, err := s.Service.GetAverageCheckByShop(ctx)
	if err != nil {
		//errors.Wrap(err, "GetAverageCheckByShop failed")
		slog.Error("%+v\n", err)
		conn.WriteJSON("Failed to retrieve average check data")
		return
	}
	conn.WriteJSON(averageCheck)
}

func (s *WSServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil) // r.Context only 1 for all
	defer conn.Close()

	if err != nil {
		errors.Wrap(err, "WebSocket upgrader error")
		return
	}

	for {
		var request struct {
			Action string          `json:"action"`
			Data   json.RawMessage `json:"data"`
		}

		if err := conn.ReadJSON(&request); err != nil {
			errors.Wrap(err, "WebSocket reading error")
			break
		}

		switch request.Action {
		case "add_order":
			s.AddOrderHandler(conn, request.Data)
		case "get_orders_by_period":
			s.GetOrdersByPeriodHandler(conn, request.Data)
		case "get_shops":
			s.GetShopsHandler(conn)
		case "get_revenue_by_shop":
			s.GetRevenueByShopHandler(conn)
		case "get_average_check_by_shop":
			s.GetAverageCheckByShopHandler(conn)
		default:
			conn.WriteJSON("Unknown action")
		}
	}
}

func (s *WSServer) WSRoute() chi.Router {
	router := chi.NewRouter() // middleware router.Use()
	router.Get("/ws", s.HandleWebSocket)
	return router
}
