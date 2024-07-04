package models

type Order struct {
	Id         string  `json:"id"`
	UserId     string  `json:"user_id"`
	ProductId  string  `json:"product_id"`
	Quantity   int32   `json:"quantity"`
	Status     string  `json:"status"`
	TotalPrice float32 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type CreateOrder struct {
	UserId     string  `json:"user_id"`
	ProductId  string  `json:"product_id"`
	Quantity   int32   `json:"quantity"`
	Status     string  `json:"status"`
	TotalPrice float32 `json:"total_price"`
}

type UpdateOrder struct {
	UserId     string  `json:"-"`
	ProductId  string  `json:"-"`
	Quantity   int32   `json:"quantity"`
	Status     string  `json:"status"`
	TotalPrice float32 `json:"total_price"`
}

type ListOrders struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
