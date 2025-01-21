package models

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Payment struct {
	Date        string  `json:"date"`
	ShopID      int     `json:"shop_id"`
	Address     string  `json:"address"`
	TotalAmount float64 `json:"total_amount"`
	Items       []Item  `json:"items"`
}

type Order struct {
	Payment Payment `json:"payment"`
}
