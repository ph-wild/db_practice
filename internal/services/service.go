package services

import (
	"time"

	"db_practice/internal/models"
	"db_practice/internal/repository"
)

type Service struct {
	repo *repository.OrderRepository
}

func NewService(repo *repository.OrderRepository) Service {
	return Service{repo: repo}
}

func (s *Service) SaveOrder(order *models.Order) error { //errors wrap or nil
	err := s.repo.SaveOrder(order)
	return err
}

func (s *Service) GetOrdersByPeriod(startTimeFormatted, endTimeFormatted time.Time) ([]models.Payment, error) {
	orders, err := s.repo.GetOrdersByPeriod(startTimeFormatted, endTimeFormatted)
	return orders, err
}

func (s *Service) GetShops() ([]string, error) {
	str, err := s.repo.GetShops()
	return str, err
}

func (s *Service) GetRevenueByShop() (map[string]float64, error) {
	rev, err := s.repo.GetRevenueByShop()
	return rev, err
}

func (s *Service) GetAverageCheckByShop() (map[string]float64, error) {
	aver, err := s.repo.GetAverageCheckByShop()
	return aver, err
}
