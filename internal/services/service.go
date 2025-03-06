package services

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"db_practice/internal/models"
	"db_practice/internal/repository"
)

var ErrTooLongPeriod = errors.New("Maximum period is 2 mounth")

type OrderRepositoryInterface interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	GetOrdersByPeriod(ctx context.Context, start, end time.Time) ([]models.Payment, error)
	GetShops(ctx context.Context) ([]string, error)
	GetRevenueByShop(ctx context.Context) (map[string]float64, error)
	GetAverageCheckByShop(ctx context.Context) (map[string]float64, error)
}

type Service struct {
	Repo OrderRepositoryInterface // *repository.OrderRepository
}

func NewService(repo *repository.OrderRepository) Service {
	return Service{Repo: repo}
}

func (s *Service) SaveOrder(ctx context.Context, order *models.Order) error { //errors wrap or nil
	err := s.Repo.SaveOrder(ctx, order)
	return err
}

func (s *Service) GetOrdersByPeriod(ctx context.Context, startTimeFormatted, endTimeFormatted time.Time) ([]models.Payment, error) {
	maxPeriod := startTimeFormatted.AddDate(0, 2, 0)
	if endTimeFormatted.After(maxPeriod) {
		return nil, ErrTooLongPeriod
	}
	orders, err := s.Repo.GetOrdersByPeriod(ctx, startTimeFormatted, endTimeFormatted)
	return orders, err // err != nil -> to wrap
}

func (s *Service) GetShops(ctx context.Context) ([]string, error) {
	str, err := s.Repo.GetShops(ctx)
	return str, err // err != nil -> to wrap
}

func (s *Service) GetRevenueByShop(ctx context.Context) (map[string]float64, error) {
	rev, err := s.Repo.GetRevenueByShop(ctx)
	return rev, err // err != nil -> to wrap
}

func (s *Service) GetAverageCheckByShop(ctx context.Context) (map[string]float64, error) {
	aver, err := s.Repo.GetAverageCheckByShop(ctx)
	return aver, err // err != nil -> to wrap
}
