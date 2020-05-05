package simpledbapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
	"github.com/tehcyx/simple-db-api/pkg/store"
)

// SimpleDBAPI struct holding info about connecting services
type SimpleDBAPI struct {
	CommerceURL string
	KymaURL     string

	dataStore store.Storage
}

// NewSimpleDBAPI returns an instance of SimpleDBAPI
func NewSimpleDBAPI() *SimpleDBAPI {
	cs := new(SimpleDBAPI)
	return cs
}

// WithStorage chainable to a NewSimpleDBAPI to set a storage
func (svc *SimpleDBAPI) WithStorage(st store.Storage) *SimpleDBAPI {
	svc.dataStore = st
	return svc
}

// IndexHandler handles root entrypoint
func (svc *SimpleDBAPI) IndexHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Info("index handler")

	w.Write([]byte("Hello world"))
}

// CreateHandler handles POST requests to persist data into a datastore
func (svc *SimpleDBAPI) CreateHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Info("create handler")

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
		return
	}

	log.Debug("parsing request body")
	reqBody, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Info(readErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read content of request")
		return
	}

	log.Debug("transforming request body")
	var mappedData store.StorageData
	marshErr := json.Unmarshal(reqBody, &mappedData)
	if marshErr != nil {
		log.Info(marshErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to unmarshal json data")
		return
	}
	mappedData.CreatedAt = time.Now()
	mappedData.RawDataEvent = reqBody

	log.Debug("persisting request")
	storErr := svc.dataStore.Write(r.Context(), mappedData)
	if storErr != nil {
		log.Info(storErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to persist data due to some internal problem")
		return
	}

	log.Debug("done > create")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mappedData)

}

// ReadHandler handles GET requests to read data from datastore persistence
func (svc *SimpleDBAPI) ReadHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Info("read handler")

	log.Debug("reading data")
	data, storErr := svc.dataStore.ReadAll(r.Context())
	if storErr != nil {
		log.Info(storErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read data due to some internal problem")
		return
	}

	log.Debug("done > read")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
