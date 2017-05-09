package middlewares

import (
	"net/http"

	"github.com/containous/traefik/log"
)

// AddPrefix is a middleware used to add prefix to an URL request
type AddPrefix struct {
	Handler http.Handler
	Prefix  string
}

func (s *AddPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newPath := s.Prefix + r.URL.Path
	log.Debugf("Adding Prefix %s to %s = %s", s.Prefix, r.URL.Path, newPath)
	r.URL.Path = newPath
	r.RequestURI = r.URL.RequestURI()
	s.Handler.ServeHTTP(w, r)
}

// SetHandler sets handler
func (s *AddPrefix) SetHandler(Handler http.Handler) {
	s.Handler = Handler
}
