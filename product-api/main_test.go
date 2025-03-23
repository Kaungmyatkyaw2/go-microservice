package main

import (
	"fmt"
	"testing"

	"github.com/Kaungmyatkyaw2/product-api/sdk/client"
	"github.com/Kaungmyatkyaw2/product-api/sdk/client/products"
)

func TestOurClient(t *testing.T) {

	cfg := client.DefaultTransportConfig()
	cfg = cfg.WithHost("localhost:9090")

	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod.GetPayload()[0])
	t.Fail()
}
