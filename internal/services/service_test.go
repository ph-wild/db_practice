package services

import (
	"context"
	"testing"
	"time"

	"db_practice/internal/models"

	"github.com/stretchr/testify/assert"
)

type FakeOrderRepository struct { // implement interface
	Orders []models.Payment
}

func (f *FakeOrderRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	return nil
}

func (f *FakeOrderRepository) GetShops(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (f *FakeOrderRepository) GetRevenueByShop(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}
func (f *FakeOrderRepository) GetAverageCheckByShop(ctx context.Context) (map[string]float64, error) {
	return nil, nil
}

func (f *FakeOrderRepository) GetOrdersByPeriod(ctx context.Context, start, end time.Time) ([]models.Payment, error) {
	if end.After(start.AddDate(0, 2, 0)) {
		return nil, ErrTooLongPeriod
	}
	return f.Orders, nil
}

func TestService_GetOrdersByPeriod_TooLongPeriod(t *testing.T) {
	fakeRepo := &FakeOrderRepository{}
	service := NewService(fakeRepo)

	start := time.Now()
	end := start.AddDate(0, 3, 0)

	orders, err := service.GetOrdersByPeriod(context.Background(), start, end)

	assert.Nil(t, orders)
	assert.ErrorIs(t, err, ErrTooLongPeriod)
}

func TestService_GetOrdersByPeriod_ValidPeriod(t *testing.T) {
	start := time.Now()
	end := start.AddDate(0, 2, 0)
	expectedOrders := []models.Payment{
		{ShopID: 1, Address: "Test", Date: start.String(), TotalAmount: 100},
	}

	fakeRepo := &FakeOrderRepository{Orders: expectedOrders}
	service := NewService(fakeRepo)

	orders, err := service.GetOrdersByPeriod(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
}
