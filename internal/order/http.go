package main

import (
	"errors"
	"fmt"

	client "github.com/FacundoChan/gorder-v1/common/client/order"

	"github.com/FacundoChan/gorder-v1/common"
	"github.com/FacundoChan/gorder-v1/order/app"
	"github.com/FacundoChan/gorder-v1/order/app/command"
	"github.com/FacundoChan/gorder-v1/order/app/dto"
	"github.com/FacundoChan/gorder-v1/order/app/query"
	"github.com/FacundoChan/gorder-v1/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	var (
		err      error
		request  client.CreateOrderRequest
		response dto.CreateOrderResponse
	)

	defer func() {
		H.Response(c, err, response)
	}()
	if err = c.ShouldBindJSON(&request); err != nil {
		return
	}
	if err = H.validate(request); err != nil {
		return
	}

	result, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: request.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(request.Items),
	})

	if err != nil {
		return
	}

	response = dto.CreateOrderResponse{
		OrderID:     result.OrderID,
		CustomerID:  request.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", request.CustomerId, result.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
	var (
		err      error
		response struct {
			Order *client.Order
		}
	)

	defer func() {
		H.Response(c, err, response)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerID: customerID,
		OrderID:    orderID,
	})

	if err != nil {
		return
	}

	response.Order = convertor.NewOrderConvertor().EntityToClient(o)
}

func (H HTTPServer) validate(request client.CreateOrderRequest) error {
	for _, v := range request.Items {
		if v.Quantity <= 0 {
			return errors.New("quantity must be positive")
		}
	}
	return nil
}
