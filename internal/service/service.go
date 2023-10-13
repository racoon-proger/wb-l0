package service

import (
	"context"

	"github.com/racoon-proger/wb-l0/internal/domain"
)

type cache interface {
	GetOrder(id int) *domain.Order
	SetOrder(order *domain.Order)
}

type storage interface {
	CreateOrder(ctx context.Context, order *domain.Order) (err error)
}

type service struct {
	cache   cache
	storage storage
}
//
func (s *service) GetOrderByID(ctx context.Context, id int) (order *domain.Order, err error) {
	order = s.cache.GetOrder(id)
	return
}

func (s *service) CreateOrder(ctx context.Context, order *domain.Order) (err error) {
	err = s.storage.CreateOrder(ctx, order)
	if err != nil {
		return err
	}
	s.cache.SetOrder(order)
	return err
}

func NewService(cache cache, storage storage) *service {
	return &service{
		cache:   cache,
		storage: storage,
	}
}
