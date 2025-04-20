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

type Product struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
	Price    float32
	ImgUrls  []string
}

type StockModel struct {
	ID          int64
	ProductID   string
	Name        string
	Quantity    int64
	Price       float32
	Description string
	Version     int64
	ImgUrls     []string
}
