package modules

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type SLBRequest struct {
	Header   http.Header
	Status   int
	Body     []byte
	Upgraded bool
	Backend  *Backend
}

type SLBRequestChan chan SLBRequest

func NewWorkerRequestErr(status int, body []byte) SLBRequest {
	return SLBRequest{
		Status: status,
		Body:   body,
	}
}

func NewWorkerRequest(status int, header http.Header, body []byte) SLBRequest {
	return SLBRequest{
		Status: status,
		Header: header,
		Body:   body,
	}
}

func NewWorkerRequestUpgraded() SLBRequest {
	return SLBRequest{
		Upgraded: true,
	}
}

func copy(dest *bufio.ReadWriter, src *bufio.ReadWriter) {
	buf := make([]byte, 40*1024)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			return
		}
		dest.Write(buf[0:n])
		dest.Flush()
	}
}

func copyBidir(frontendConn io.ReadWriteCloser, rwFront *bufio.ReadWriter,
	backendConn io.ReadWriteCloser, rwBack *bufio.ReadWriter) {

	finished := make(chan bool)

	go func() {
		copy(rwBack, rwFront)
		backendConn.Close()
		finished <- true
	}()

	go func() {
		copy(rwFront, rwBack)
		frontendConn.Close()
		finished <- true
	}()

	<-finished
	<-finished
}

func (s *SLBRequest) HijackWebSocket(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)

	if !ok {
		log.Println("Error: Webserver doesn't support hijacking")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	frontendConn, buffer, err := hj.Hijack()
	defer frontendConn.Close()

	URL := &url.URL{}
	UrlParsed, _ := URL.Parse(s.Backend.BackendConfig.Address)

	backendConn, err := net.Dial("tcp", UrlParsed.Host)
	if err != nil {
		log.Println("Error: Couldn't connect to backend server")
		http.Error(w, "Internal Error", http.StatusServiceUnavailable)
		return
	}
	defer backendConn.Close()

	err = r.Write(backendConn)
	if err != nil {
		log.Printf("Writing WebSocket request to backend server failed: %v", err)
		return
	}

	copyBidir(frontendConn, buffer, backendConn,
		bufio.NewReadWriter(bufio.NewReader(backendConn), bufio.NewWriter(backendConn)))
}
