package store

import (
	"context"
	"time"
)

// StorageData should hold arbitrary data, nothing fine-grained at the moment
type StorageData struct {
	Date time.Time `json:"date"`
	Data []byte    `json:"data"`
}

// Storage is an interface to support handling of different storage options
type Storage interface {
	Write(context.Context, StorageData) error
	ReadAll(context.Context) ([]StorageData, error)
}
