package store

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
)

// Model represents basic fields for all persisted values
type Model struct {
	ID        uint       `gorm:"primary_key,auto_increment" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// Order represents a Commerce Order
type Order struct {
	Model
	Firstname       string `json:"firstName"`
	Lastname        string `json:"lastName"`
	OrderCode       string `json:"orderCode"`
	BaseSiteUID     string `json:"baseSiteUid"`
	Items           []Item `json:"items"`
	RawDataEvent    []byte `gorm:"-" json:"-"`
	RawDataCommerce []byte `gorm:"-" json:"-"`
}

// Item represents a Commerce Order line item
type Item struct {
	Model
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	OrderID  uint
}

// StorageData should hold arbitrary data, nothing fine-grained at the moment
type StorageData struct {
	Order
}

// Storage is an interface to support handling of different storage options
type Storage interface {
	Write(context.Context, Order) error
	ReadAll(context.Context) ([]Order, error)
}

type CommerceResponse struct {
	DeliveryAddress Address `json:"deliveryAddress"`
	Entries         []Entry `json:"entries"`
}

type Address struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Entry struct {
	Quantity int     `json:"quantity"`
	Product  Product `json:"product"`
}

type Product struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Validate ensures, that the required parameters for subsequent calls are set
func (o Order) Validate() error {
	if o.BaseSiteUID == "" {
		return fmt.Errorf("Order not valid without a baseSiteUid")
	}
	if o.OrderCode == "" {
		return fmt.Errorf("Order not valid without a orderCode")
	}
	return nil
}

func (o *Order) Enrich(ctx context.Context, url string) error {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	if err := o.Validate(); err != nil {
		log.Errorf("Can't produce subsequent order detail call, minimum requirements are not met.")
		return fmt.Errorf("Can't produce subsequent order detail call, minimum requirements are not met: %w", err)
	}
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("Constructing API request failed")
		return fmt.Errorf("Couldn't create request")
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	response, clientErr := client.Do(req)
	if clientErr != nil {
		log.Error(clientErr)
		return fmt.Errorf("Executing request failed: %w", clientErr)
	}
	defer response.Body.Close()

	// Reading the response
	responseByteArray, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Error(readErr)
		return fmt.Errorf("Failed to read response: %w", readErr)
	}

	log.Debugf("Response from commerce: %+v", string(responseByteArray))

	o.RawDataCommerce = responseByteArray

	var cresp CommerceResponse
	marshErr := json.Unmarshal(responseByteArray, &cresp)
	if marshErr != nil {
		log.Error(marshErr)
		return fmt.Errorf("Failed to parse response json: %w", marshErr)
	}

	// fill orderdata fields here with more info
	o.Firstname = cresp.DeliveryAddress.FirstName
	o.Lastname = cresp.DeliveryAddress.LastName

	for _, i := range cresp.Entries {
		o.Items = append(o.Items, Item{Name: i.Product.Code, Quantity: i.Quantity})
	}

	return nil
}
