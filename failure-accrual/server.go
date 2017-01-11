package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests",
	Help: "Number of requests",
})

type handler struct {
	successRate float64
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	requests.Inc()
	if rand.Float64() > h.successRate {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal service error"))
		return
	}
	w.Write([]byte("pong"))
}

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	addr := flag.String("addr", ":8501", "service port to run on")
	successRate := flag.Float64("success-rate", 1.0, "service success rate")
	flag.Parse()

	fmt.Printf("serving on %s with %2f success rate\n", *addr, *successRate)

	httpHandler := handler{successRate: *successRate}
	http.HandleFunc("/", httpHandler.HandleRequest)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*addr, nil)
}
