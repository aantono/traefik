package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/containous/traefik/types"
	"github.com/stvp/go-udp-testing"
)

func TestDatadog(t *testing.T) {
	udp.SetAddr(":18125")
	// This is needed to make sure that UDP Listener listens for data a bit longer, otherwise it will quit after a millisecond
	udp.Timeout = 5 * time.Second
	recorder := httptest.NewRecorder()
	InitDatadogClient(&types.Datadog{":18125", "1s"})

	n := negroni.New()
	dd := NewDataDog("test")
	defer dd.Stop()
	metricsMiddlewareBackend := NewMetricsWrapper(dd)

	n.Use(metricsMiddlewareBackend)
	r := http.NewServeMux()
	//r.Handle("/metrics", promhttp.Handler())
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
		"traefik.traefik_requests_total:1.000000|c|#service:test,code:404,method:GET\n",
		"traefik.traefik_requests_total:1.000000|c|#service:test,code:200,method:GET\n",
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
