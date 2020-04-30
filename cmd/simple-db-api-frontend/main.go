package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
	frontend "github.com/tehcyx/simple-db-api/pkg/simple-db-api-frontend"
	"github.com/tehcyx/simple-db-api/pkg/simple-db-api/cmd"
	"github.com/tehcyx/simple-db-api/pkg/util"
)

const (
	port = "8081"
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
	log.Infof("poc-simple-db-api-frontend %s@%s started", cmd.Version, cmd.GitCommit)
	flag.Parse()

	srvPort := port
	if os.Getenv("PORT") != "" {
		srvPort = os.Getenv("PORT")
	}
	addr := os.Getenv("LISTEN_ADDR")
	svc := frontend.NewSimpleDBAPIFrontend()
	util.MustMapEnv(&svc.BackendURL, "BACKEND_URL")

	r := mux.NewRouter()
	r.HandleFunc("/", svc.IndexHandler).Methods(http.MethodHead, http.MethodGet)

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	var handler http.Handler = r
	handler = &logging.LogHandler{Log: log, Next: handler} // add logging

	log.Infof("starting simple-db-api-frontend on " + addr + ":" + srvPort)
	log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}
