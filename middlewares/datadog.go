package middlewares

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/dogstatsd"
	kitlog "github.com/go-kit/kit/log"
	"github.com/containous/traefik/log"
	"net/http"
)

var _ Metrics = (Metrics)(nil)

var datadogClient = dogstatsd.New("traefik.", kitlog.LoggerFunc(func(keyvals ...interface{}) error {
	log.Info(keyvals)
	return nil
}))

// DataDog is an Implementation for Metrics that exposes datadog metrics for the latency
// and the number of requests partitioned by status code and method.
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

func initDatadogClient(address string, pushInterval time.Duration) *time.Ticker {
	report := time.NewTicker(pushInterval)

	go datadogClient.SendLoop(report.C, "udp", address)

	return report
}

func (dd *Datadog) Stop(report *time.Ticker) {
	report.Stop()
}