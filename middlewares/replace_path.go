package middlewares

import (
	"github.com/containous/traefik/log"
	"net/http"
)

// ReplacePath is a middleware used to replace the path of a URL request
type ReplacePath struct {
	Handler http.Handler
	Path    string
}

// ReplacedPathHeader is the default header to set the old path to
const ReplacedPathHeader = "X-Replaced-Path"

func (s *ReplacePath) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Add(ReplacedPathHeader, r.URL.Path)
	log.Debugf("Replacing Path %s with %s", r.URL.Path, s.Path)
	r.URL.Path = s.Path
	s.Handler.ServeHTTP(w, r)
}

// SetHandler sets handler
func (s *ReplacePath) SetHandler(Handler http.Handler) {
	s.Handler = Handler
}
