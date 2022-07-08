package entity

type Item struct {
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	OrderId     int    `json:"order_id"`
}

type ItemType struct {
	ItemCode    string `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
