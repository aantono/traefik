package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stvp/go-udp-testing"
	"github.com/codegangsta/negroni"
)

func TestStatsD(t *testing.T) {
	udp.SetAddr(":18125")
	// This is needed to make sure that UDP Listener listens for data a bit longer, otherwise it will quit after a millisecond
	udp.Timeout = 5 * time.Second
	recorder := httptest.NewRecorder()
	ticker := initStatsdClient(":18125", 1 * time.Second)

	n := negroni.New()
	c := NewStatsD("test")
	defer c.Stop(ticker)
	metricsMiddlewareBackend := NewMetricsWrapper(c)

	n.Use(metricsMiddlewareBackend)
	r := http.NewServeMux()
	r.HandleFunc(`/ok`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	r.HandleFunc(`/not-found`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "not-found")
	})
	n.UseHandler(r)

	req1, err := http.NewRequest("GET", "http://localhost:3000/ok", nil)
	if err != nil {
		t.Error(err)
	}
	req2, err := http.NewRequest("GET", "http://localhost:3000/not-found", nil)
	if err != nil {
		t.Error(err)
	}

	expected := []string{
		// We are only validating counts, as it is nearly impossible to validate latency, since it varies every run
		"traefik.traefik_requests_total:2.000000|c\n",
	}

	udp.ShouldReceiveAll(t, expected, func() {
		n.ServeHTTP(recorder, req1)
		n.ServeHTTP(recorder, req2)
		//		body := recorder.Body.String()

		//if !strings.Contains(body, ddReqsName) {
		//	t.Errorf("body does not contain request total entry '%s'", ddReqsName)
		//}
		//if !strings.Contains(body, ddLatencyName) {
		//	t.Errorf("body does not contain request duration entry '%s'", ddLatencyName)
		//}
	})
}