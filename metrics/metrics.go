package metrics

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	requestsFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_tester_requests_failures_total",
			Help: "The total number of request failures",
		},
		[]string{"handler", "method", "code", "endpoint"},
	)

	totalRequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_tester_http_requests_total",
			Help: "The total number of requests",
		},
		[]string{"handler", "method", "code", "endpoint"},
	)

	journeysFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_tester_journey_failures_total",
			Help: "The total number of journey failures",
		},
		[]string{"handler", "journey"},
	)

	totalJourneysCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_tester_journeys_total",
			Help: "The total number of journeys",
		},
		[]string{"handler", "journey"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_tester_request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method", "code", "endpoint"},
	)

	journeyDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_tester_journey_duration_seconds",
			Help:    "A histogram of latencies for different journeys.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "journey"},
	)
)

type Metrics struct {
	TotalRequestCounter    *prometheus.CounterVec
	RequestsFailureCounter *prometheus.CounterVec
	TotalJourneyCounter    *prometheus.CounterVec
	JourneysFailureCounter *prometheus.CounterVec
	RequestDuration        *prometheus.HistogramVec
	JourneyDuration        *prometheus.HistogramVec
}

func New() *Metrics {
	prometheus.MustRegister(totalRequestsCounter, requestsFailureCounter, requestDuration)
	prometheus.MustRegister(totalJourneysCounter, journeysFailureCounter, journeyDuration)

	return &Metrics{
		TotalRequestCounter:    totalRequestsCounter,
		RequestsFailureCounter: requestsFailureCounter,
		TotalJourneyCounter:    totalJourneysCounter,
		JourneysFailureCounter: journeysFailureCounter,
		RequestDuration:        requestDuration,
		JourneyDuration:        journeyDuration,
	}
}

func StartMetricsServer(host, port string) {
	s := http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: promhttp.Handler(),
	}

	go func() {
		err := s.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Metric server ListenAndServe: %v", err)
		}
	}()

	log.WithFields(log.Fields{
		"addr": s.Addr,
	}).Info("Metrics server listening for requests")
}
