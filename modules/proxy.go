package modules

import (
	"github.com/linsheng9731/slb/logger"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var lg = logger.Server
var access = logger.Access

// Status writer wrap http response writer
// to record http status code
type StatusRepWriter struct {
	status int
	w      http.ResponseWriter
}

func (w *StatusRepWriter) Flush() {
	if wf, ok := w.w.(http.Flusher); ok {
		wf.Flush()
	}
}
func (w *StatusRepWriter) Header() http.Header         { return w.w.Header() }
func (w *StatusRepWriter) Write(d []byte) (int, error) { return w.w.Write(d) }
func (w *StatusRepWriter) WriteHeader(code int) {
	w.status = code
	w.w.WriteHeader(code)
}

type HttpProxy struct {
	Transport http.RoundTripper
	Lookup    func(r *http.Request) *Route
}

func (p *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	route := p.Lookup(r)
	if route == nil {
		w.WriteHeader(404)
		return
	}
	statusW := &StatusRepWriter{-1, w}
	schema, host := schemaHost(route.Dst)
	targetURL := r.URL
	targetURL.Host = host
	targetURL.Scheme = schema

	upgrade, accept := r.Header.Get("Upgrade"), r.Header.Get("Accept")
	var httpProxy http.Handler
	switch {
	case upgrade == "websocket" || upgrade == "Websocket":
		r.URL = targetURL
		if targetURL.Scheme == "https" || targetURL.Scheme == "wss" {
			lg.Error("https and wss are not implement yet!")
			httpProxy = newRawHTTPProxy(targetURL.Host, net.Dial)
		} else {
			httpProxy = newRawHTTPProxy(targetURL.Host, net.Dial)
		}
	case accept == "text/event-stream":
		httpProxy = newHTTPProxy(targetURL, p.Transport, time.Duration(10))
	default:
		httpProxy = newHTTPProxy(targetURL, p.Transport, time.Duration(0))
	}
	httpProxy.ServeHTTP(statusW, r)
	timeEnd := time.Now()
	dur := timeEnd.Sub(timeStart)
	reqUrl := r.URL
	reqUrl.Host = r.Host
	access.Info(r.Method, "From", reqUrl, "To", route.Dst, dur.Seconds()*1000, "ms", statusW.status)
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

type dialFunc func(net, add string) (net.Conn, error)

func newRawHTTPProxy(host string, dial dialFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "not a hijacker", http.StatusInternalServerError)
			return
		}

		in, _, err := hj.Hijack()
		if err != nil {
			lg.Error(" Hijack error for %s. %s", r.URL, err)
			http.Error(w, "hijack error", http.StatusInternalServerError)
			return
		}

		defer in.Close()

		out, err := dial("tcp", host)

		if err != nil {
			lg.Error(" WS error for %s. %s", r.URL, err)
			http.Error(w, "error contacting backend server", http.StatusInternalServerError)
			return
		}
		err = r.Write(out)
		if err != nil {
			lg.Error(" WS error for %s. %s", r.URL, err)
			http.Error(w, "error contacting backend server", http.StatusInternalServerError)
			return
		}

		errCh := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err := io.Copy(dst, src)
			errCh <- err
		}
		go cp(out, in)
		go cp(in, out)
		err = <-errCh
		if err != nil && err != io.EOF {
			lg.Info(" WS error for %s. %s", r.URL, err)
		}
	})
}
