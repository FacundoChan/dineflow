package entity

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*Item
}

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
	Price    float32
}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}
