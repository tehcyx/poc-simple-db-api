package store

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
)

// InMemory simply holds data for the runtime of the application
type InMemory struct {
	data []Order
}

// NewInMemoryStore returns an instance of an inmemory store
func NewInMemoryStore() *InMemory {
	return &InMemory{}
}

// Write writes the storage object to the in memory store
func (im *InMemory) Write(ctx context.Context, data Order) error {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Debugf("writing: %+v", data)
	im.data = append(im.data, data)
	return nil
}

// ReadAll returns all data stored in memory
func (im *InMemory) ReadAll(ctx context.Context) ([]Order, error) {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Debugf("reading: %+v", im.data)
	return im.data, nil
}
