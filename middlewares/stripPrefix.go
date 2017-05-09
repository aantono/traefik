package middlewares

import (
	"github.com/containous/traefik/log"
	"net/http"
	"strings"
)

const (
	forwardedPrefixHeader = "X-Forwarded-Prefix"
)

// StripPrefix is a middleware used to strip prefix from an URL request
type StripPrefix struct {
	Handler  http.Handler
	Prefixes []string
}

func (s *StripPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, prefix := range s.Prefixes {
		if p := strings.TrimPrefix(r.URL.Path, strings.TrimSpace(prefix)); len(p) < len(r.URL.Path) {
			log.Debugf("Striping Path Prefix %s from %s = %s", prefix, r.URL.Path, p)
			r.URL.Path = p
			r.Header[forwardedPrefixHeader] = []string{prefix}
			r.RequestURI = r.URL.RequestURI()
			s.Handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// SetHandler sets handler
func (s *StripPrefix) SetHandler(Handler http.Handler) {
	s.Handler = Handler
}
