package modules

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type HttpProxy struct {
	Transport http.RoundTripper
	Lookup    func(r *http.Request) *Route
}

func (p *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := p.Lookup(r)
	if route == nil {
		w.WriteHeader(404)
		return
	}

	schema, host := schemaHost(route.Dst)
	targetURL := r.URL
	targetURL.Host = host
	targetURL.Scheme = schema

	upgrade, accept := r.Header.Get("Upgrade"), r.Header.Get("Accept")
	var httpProxy http.Handler
	switch {
	case upgrade == "websocket" || upgrade == "Websocket":
		log.Fatal("not impl websocket!")
		return
	case accept == "text/event-stream":
		httpProxy = newHTTPProxy(targetURL, p.Transport, time.Duration(10))
	default:
		httpProxy = newHTTPProxy(targetURL, p.Transport, time.Duration(0))
	}
	httpProxy.ServeHTTP(w, r)

}

func newHTTPProxy(target *url.URL, tr http.RoundTripper, flush time.Duration) http.Handler {
	return &httputil.ReverseProxy{
		// this is a simplified director function based on the
		// httputil.NewSingleHostReverseProxy() which does not
		// mangle the request and target URL since the target
		// URL is already in the correct format.
		Director: func(req *http.Request) {
			//req.URL.Scheme = target.Scheme
			//req.URL.Host = target.Host
			//req.URL.Path = target.Path
			//req.URL.RawQuery = target.RawQuery
			req.URL = target
			if _, ok := req.Header["User-Agent"]; !ok {
				// explicitly disable User-Agent so it's not set to default value
				req.Header.Set("User-Agent", "")
			}
		},
		FlushInterval: flush,
		Transport:     &transport{tr, nil, nil},
	}
}

// transport executes the roundtrip and captures the response. It is not
// safe for multiple or concurrent use since it only captures a single
// response.
type transport struct {
	http.RoundTripper
	resp *http.Response
	err  error
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.resp, t.err = t.RoundTripper.RoundTrip(r)
	return t.resp, t.err
}
