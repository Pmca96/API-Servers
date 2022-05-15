package Service

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

// TODO remove log.fatalln because they stop program, proxy cant be stopped

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

type AllServices struct {
	Services   [][] Service 
	Names 		[] string
}

func (s *SimpleService) Address() string {
	return s.Addr
}


func (s *SimpleService) IsAlive() bool {
	resp, err := http.Get(s.Addr+"/ping")
	if err != nil {
		log.Fatalln(err)
		return false
	}
	
	body, err := io.ReadAll(resp.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
		return false
	}

    _ , err = strconv.ParseBool(string(body))
    return err == nil
}

func (s *SimpleService) Service(rw http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(rw, req)
}
