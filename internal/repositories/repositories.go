package repositories

import (
	"db_practice/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	DB *sqlx.DB
}

func (r *OrderRepository) SaveOrder(order *models.Order) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	// Сохранение заказа
	var orderID int64
	err = tx.QueryRowx(`INSERT INTO orders (shop_id, address, date, total_amount)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		order.Payment.ShopID, order.Payment.Address, order.Payment.Date, order.Payment.TotalAmount).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение позиций заказа
	for _, item := range order.Payment.Items {
		_, err := tx.Exec(`INSERT INTO items (name, price) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`,
			item.Name, item.Price)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(`INSERT INTO order_items (order_id, item_id, quantity, total_price)
			VALUES ($1, (SELECT id FROM items WHERE name=$2), $3, $4)`,
			orderID, item.Name, item.Quantity, item.Price*float64(item.Quantity))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrdersByPeriod(start, end string) ([]models.Order, error) {
	var orders []models.Order
	query := `
		SELECT shop_id, address, date, total_amount
		FROM orders
		WHERE date BETWEEN $1 AND $2
	`
	err := r.DB.Select(&orders, query, start, end)
	return orders, err
}

func (r *OrderRepository) GetShops() ([]string, error) {
	var shops []string
	query := `
		SELECT DISTINCT address
		FROM orders
	`
	err := r.DB.Select(&shops, query)
	return shops, err
}

func (r *OrderRepository) GetRevenueByShop() (map[string]float64, error) {
	rows, err := r.DB.Queryx(`
		SELECT address, SUM(total_amount) AS revenue
		FROM orders
		GROUP BY address
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	revenue := make(map[string]float64)
	for rows.Next() {
		var address string
		var total float64
		if err := rows.Scan(&address, &total); err != nil {
			return nil, err
		}
		revenue[address] = total
	}
	return revenue, nil
}

func (r *OrderRepository) GetAverageCheckByShop() (map[string]float64, error) {
	rows, err := r.DB.Queryx(`
		SELECT address, AVG(total_amount) AS avg_check
		FROM orders
		GROUP BY address
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	averageCheck := make(map[string]float64)
	for rows.Next() {
		var address string
		var avg float64
		if err := rows.Scan(&address, &avg); err != nil {
			return nil, err
		}
		averageCheck[address] = avg
	}
	return averageCheck, nil
}
