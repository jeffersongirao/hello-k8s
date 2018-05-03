package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
)

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
	http.Handle("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, errorHandler))

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(bind, nil))
}
