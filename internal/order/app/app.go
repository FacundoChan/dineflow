package app

import (
	"github.com/FacundoChan/gorder-v1/order/app/command"
	"github.com/FacundoChan/gorder-v1/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
	GetProducts      query.GetProductsHandler
}
