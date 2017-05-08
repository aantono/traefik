package middlewares

import (
	"fmt"
	"time"

	"github.com/containous/traefik/log"
	"github.com/containous/traefik/types"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/dogstatsd"
	"net/http"
)

var _ Metrics = (Metrics)(nil)

var datadogClient = dogstatsd.New("traefik.", kitlog.LoggerFunc(func(keyvals ...interface{}) error {
	log.Info(keyvals)
	return nil
}))
var datadogTicker *time.Ticker

// Datadog is an Implementation for Metrics that exposes datadog metrics for the latency and the number of requests partitioned by status code and method.
type Datadog struct {
	reqsCounter      metrics.Counter
	latencyHistogram metrics.Histogram
}

func (dd *Datadog) getReqsCounter() metrics.Counter {
	return dd.reqsCounter
}

func (dd *Datadog) getLatencyHistogram() metrics.Histogram {
	return dd.latencyHistogram
}

// NewDataDog creates new instance of Datadog
func NewDataDog(name string) *Datadog {
	var m Datadog

	m.reqsCounter = datadogClient.NewCounter(metricsReqsName, 1.0).With("service", name)
	m.latencyHistogram = datadogClient.NewHistogram(metricsLatencyName, 1.0).With("service", name)

	return &m
}

func (dd *Datadog) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// DataDog is a push-metrics reporter, so nothing to pull
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
}

// InitDatadogClient initializes metrics pusher and creates a datadogClient if not created already
func InitDatadogClient(config *types.Datadog) *time.Ticker {
	if datadogTicker == nil {
		address := config.Address
		if len(address) == 0 {
			address = "localhost:8125"
		}
		pushInterval, err := time.ParseDuration(config.PushInterval)
		if err != nil {
			pushInterval = 10 * time.Second
		}

		report := time.NewTicker(pushInterval)

		go datadogClient.SendLoop(report.C, "udp", address)

		datadogTicker = report
	}
	return datadogTicker
}

// Stop stops internal datadogTicker which controls the pushing of metrics to DD Agent and resets it to `nil`
func (dd *Datadog) Stop() {
	datadogTicker.Stop()
	datadogTicker = nil
}
