package Service

import (
	"net/http"
	"net/http/httputil"
)

type Service interface {
	// Address returns the address with which to access the server
	Address() string

	// IsAlive returns true if the server is alive and able to serve requests
	IsAlive() bool

	// Serve uses this server to process the request
	Service(rw http.ResponseWriter, req *http.Request)
}

type SimpleService struct {
	Addr  string
	Proxy *httputil.ReverseProxy
}

func (s *SimpleService) Address() string {
	return s.Addr
}

func (s *SimpleService) IsAlive() bool {

	return true
}

func (s *SimpleService) Service(rw http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(rw, req)
}
