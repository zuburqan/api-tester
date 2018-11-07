package transport

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/zuburqan/api-tester/metrics"
	"github.com/zuburqan/api-tester/utils"
)

func SetupRoundTrippers(metrics *metrics.Metrics, r http.RoundTripper, handler string) http.RoundTripper {
	return instrumentRoundTripperDuration(metrics.RequestDuration.MustCurryWith(prometheus.Labels{"handler": handler}),
		instrumentRoundTripperCounter(metrics.TotalRequestCounter, metrics.RequestsFailureCounter, r, handler), handler)
}

func instrumentRoundTripperCounter(totalRequestsCounter *prometheus.CounterVec,
	requestsFailureCounter *prometheus.CounterVec, next http.RoundTripper, handler string) promhttp.RoundTripperFunc {

	return promhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(r)

		if err != nil {
			totalRequestsCounter.With(prometheus.Labels{"handler": handler, "method": r.Method,
				"code": "Request Failed", "endpoint": utils.RemoveIDFrom(r.URL.Path, handler)}).Add(1)
			requestsFailureCounter.With(prometheus.Labels{"handler": handler, "method": r.Method,
				"code": "Request Failed", "endpoint": utils.RemoveIDFrom(r.URL.Path, handler)}).Add(1)

			return nil, err
		}

		totalRequestsCounter.With(prometheus.Labels{"handler": handler, "method": r.Method,
			"code": resp.Status, "endpoint": utils.RemoveIDFrom(r.URL.Path, handler)}).Add(1)

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithFields(log.Fields{
				"method":   r.Method,
				"endpoint": r.URL.Path,
			}).Debugf("Failed to read response body")
		}

		if (resp.StatusCode >= 400) && (resp.StatusCode != 401) {
			log.WithFields(log.Fields{
				"status":   resp.Status,
				"method":   r.Method,
				"endpoint": r.URL.Path,
				"body":     string(body),
			}).Warnf("Request did not succeed")

			requestsFailureCounter.With(prometheus.Labels{"handler": handler, "method": r.Method,
				"code": resp.Status, "endpoint": utils.RemoveIDFrom(r.URL.Path, handler)}).Add(1)
			return nil, fmt.Errorf("Request failed")
		}
		return resp, nil
	})
}

func instrumentRoundTripperDuration(obs prometheus.ObserverVec, next http.RoundTripper, handler string) promhttp.RoundTripperFunc {
	return promhttp.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
		start := time.Now()
		resp, err := next.RoundTrip(r)
		if err == nil {
			obs.With(prometheus.Labels{"method": r.Method, "code": resp.Status, "endpoint": utils.RemoveIDFrom(r.URL.Path, handler)}).Observe(time.Since(start).Seconds())
		}
		return resp, err
	})
}
