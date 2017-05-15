package middlewares

import (
	"fmt"
	"time"
	"net/http"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/statsd"
	kitlog "github.com/go-kit/kit/log"
	"github.com/containous/traefik/log"
)

var _ Metrics = (Metrics)(nil)

var statsdClient = statsd.New("traefik.", kitlog.LoggerFunc(func(keyvals ...interface{}) error {
	log.Info(keyvals)
	return nil
}))

// StatsD is an Implementation for Metrics that exposes statsd metrics for the latency
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

func (s *Statsd) Stop(report *time.Ticker) {
	report.Stop()
}