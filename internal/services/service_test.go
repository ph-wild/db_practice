package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"db_practice/internal/models"
)

type MockOrderRepository struct {
	mock.Mock // embedding - встраивание структуры в структуру, можно MockInstance mock.Mock, тогда m.MockInstance.On вместо m.On
}

func (m *MockOrderRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersByPeriod(ctx context.Context, start, end time.Time) ([]models.Payment, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]models.Payment), args.Error(1) // 1 (not 0) - error is second returning arg
}

func (m *MockOrderRepository) GetShops(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockOrderRepository) GetRevenueByShop(ctx context.Context) (map[string]float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func (m *MockOrderRepository) GetAverageCheckByShop(ctx context.Context) (map[string]float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]float64), args.Error(1)
}

func TestService_GetOrdersByPeriod_TooLongPeriod(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := Service{Repo: mockRepo}

	start := time.Now()
	end := start.AddDate(0, 3, 0)

	orders, err := service.GetOrdersByPeriod(context.Background(), start, end)
	assert.Nil(t, orders)
	assert.ErrorIs(t, err, ErrTooLongPeriod)
}

func TestService_GetOrdersByPeriod_ValidPeriod(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := Service{Repo: mockRepo}

	start := time.Now()
	end := start.AddDate(0, 2, 0)
	expectedOrders := []models.Payment{{ShopID: 1, Address: "Test", Date: start.String(), TotalAmount: 100}}

	// as db result
	mockRepo.On("GetOrdersByPeriod", mock.Anything, start, end).Return(expectedOrders, nil)

	orders, err := service.GetOrdersByPeriod(context.Background(), start, end) // call mockRepo.On
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	mockRepo.AssertExpectations(t)
}
