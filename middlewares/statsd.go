package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/containous/traefik/log"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/statsd"
)

var _ Metrics = (Metrics)(nil)

var statsdClient = statsd.New("traefik.", kitlog.LoggerFunc(func(keyvals ...interface{}) error {
	log.Info(keyvals)
	return nil
}))

// Statsd is an Implementation for Metrics that exposes statsd metrics for the latency
// and the number of requests partitioned by status code and method.
type Statsd struct {
	reqsCounter      metrics.Counter
	latencyHistogram metrics.Histogram
}

func (s *Statsd) getReqsCounter() metrics.Counter {
	return s.reqsCounter
}

func (s *Statsd) getLatencyHistogram() metrics.Histogram {
	return s.latencyHistogram
}

// NewStatsD creates new instance of StatsD
func NewStatsD(name string) *Statsd {
	var m Statsd

	m.reqsCounter = statsdClient.NewCounter(metricsReqsName, 1.0).With("service", name)
	m.latencyHistogram = statsdClient.NewTiming(metricsLatencyName, 1.0).With("service", name)

	return &m
}

func (s *Statsd) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// StatsD is a push-metrics reporter, so nothing to pull
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
}

func initStatsdClient(address string, pushInterval time.Duration) *time.Ticker {
	report := time.NewTicker(pushInterval)

	go statsdClient.SendLoop(report.C, "udp", address)

	return report
}

// Stop stops internal datadogTicker which controls the pushing of metrics to DD Agent and resets it to `nil`
func (s *Statsd) Stop(report *time.Ticker) {
	report.Stop()
}
