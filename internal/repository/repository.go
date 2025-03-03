package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"db_practice/internal/models"
)

var NoRows = errors.New("No rows selection")

type OrderRepository struct {
	DB *sqlx.DB
}

func NewOrderRepository(DB *sqlx.DB) *OrderRepository {
	return &OrderRepository{DB: DB}
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := r.DB.BeginTxx(dbCtx, nil)
	if err != nil {
		return err
	}

	var orderID int64
	err = tx.QueryRowxContext(ctx, `INSERT INTO orders (shop_id, address, date, total_amount)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		order.Payment.ShopID, order.Payment.Address, order.Payment.Date, order.Payment.TotalAmount).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertItemStmt, err := tx.PreparexContext(ctx, `INSERT INTO items (name, price) 
		VALUES ($1, $2) ON CONFLICT (name) DO NOTHING`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer insertItemStmt.Close()

	insertOrderItemStmt, err := tx.PreparexContext(ctx, `INSERT INTO order_items (order_id, item_id, quantity, total_price)
		VALUES ($1, (SELECT id FROM items WHERE name=$2), $3, $4)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer insertOrderItemStmt.Close()

	for _, item := range order.Payment.Items {
		_, err := insertItemStmt.ExecContext(ctx, item.Name, item.Price)
		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = insertOrderItemStmt.ExecContext(ctx, orderID, item.Name, item.Quantity, item.Price*float64(item.Quantity))
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrdersByPeriod(ctx context.Context, start, end time.Time) ([]models.Payment, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	var orders []models.Payment
	query := `
		SELECT shop_id, address, date, total_amount
		FROM orders
		WHERE date BETWEEN $1 AND $2
	`
	err := r.DB.SelectContext(dbCtx, &orders, query, start, end)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "Can't perform select")
	}

	return orders, nil
}

func (r *OrderRepository) GetShops(ctx context.Context) ([]string, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	var shops []string
	query := `
		SELECT DISTINCT address
		FROM orders
	`
	err := r.DB.SelectContext(dbCtx, &shops, query)
	return shops, err
}

func (r *OrderRepository) GetRevenueByShop(ctx context.Context) (map[string]float64, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	rows, err := r.DB.QueryxContext(dbCtx, `
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

func (r *OrderRepository) GetAverageCheckByShop(ctx context.Context) (map[string]float64, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rows, err := r.DB.QueryxContext(dbCtx, `
		SELECT address, AVG(total_amount) AS avg_check
		FROM orders
		GROUP BY address
	`)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errors.Wrap(err, "repo: GetAverageCheckByShop ended by timeout")
	}
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
