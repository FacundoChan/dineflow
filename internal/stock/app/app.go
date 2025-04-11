package app

import "github.com/FacundoChan/gorder-v1/stock/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
	GetItems            query.GetItemsHandler
	GetAllProducts      query.GetAllProductsHandler // temporary
}
