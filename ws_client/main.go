package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverAddr := "ws://localhost:8282/ws"

	fmt.Println("🔌 Connecting to WebSocket server...")
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to", serverAddr)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			fmt.Println("Server response:", string(message))
		}
	}()
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("\nAvailable actions:")
		fmt.Println("1. Add order")
		fmt.Println("2. Get orders by period")
		fmt.Println("3. Get shops")
		fmt.Println("4. Get revenue by shop")
		fmt.Println("5. Get average check by shop")
		fmt.Print("Select action (1-5 or 'exit'): ")

		if !reader.Scan() {
			break
		}
		input := reader.Text()

		var request map[string]interface{}
		switch input {
		case "1":
			request = map[string]interface{}{
				"action": "add_order",
				"data": map[string]interface{}{
					"Payment": map[string]interface{}{
						"ShopID":      1,
						"Address":     "Test Address",
						"TotalAmount": 100.5,
						"Date":        "2006-01-02T15:04:05.000",
					},
				},
			}
		case "2":
			request = map[string]interface{}{
				"action": "get_orders_by_period",
				"data": map[string]string{
					"start": time.Now().AddDate(0, -1, 0).Format("2006-01-02T15:04:05.000"),
					"end":   time.Now().Format("2006-01-02T15:04:05.000"),
				},
			}
		case "3":
			request = map[string]interface{}{
				"action": "get_shops",
			}
		case "4":
			request = map[string]interface{}{
				"action": "get_revenue_by_shop",
			}
		case "5":
			request = map[string]interface{}{
				"action": "get_average_check_by_shop",
			}
		case "exit":
			fmt.Println("Closing connection...")
			return
		default:
			fmt.Println("Invalid action")
			continue
		}

		jsonData, err := json.Marshal(request)
		if err != nil {
			log.Println("JSON error:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Println("Send error:", err)
			break
		}
		fmt.Println("Sent:", string(jsonData))
	}

	fmt.Println("Closing WebSocket connection...")
}
