package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/racoon-proger/wb-l0/internal/domain"
)

// Storage this is an order storage
type Storage struct {
	db *sql.DB
}

// CreateOrder inserts data into the order table
func (s *Storage) CreateOrder(
	ctx context.Context,
	order *domain.Order,
) (err error) {
	var delivery, payment, items []byte
	delivery, err = json.Marshal(&order.Delivery)
	if err != nil {
		return err
	}
	payment, err = json.Marshal(&order.Payment)
	if err != nil {
		return err
	}
	items, err = json.Marshal(&order.Items)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO orders (
			order_uuid,
			track_number,
			entry,
			locale,
			internal_signature,
			customer_id,
			delivery_service,
			shardkey,
			sm_id,
			date_created,
			oof_shard,
			delivery,
			payment,
			items
			)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
		delivery,
		payment,
		items,
	)
	return err
}

// GetOrders fills the cache from the database.
func (s *Storage) GetOrders(ctx context.Context) (orders []domain.Order, err error) {
	rows, err := s.db.Query(`SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			order                    domain.Order
			delivery, payment, items []byte
		)
		err = rows.Scan(
			&order.ID,
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&delivery,
			&payment,
			&items,
		)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(delivery, &order.Delivery)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(payment, &order.Payment)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(items, &order.Items)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}
