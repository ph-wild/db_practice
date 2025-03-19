package cache

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"db_practice/internal/models"
	"db_practice/internal/services"

	"github.com/pkg/errors"
)

type Cache struct {
	data       map[string]any
	repository services.OrderRepositoryInterface
	mu         sync.RWMutex
}

func NewCache(c services.OrderRepositoryInterface) *Cache {
	return &Cache{
		data:       make(map[string]any),
		repository: c,
	}
}

func (c *Cache) set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

func (c *Cache) get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, found := c.data[key]
	return val, found
}

func (c *Cache) SaveOrder(ctx context.Context, order *models.Order) error {
	return c.repository.SaveOrder(ctx, order)
}

func (c *Cache) GetShops(ctx context.Context) ([]string, error) {
	key := "all_shops"
	data, found := c.get(key)
	if found {
		shops, err := data.([]string)
		if err {
			return shops, errors.New("can't convert any to []strings")
		}
		return shops, nil
	}

	shops, err := c.repository.GetShops(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cache: can't get shops")
	}

	c.set(key, strings.Join(shops, ", "))
	return shops, nil
}

func (c *Cache) GetRevenueByShop(ctx context.Context) (map[string]float64, error) {
	key := "shops_revenue"
	data, found := c.get(key)
	if found {
		revenue, err := data.(map[string]float64)
		if err {
			return revenue, errors.New("can't convert any to map[string]float64")
		}
		return revenue, nil
	}

	revenue, err := c.repository.GetRevenueByShop(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cache: can't get revenue by shops")
	}

	revenueArr := []string{}
	for k, v := range revenue {
		revenueArr = append(revenueArr, k+": "+strconv.FormatFloat(v, 'f', 4, 64))
	}

	c.set(key, strings.Join(revenueArr, ", "))
	return revenue, nil
}

func (c *Cache) GetAverageCheckByShop(ctx context.Context) (map[string]float64, error) {
	key := "avarage_check"
	data, found := c.get(key)
	if found {
		averageCheck, err := data.(map[string]float64)
		if err {
			return averageCheck, errors.New("can't convert any to map[string]float64")
		}
		return averageCheck, nil
	}

	averageCheck, err := c.repository.GetAverageCheckByShop(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cache: can't get avarage check by shops")
	}

	averageCheckArr := []string{}
	for k, v := range averageCheck {
		averageCheckArr = append(averageCheckArr, k+": "+strconv.FormatFloat(v, 'f', 4, 64))
	}

	c.set(key, strings.Join(averageCheckArr, ", "))
	return averageCheck, nil
}

func (c *Cache) GetOrdersByPeriod(ctx context.Context, start, end time.Time) ([]models.Payment, error) {
	key := "orders from " + start.String() + " to " + end.String()
	data, found := c.get(key)
	if found {
		ordersByPeriod, err := data.([]models.Payment)
		if err {
			return ordersByPeriod, errors.New("can't convert any to []models.Payment")
		}
		return ordersByPeriod, nil
	}

	ordersByPeriod, err := c.repository.GetOrdersByPeriod(ctx, start, end)
	if err != nil {
		return nil, errors.Wrap(err, "cache: can't get orders by period")
	}
	return ordersByPeriod, nil
}
