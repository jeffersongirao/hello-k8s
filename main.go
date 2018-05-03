package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.Parse(os.Args[1:])

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)

	successHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello Kubernetes."))
	})

	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	http.Handle("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, successHandler))
	log.WithFields(log.Fields{"handler": "success"}).Info("initialized")
	http.Handle("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, errorHandler))
	log.WithFields(log.Fields{"handler": "err"}).Info("initialized")

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.WithFields(log.Fields{"handler": "metrics"}).Info("initialized")
	log.Fatal(http.ListenAndServe(bind, nil))
}
