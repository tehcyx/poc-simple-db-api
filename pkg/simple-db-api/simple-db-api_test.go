package simpledbapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tehcyx/simple-db-api/pkg/logging"
	"github.com/tehcyx/simple-db-api/pkg/store"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.ErrorLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

func TestSimpleDBAPI_WithStorage(t *testing.T) {
	type fields struct {
		CommerceURL string
		KymaURL     string
		dataStore   store.Storage
	}
	type args struct {
		st store.Storage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *SimpleDBAPI
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &SimpleDBAPI{
				CommerceURL: tt.fields.CommerceURL,
				KymaURL:     tt.fields.KymaURL,
				dataStore:   tt.fields.dataStore,
			}
			if got := svc.WithStorage(tt.args.st); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SimpleDBAPI.WithStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleDBAPI_IndexHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	req = setupLogForRequest(req)

	NewSimpleDBAPI().IndexHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusOK)
	}
}

func TestSimpleDBAPI_CreateHandler(t *testing.T) {
	tests := []struct {
		name        string
		CommerceURL string
		request     *http.Request
		expect      int
	}{
		{name: "Get Request should fail",
			request:     httptest.NewRequest(http.MethodGet, "/create", nil),
			CommerceURL: "",
			expect:      http.StatusMethodNotAllowed},
		{name: "Post Request should not fail",
			request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{"orderCode": "hello", "baseSiteUid": "abc"}`))),
			CommerceURL: "",
			expect:      http.StatusCreated},
		{name: "Post Request without required parameters should fail 1/3",
			request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{}`))),
			CommerceURL: "",
			expect:      http.StatusBadRequest},
		{name: "Post Request without required parameters should fail 2/3",
			request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{"orderCode": "hello"}`))),
			CommerceURL: "",
			expect:      http.StatusBadRequest},
		{name: "Post Request without required parameters should fail 3/3",
			request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{"baseSiteUid": "abc"}`))),
			CommerceURL: "",
			expect:      http.StatusBadRequest},
		{name: "Post Request with broken json should fail",
			request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{`))),
			CommerceURL: "",
			expect:      http.StatusBadRequest},
		// !!Disabled due to subsequent request not yet mockable
		// {name: "Post Request with commerce url should not fail",
		// 	request:     httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{"orderCode": "hello", "baseSiteUid": "abc"}`))),
		// 	CommerceURL: "http://localhost",
		// 	expect:      http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &SimpleDBAPI{
				CommerceURL: tt.CommerceURL,
				dataStore:   store.NewInMemoryStore(),
			}
			res := httptest.NewRecorder()
			tt.request = setupLogForRequest(tt.request)
			svc.CreateHandler(res, tt.request)

			assert.Equal(t, res.Code, tt.expect, "got status %d but wanted %d", res.Code, tt.expect)
		})
	}
}

func TestSimpleDBAPI_ReadHandler(t *testing.T) {
	tests := []struct {
		name        string
		CommerceURL string
		request     *http.Request
		expect      int
	}{
		{name: "Get Request should not fail",
			request:     httptest.NewRequest(http.MethodGet, "/read", nil),
			CommerceURL: "",
			expect:      http.StatusOK},
		{name: "Post Request should fail",
			request:     httptest.NewRequest(http.MethodPost, "/read", bytes.NewBuffer([]byte{})),
			CommerceURL: "",
			expect:      http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &SimpleDBAPI{
				CommerceURL: tt.CommerceURL,
				dataStore:   store.NewInMemoryStore(),
			}
			res := httptest.NewRecorder()
			tt.request = setupLogForRequest(tt.request)
			svc.ReadHandler(res, tt.request)

			assert.Equal(t, res.Code, tt.expect, "got status %d but wanted %d", res.Code, tt.expect)
		})
	}
}

func setupLogForRequest(req *http.Request) *http.Request {
	ctx := req.Context()
	requestID, _ := uuid.NewRandom()
	ctx = context.WithValue(ctx, logging.CtxKeyRequestID{}, requestID.String())
	testLogger := log.WithFields(logrus.Fields{
		"http.req.path":   req.URL.Path,
		"http.req.method": req.Method,
		"http.req.id":     requestID.String(),
	})
	ctx = context.WithValue(ctx, logging.CtxKeyLog{}, testLogger)
	return req.WithContext(ctx)
}

func BenchmarkSimpleDBAPI_IndexHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	req = setupLogForRequest(req)

	for n := 0; n < b.N; n++ {
		NewSimpleDBAPI().IndexHandler(res, req)
	}
}

func BenchmarkSimpleDBAPI_ReadHandler(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/read", nil)
	res := httptest.NewRecorder()

	req = setupLogForRequest(req)
	svc := NewSimpleDBAPI().WithStorage(store.NewInMemoryStore())

	for n := 0; n < b.N; n++ {
		svc.ReadHandler(res, req)
	}
}

func BenchmarkSimpleDBAPI_CreateHandler(b *testing.B) {
	b.Run("Test empty (invalid) POST body", func(b *testing.B) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{}`)))
		res := httptest.NewRecorder()

		req = setupLogForRequest(req)
		svc := NewSimpleDBAPI().WithStorage(store.NewInMemoryStore())

		for n := 0; n < b.N; n++ {
			svc.CreateHandler(res, req)
		}
	})
	b.Run("Test valid POST body", func(b *testing.B) {
		req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer([]byte(`{"orderCode": "hello", "baseSiteUid": "abc"}`)))
		res := httptest.NewRecorder()

		req = setupLogForRequest(req)
		svc := NewSimpleDBAPI().WithStorage(store.NewInMemoryStore())

		for n := 0; n < b.N; n++ {
			svc.CreateHandler(res, req)
		}
	})
}
