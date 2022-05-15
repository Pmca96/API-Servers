package LoadBalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Pmca96/API-Servers/ProxyService/Modules/Service"
)

type LoadBalancer struct {
	Port             string
	RoundRobinCount  int
	AllServices Service.AllServices
}

func NewSimpleServer(addr string) *Service.SimpleService {
	serviceUrl, err := url.Parse(addr)
	handleErr(err)

	return &Service.SimpleService{
		Addr:  addr,
		Proxy: httputil.NewSingleHostReverseProxy(serviceUrl),
	}
}

// TODO PassRoundRobinCount to inside Services (Depend on the services organizations)
// NOT NEEDED if a single Server would contain all services one time
func NewLoadBalancer(port string, allServices Service.AllServices) *LoadBalancer {
	return &LoadBalancer{
		Port:             port,
		RoundRobinCount:  0,
		AllServices: allServices,
	}
}

// TODO: IMPLEMENT BETTER ERROR HANDLER
// handleErr prints the error and exits the program
// Note: this is not how one would want to handle errors in production, but
// serves well for demonstration purposes.
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

// Find MicroService to call and then procceed to the availableServer of that same microService
func (lb *LoadBalancer) findMicroServiceToCall(path string) Service.Service {
	
	server :=  lb.AllServices.Services[0][0]
	for i := 0; i < len(lb.AllServices.Names); i++ {
		// service would be like "Users"
		service:= "/"+lb.AllServices.Names[i];

		// Check if path starts with '/service/' or equals '/service'
		if (strings.HasPrefix(service+"/", path) || service == path  ) {
			server= lb.getNextAvailableServer(i)
			break
		}
    }
	
	return server
}


// getNextServerAddr returns the address of the next available server to send a
// request to, using a simple round-robin algorithm
func (lb *LoadBalancer) getNextAvailableServer(serviceIndex int) Service.Service {
	microService:=lb.AllServices.Services[serviceIndex]

	server :=  microService[lb.RoundRobinCount%len(microService)]
	for !server.IsAlive() {
		lb.RoundRobinCount++
		server = microService[lb.RoundRobinCount%len(microService)]
	}
	lb.RoundRobinCount++

	return server
}


func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	targetService := lb.findMicroServiceToCall(req.URL.Path)

	// could optionally log stuff about the request here!
	fmt.Printf("forwarding request to address %q\n", targetService.Address())

	// could delete pre-existing X-Forwarded-For header to prevent IP spoofing
	targetService.Service(rw, req)
}
