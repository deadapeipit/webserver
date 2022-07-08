package entity

import "time"

type Order struct {
	OrderId      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
}

type OrderWithItems struct {
	Order
	Items []Item `json:"items"`
}
