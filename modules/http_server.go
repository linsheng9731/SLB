package modules

import (
	"context"
	"crypto/tls"
	"fmt"
	proxyproto "github.com/armon/go-proxyproto"
	"net"
	"net/http"
	"sync"
	"time"
)

type HttpServer interface {
	Close() error
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
}

var (
	// mu guards servers which contains the list
	// of running proxy servers.
	mu      sync.Mutex
	servers []HttpServer
)

type Listen struct {
	Addr         string
	Proto        string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func ListenAndServeHTTP(l Listen, h http.Handler) error {
	ln, err := ListenTCP(l.Addr, nil)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:         l.Addr,
		Handler:      h,
		ReadTimeout:  l.ReadTimeout,
		WriteTimeout: l.WriteTimeout,
	}
	return serve(ln, srv)
}

func serve(ln net.Listener, srv HttpServer) error {
	mu.Lock()
	servers = append(servers, srv)
	mu.Unlock()
	return srv.Serve(ln)
}

func ListenTCP(laddr string, cfg *tls.Config) (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp", laddr)
	if err != nil {
		return nil, fmt.Errorf("listen: Fail to resolve tcp addr. %s", laddr)
	}

	var ln net.Listener
	ln, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listen: Fail to listen. %s", err)
	}

	// enable TCPKeepAlive support
	ln = tcpKeepAliveListener{ln.(*net.TCPListener)}

	// enable PROXY protocol support
	ln = &proxyproto.Listener{Listener: ln}

	// enable TLS
	if cfg != nil {
		ln = tls.NewListener(ln, cfg)
	}

	return ln, nil
}

// copied from http://golang.org/src/net/http/server.go?s=54604:54695#L1967
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}
