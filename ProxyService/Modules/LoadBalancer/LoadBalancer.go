package LoadBalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/Pmca96/API-Servers/ProxyService/Modules/Service"
)

type LoadBalancer struct {
	Port            string
	RoundRobinCount int
	Services        []Service.Service
}

func NewSimpleServer(addr string) *Service.SimpleService {
	serviceUrl, err := url.Parse(addr)
	handleErr(err)

	return &Service.SimpleService{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serviceUrl),
	}
}

func NewLoadBalancer(port string, services []Service.Service) *LoadBalancer {
	return &LoadBalancer{
		Port:            port,
		RoundRobinCount: 0,
		Services:        services,
	}
}

// handleErr prints the error and exits the program
// Note: this is not how one would want to handle errors in production, but
// serves well for demonstration purposes.
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

// getNextServerAddr returns the address of the next available server to send a
// request to, using a simple round-robin algorithm
func (lb *LoadBalancer) getNextAvailableServer() Service.Service {
	server := lb.Services[lb.RoundRobinCount%len(lb.Services)]
	for !server.IsAlive() {
		lb.RoundRobinCount++
		server = lb.Services[lb.RoundRobinCount%len(lb.Services)]
	}
	lb.RoundRobinCount++

	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	targetService := lb.getNextAvailableServer()

	// could optionally log stuff about the request here!
	fmt.Printf("forwarding request to address %q\n", targetService.Address())

	// could delete pre-existing X-Forwarded-For header to prevent IP spoofing
	targetService.Service(rw, req)
}
