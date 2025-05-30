package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/FacundoChan/dineflow/common/consts"

	sw "github.com/FacundoChan/dineflow/common/client/order"
	_ "github.com/FacundoChan/dineflow/common/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	ctx        = context.Background()
	server     = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
	customerID = "customerID-test"
)

func TestMain(m *testing.M) {
	before()
	m.Run()

}

func before() {
	log.Printf("server=%s", server)
}

func TestCreateOrder_success(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIDOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S3CrGrzAS1MZsK",
				Quantity: 10,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: 5,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: 1,
			},
		},
	}

	response := getResponse(t, customerID, requestBody)

	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, consts.ErrnoSuccess, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_prams(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIDOrdersJSONRequestBody{
		CustomerId: customerID,
		Items:      nil,
	}
	response := getResponse(t, customerID, requestBody)

	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, consts.ErrnoRequestNilItemsError, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_item_id(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIDOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S3CrGrzAS1MZsK",
				Quantity: 10,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: 5,
			},
			{
				Id:       "invalid_id",
				Quantity: 1,
			},
		},
	}

	response := getResponse(t, customerID, requestBody)

	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, consts.ErrnoStripeResourceMissingError, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_item_quantity(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIDOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S3CrGrzAS1MZsK",
				Quantity: 10,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: -5,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: 1,
			},
		},
	}

	response := getResponse(t, customerID, requestBody)

	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, consts.ErrnoRequestValidateError, response.JSON200.Errorno)
}

func TestCreateOrder_idempotent(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIDOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S3CrGrzAS1MZsK",
				Quantity: 2,
			},
			{
				Id:       "prod_S3Cr3l2WHdiL53",
				Quantity: 3,
			},
		},
	}

	// First order
	response1 := getResponse(t, customerID, requestBody)
	assert.Equal(t, 200, response1.StatusCode())
	assert.Equal(t, consts.ErrnoSuccess, response1.JSON200.Errorno)
	orderID1, ok1 := response1.JSON200.Data["order_id"].(string)
	assert.True(t, ok1)
	assert.NotEmpty(t, orderID1)

	// Second order with the same content
	response2 := getResponse(t, customerID, requestBody)
	assert.Equal(t, 200, response2.StatusCode())
	assert.Equal(t, consts.ErrnoSuccess, response2.JSON200.Errorno)
	orderID2, ok2 := response2.JSON200.Data["order_id"].(string)
	assert.True(t, ok2)
	assert.NotEmpty(t, orderID2)

	assert.Equal(t, orderID1, orderID2)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIDOrdersJSONRequestBody) *sw.PostCustomerCustomerIDOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("getResponse body=%+v", body)
	response, err := client.PostCustomerCustomerIDOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
