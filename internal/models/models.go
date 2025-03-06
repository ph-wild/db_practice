package models

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Payment struct {
	Date        string  `json:"date" db:"date"`
	ShopID      int     `json:"shop_id" db:"shop_id"`
	Address     string  `json:"address" db:"address"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	Items       []Item  `json:"items" db:"items"`
}

type Order struct {
	Payment Payment `json:"payment"`
}
