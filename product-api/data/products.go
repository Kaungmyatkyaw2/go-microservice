package data

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protos "github.com/Kaungmyatkyaw2/go-microservice/currency/protos/currency"
)

// Product defines the strcture for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

type ProductsDB struct {
	currency protos.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client   protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {

	pb := &ProductsDB{c, l, make(map[string]float64), nil}
	go pb.handleUpdates()
	return pb
}

func (p *ProductsDB) handleUpdates() {
	sub, err := p.currency.SubscribeRates(context.Background())

	if err != nil {
		p.log.Error("Unable to subscribe for rates", "error", err)
	}

	p.client = sub

	for {
		rr, err := sub.Recv()

		if grpcError := rr.GetError(); grpcError != nil {
			p.log.Error("Error subscribing for rates", "error", err)
			continue
		}

		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("Received updated rate from server", "dest", resp.GetDestination().String())
			if err != nil {
				p.log.Error("Error receiving message", "error", err)
				return
			}
			p.rates[resp.Destination.String()] = resp.Rate

		}

	}

}

func (p *ProductsDB) GetProducts(currency string) (Products, error) {
	if currency == "" {
		return productList, nil
	}

	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * rate
		pr = append(pr, &np)
	}

	return pr, nil
}

func (p *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency == "" {
		return productList[i], nil
	}
	rate, err := p.getRate(currency)
	if err != nil {
		p.log.Error("Unable to get rate", "currency", currency)
		return nil, err
	}

	np := *productList[i]
	np.Price = np.Price * rate

	return &np, nil
}

func (p *ProductsDB) AddProduct(prod *Product) {
	prod.ID = getNextID()
	productList = append(productList, prod)
}

func (p *ProductsDB) UpdateProduct(id int, product *Product) error {
	_, pos, err := findProduct(id)

	if err != nil {
		return err
	}

	product.ID = id
	productList[pos] = product
	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	// To demo our subscription duplication err handling we have to comment off following cached mechanism
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	// get initial rate

	resp, err := p.currency.GetRate(context.Background(), rr)

	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("Unable to get rate from currency server, destination and base can't be the same, base:%s dest:%ss", md.Base.String(), md.Destination.String())

			}
			return -1, fmt.Errorf("Unable to get rate from currency server, base:%s dest:%ss", md.Base.String(), md.Destination.String())
		}

		return -1, err
	}

	p.rates[destination] = resp.Rate

	// subscribe for updated rates
	err = p.client.Send(rr)

	if err != nil {
		return -1, err
	}

	return resp.Rate, err

}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
