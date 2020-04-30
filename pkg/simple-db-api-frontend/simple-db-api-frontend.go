package simpledbapifrontend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/internal/tmpl"
	"github.com/tehcyx/simple-db-api/pkg/logging"
	"github.com/tehcyx/simple-db-api/pkg/store"
)

// SimpleDBAPIFrontend struct holding info about connecting services
type SimpleDBAPIFrontend struct {
	BackendURL string
}

var (
	templates = template.Must(template.New("").
		Funcs(template.FuncMap{
			"toString":     func(data []byte) string { return string(data) },
			"base64Decode": func(input string) string { data, _ := base64.StdEncoding.DecodeString(input); return string(data) },
			"safe":         func(s string) template.HTML { return template.HTML(s) },
			"attr":         func(s string) template.HTMLAttr { return template.HTMLAttr(s) },
		}).Parse(""))
)

func init() {
	for _, tpl := range tmpl.TMPLMap {
		templates = template.Must(templates.Parse(tpl))
	}
}

// NewSimpleDBAPIFrontend returns an instance of SimpleDBAPI
func NewSimpleDBAPIFrontend() *SimpleDBAPIFrontend {
	cs := new(SimpleDBAPIFrontend)
	return cs
}

// IndexHandler handles root entrypoint
func (svc *SimpleDBAPIFrontend) IndexHandler(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Info("index handler")

	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, fmt.Sprintf("http://%s/read", svc.BackendURL), nil)
	if err != nil {
		log.Debug("API request failed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occured"))
	}
	response, clientErr := client.Do(req)
	if clientErr != nil {
		log.Info(clientErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to do request to backend")
		return
	}
	defer response.Body.Close()

	// Reading the response
	responseByteArray, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Info(readErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to read content of request")
		return
	}
	var events []store.StorageData
	marshErr := json.Unmarshal(responseByteArray, &events)
	if marshErr != nil {
		log.Info(marshErr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "The received data could not be read")
		return
	}

	if err := templates.ExecuteTemplate(w, "home", map[string]interface{}{
		"request_id": r.Context().Value(logging.CtxKeyRequestID{}),
		"APIData":    events,
	}); err != nil {
		log.Error(err)
	}
}
