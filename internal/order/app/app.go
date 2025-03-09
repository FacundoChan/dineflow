package app

import "github.com/FacundoChan/gorder-v1/order/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct{}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
