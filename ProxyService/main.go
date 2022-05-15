package main

import (
	"fmt"
	"net/http"

	"github.com/Pmca96/API-Servers/ProxyService/Modules/LoadBalancer"
	"github.com/Pmca96/API-Servers/ProxyService/Modules/Service"
)

func main() {

	// service registration
	userServices := []Service.Service{
		LoadBalancer.NewSimpleServer("http://localhost:8101"),
		LoadBalancer.NewSimpleServer("http://localhost:8102"),
		LoadBalancer.NewSimpleServer("http://localhost:8103"),
	}
	messageServices := []Service.Service{
		LoadBalancer.NewSimpleServer("http://localhost:8111"),
		LoadBalancer.NewSimpleServer("http://localhost:8112"),
		LoadBalancer.NewSimpleServer("http://localhost:8113"),
	}

	// attatch all and declare service name by order
	allTypesService := Service.AllServices{
		Services: [][]Service.Service{
			userServices,
			messageServices,
		},
		Names: []string {
			"Users", 
			"Message",
		},
	}

	// init loadBalancer with all microServices
	loadB := LoadBalancer.NewLoadBalancer("8080", allTypesService)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadB.ServeProxy(rw, req)
	}

	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", loadB.Port)
	http.ListenAndServe(":"+loadB.Port, nil)
}
