package cache

import (
	"sync"

	"github.com/racoon-proger/wb-l0/internal/domain"
)

type cache struct {
	sync.RWMutex
	orders map[int]*domain.Order
}

// SetOrder saves an order in the cache
func (c *cache) SetOrder(order *domain.Order) {
	c.Lock()
	c.orders[order.ID] = order
	c.Unlock()
}

// GetOrder retrieves data from cache
func (c *cache) GetOrder(id int) *domain.Order {
	c.RLock()
	order := c.orders[id]
	c.RUnlock()
	return order
}

func NewCache() *cache {
	return &cache{
		orders: make(map[int]*domain.Order),
	}
}
