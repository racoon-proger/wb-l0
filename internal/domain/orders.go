package domain

import "time"

// order structure
type Order struct {
	ID                int          `json:"id,omitempty"`
	OrderUID          string       `json:"order_uid,omitempty"`
	TrackNumber       string       `json:"track_number,omitempty"`
	Entry             string       `json:"entry,omitempty"`
	Locale            string       `json:"locale,omitempty"`
	InternalSignature string       `json:"internal_signature,omitempty"`
	CustomerID        string       `json:"customer_id,omitempty"`
	DeliveryService   string       `json:"delivery_service,omitempty"`
	Shardkey          string       `json:"shardkey,omitempty"`
	SmID              int          `json:"sm_id,omitempty"`
	DateCreated       time.Time    `json:"date_created,omitempty"`
	OofShard          string       `json:"oof_shard,omitempty"`
	Delivery          Delivery     `json:"delivery,omitempty"`
	Payment           OrderPayment `json:"payment,omitempty"`
	Items             []OrderItem  `json:"items,omitempty"`
}
