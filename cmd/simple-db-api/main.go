package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
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

	if url := os.Getenv("GATEWAY_URL"); url != "" {
		svc.CommerceURL = url
	} else {
		log.Infof("GATEWAY_URL is not set, skipping calls to commerce backend")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.IndexHandler).Methods(http.MethodHead, http.MethodGet)
	r.HandleFunc("/create", svc.CreateHandler).Methods(http.MethodHead, http.MethodPost)
	r.HandleFunc("/read", svc.ReadHandler).Methods(http.MethodHead, http.MethodGet)

	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	var handler http.Handler = r
	handler = &logging.LogHandler{Log: log, Next: handler} // add logging

	if os.Getenv("DISABLE_PROFILER") == "" {
		log.Info("Profiling enabled")
		r.PathPrefix("/debug").Handler(http.DefaultServeMux)
	} else {
		log.Info("Profiling disabled")
	}

	log.Infof("starting simple-db-api on " + addr + ":" + srvPort)
	log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}
