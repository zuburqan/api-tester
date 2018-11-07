package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/zuburqan/api-tester/config"
	"github.com/zuburqan/api-tester/digest"
	"github.com/zuburqan/api-tester/metrics"
	"github.com/zuburqan/api-tester/transport"
	"github.com/zuburqan/api-tester/utils"
)

var (
	configFile           = flag.String("c", "./etc/api-tester.conf", "Path API config file")
	connectionConfigFile = flag.String("nc", "./etc/api-tester-connection.conf", "Path to connection config file")
	cfg                  *config.APIConfig
	connectionCfg        *config.ConnectionConfig
)

func main() {
	loadConfig()
	startMetricsServer()
	m := metrics.New()
	client := buildHTTPClient(m)
	spawnApiUsers(client, m)
}

func loadConfig() {
	flag.Parse()
	cfg = config.Load(*configFile)
	connectionCfg = config.LoadConnection(*connectionConfigFile)
	setupLogging()
}

func setupLogging() {
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal("Failed to parse log level: ", err.Error())
	}
	log.SetLevel(logLevel)
}

func startMetricsServer() {
	metrics.StartMetricsServer(cfg.StatsHost, cfg.StatsPort)
}

func buildHTTPClient(m *metrics.Metrics) http.Client {
	transport := transport.SetupRoundTrippers(m, http.DefaultTransport, cfg.Destination)
	return http.Client{Transport: transport, Timeout: time.Duration(cfg.ClientTimeout) * time.Second}
}

func spawnApiUsers(client http.Client, m *metrics.Metrics) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}

	// spawn goroutines
	for id := 1; id <= cfg.Users; id++ {
		wg.Add(1)
		go executeJourneySet(ctx, &wg, client, m, id)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// tell the goroutines to stop and wait for them to finish
	log.Info("Shutdown request received - waiting for any in-progress journeys to finish...")
	cancel()
	wg.Wait()
	log.Info("Stopped")
}

func executeJourneySet(ctx context.Context, wg *sync.WaitGroup, client http.Client, m *metrics.Metrics, id int) {
	defer wg.Done()
	interval := time.Duration(cfg.Sleep) * time.Second
	timer := time.NewTimer(interval)

	for {
		doJourneys(ctx, client, m, utils.GenerateID(id, cfg.Destination))
		timer.Reset(interval)
		select {
		case <-timer.C:
		case <-ctx.Done():
			return
		}
	}
}

func doJourneys(ctx context.Context, client http.Client, m *metrics.Metrics, ID string) {
	obs := m.JourneyDuration.MustCurryWith(prometheus.Labels{"handler": cfg.Destination})

	for _, journey := range cfg.Journeys {
		if ctx.Err() != nil {
			return
		}

		m.TotalJourneyCounter.With(prometheus.Labels{"handler": cfg.Destination, "journey": journey.Name}).Add(1)

		err := runRequests(journey.Setup, client, "SETUP: "+journey.Name, ID)
		if err != nil {
			log.WithFields(log.Fields{
				"journey": journey.Name,
			}).Warnf("Journey setup failed. Moving to next journey...")
			m.JourneysFailureCounter.With(prometheus.Labels{"handler": cfg.Destination, "journey": journey.Name}).Add(1)
			continue
		}

		start := time.Now()
		err = runRequests(journey.Requests, client, journey.Name, ID)
		if err != nil {
			log.WithFields(log.Fields{
				"journey": journey.Name,
			}).Warnf("Journey requests failed. Moving to next journey...")
			m.JourneysFailureCounter.With(prometheus.Labels{"handler": cfg.Destination, "journey": journey.Name}).Add(1)
			continue
		}
		duration := time.Since(start).Seconds()
		obs.With(prometheus.Labels{"journey": journey.Name}).Observe(duration)

		log.WithFields(log.Fields{
			"journey":  journey.Name,
			"duration": duration,
		}).Infof("Journey completed")

		err = runRequests(journey.Cleanup, client, "CLEAN-UP: "+journey.Name, ID)

		if err != nil {
			log.WithFields(log.Fields{
				"journey": journey.Name,
			}).Warnf("Journey cleanup failed")
		}
	}
}

func runRequests(requests []config.Request, client http.Client, journeyName, ID string) error {
	var resp *http.Response
	var req *http.Request
	var err error

	for _, request := range requests {
		endpoint := utils.Inject(ID, request.Endpoint)
		payload := utils.Inject(ID, request.Payload)

		if cfg.Auth == "digest" {
			resp, err = digest.Digest(connectionCfg, endpoint, request.Method, []byte(payload), &client)
		} else {
			req, err = http.NewRequest(request.Method, connectionCfg.Host+endpoint, bytes.NewBuffer([]byte(payload)))
			if err != nil {
				log.WithFields(log.Fields{
					"method":   req.Method,
					"endpoint": endpoint,
				}).Warnf("Failed to create request")

				return err
			}
			if cfg.Auth == "basic" {
				req.SetBasicAuth(connectionCfg.Username, connectionCfg.Password)
			}
			resp, err = client.Do(req)
		}
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Warnf("Request failed")

			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.WithFields(log.Fields{
				"method":   request.Method,
				"endpoint": endpoint,
				"error":    err,
			}).Debugf("Failed to read response body")
		}

		log.WithFields(log.Fields{
			"journey":  journeyName,
			"method":   request.Method,
			"endpoint": endpoint,
		}).Infof("Request completed")

		log.WithFields(log.Fields{
			"status": resp.Status,
			"body":   string(body),
		}).Debugf("Response status & body")
	}
	return nil
}
