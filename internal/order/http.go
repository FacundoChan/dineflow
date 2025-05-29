package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	client "github.com/FacundoChan/dineflow/common/client/order"
	"github.com/FacundoChan/dineflow/common/consts"

	"github.com/FacundoChan/dineflow/common"
	"github.com/FacundoChan/dineflow/common/handler/errors"
	"github.com/FacundoChan/dineflow/order/app"
	"github.com/FacundoChan/dineflow/order/app/command"
	"github.com/FacundoChan/dineflow/order/app/dto"
	"github.com/FacundoChan/dineflow/order/app/query"
	"github.com/FacundoChan/dineflow/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var (
		err      error
		request  client.CreateOrderRequest
		response dto.CreateOrderResponse
	)

	defer func() {
		H.Response(c, err, response)
	}()
	if err = c.ShouldBindJSON(&request); err != nil {
		err = errors.NewWithError(consts.ErrnoBindRequestError, err)
		return
	}
	if err = H.validate(request); err != nil {
		if request.Items == nil {
			err = errors.NewWithError(consts.ErrnoRequestNilItemsError, err)
		} else {
			err = errors.NewWithError(consts.ErrnoRequestValidateError, err)
		}
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

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
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

func (H HTTPServer) GetProducts(c *gin.Context) {
	var (
		err            error
		productsResult *query.GetProductsResult
		products       []dto.ProductDTO
		response       dto.GetProductsResponse
	)

	defer func() {
		H.Response(c, err, response)
	}()

	productsResult, err = H.app.Queries.GetProducts.Handle(c.Request.Context(), query.GetProducts{})
	logrus.Debug("productsResult", productsResult)

	if err != nil {
		return
	}

	for _, item := range productsResult.Products {
		products = append(products, dto.ProductDTO{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			ImgUrls:  item.ImgUrls,
		})
	}

	response.Products = products
}

func (H HTTPServer) validate(request client.CreateOrderRequest) error {
	if request.Items == nil {
		return fmt.Errorf("request.Items is nil")
	}
	for _, v := range request.Items {
		if v.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive, got %d from %s", v.Quantity, v.Id)
		}
	}
	return nil
}
