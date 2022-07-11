package entity

type Subscription struct {
	Plan          string `json:"plan"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	PaymentTerm   string `json:"term"`
}
