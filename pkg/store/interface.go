package store

import (
	"context"
	"time"
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
	BaseSiteID      string `json:"baseSiteId"`
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
