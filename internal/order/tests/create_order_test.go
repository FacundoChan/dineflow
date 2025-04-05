package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	sw "github.com/FacundoChan/gorder-v1/common/client/order"
	_ "github.com/FacundoChan/gorder-v1/common/config"
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
	requestBody := sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S38OduC65V5pGR",
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
	assert.Equal(t, 0, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_prams(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: customerID,
		Items:      nil,
	}
	response := getResponse(t, customerID, requestBody)

	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 2, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_item_id(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S38OduC65V5pGR",
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
	assert.Equal(t, 2, response.JSON200.Errorno)
}

func TestCreateOrder_invalid_item_quantity(t *testing.T) {
	requestBody := sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: customerID,
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_S38OduC65V5pGR",
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
	assert.Equal(t, 2, response.JSON200.Errorno)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("getResponse body=%+v", body)
	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}
