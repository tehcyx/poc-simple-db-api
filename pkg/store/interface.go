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
	RawDataEvent    []byte `gorm:"-" json:"-"`
	RawDataCommerce []byte `gorm:"-" json:"-"`
}

// StorageData should hold arbitrary data, nothing fine-grained at the moment
type StorageData struct {
	Order
}

// Storage is an interface to support handling of different storage options
type Storage interface {
	Write(context.Context, StorageData) error
	ReadAll(context.Context) ([]StorageData, error)
}

type commerceResponse struct {
}

func (o Order) Validate() error {
	if o.BaseSiteUID == "" {
		return fmt.Errorf("Order not valid without a baseSiteUid")
	}
	if o.OrderCode == "" {
		return fmt.Errorf("Order not valid without a orderCode")
	}
	return nil
}

func (o Order) Enrich(ctx context.Context, url string) error {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error("Constructing API request failed")
		return fmt.Errorf("Couldn't create request")
	}
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

	var cresp commerceResponse
	marshErr := json.Unmarshal(responseByteArray, &cresp)
	if marshErr != nil {
		log.Error(marshErr)
		return fmt.Errorf("Failed to parse response json: %w", marshErr)
	}

	// fill orderdata fields here with more info
	// o.Firstname = ... cresp.something.firstname
	// o.Lastname = ... cresp.something.lastname
	o.RawDataCommerce = responseByteArray

	return nil
}
