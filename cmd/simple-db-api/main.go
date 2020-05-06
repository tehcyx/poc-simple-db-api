package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
	service "github.com/tehcyx/simple-db-api/pkg/simple-db-api"
	"github.com/tehcyx/simple-db-api/pkg/simple-db-api/cmd"
	"github.com/tehcyx/simple-db-api/pkg/store"
	"github.com/tehcyx/simple-db-api/pkg/util"
)

const (
	port = "8080"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logging.DetermineLogLevel("LOG_LEVEL")
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

type status struct {
	Status string `json:"status"`
}

func main() {
	// ctx := context.Background()
	log.Infof("poc-simple-db-api %s@%s started", cmd.Version, cmd.GitCommit)
	flag.Parse()

	srvPort := port
	if os.Getenv("PORT") != "" {
		srvPort = os.Getenv("PORT")
	}
	addr := os.Getenv("LISTEN_ADDR")

	var DBUser, DBPass, DBDBase, DBHost, DBPort string
	util.MustMapEnv(&DBUser, "POSTGRES_USER")
	util.MustMapEnv(&DBPass, "POSTGRES_PASSWORD")
	util.MustMapEnv(&DBDBase, "POSTGRES_DB")
	util.MustMapEnv(&DBHost, "POSTGRES_HOST")
	util.MustMapEnv(&DBPort, "POSTGRES_PORT")

	svc := service.NewSimpleDBAPI().WithStorage(store.NewPostgresStore(log, DBUser, DBPass, DBHost, DBPort, DBDBase))
	// util.MustMapEnv(&svc.KymaURL, "KYMA_URL")
	// util.MustMapEnv(&svc.CommerceURL, "COMMERCE_URL")

	r := mux.NewRouter()
	r.HandleFunc("/", svc.IndexHandler).Methods(http.MethodHead, http.MethodGet)
	r.HandleFunc("/create", svc.CreateHandler).Methods(http.MethodHead, http.MethodPost)
	r.HandleFunc("/read", svc.ReadHandler).Methods(http.MethodHead, http.MethodGet)

	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status{Status: "ok"})
	})

	var handler http.Handler = r
	handler = &logging.LogHandler{Log: log, Next: handler} // add logging

	log.Infof("starting simple-db-api on " + addr + ":" + srvPort)
	log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}
