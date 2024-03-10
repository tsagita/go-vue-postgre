package domain

import (
	"database/sql"
)

type OrderItem struct {
	OrderID      int     `json:"order_id"`
	Product      string  `json:"product"`
	OrderName    string  `json:"order_name"`
	CreatedAt    string  `json:"created_at"`
	DeliveryAmt  sql.NullFloat64 `json:"delivery_amount"`
	TotalAmount  sql.NullFloat64 `json:"total_amount"`
	CustomerName string  `json:"customer_name"`
	CompanyName  string  `json:"company_name"`
}