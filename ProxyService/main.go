package main

import (
	"fmt"
	"net/http"

	"github.com/Pmca96/API-Servers/ProxyService/Modules/LoadBalancer"
	"github.com/Pmca96/API-Servers/ProxyService/Modules/Service"
)

func main() {
	userService := []Service.Service{
		LoadBalancer.NewSimpleServer("http://localhost:8091"),
		LoadBalancer.NewSimpleServer("http://localhost:8092"),
		LoadBalancer.NewSimpleServer("http://localhost:8093"),
	}

	loadB := LoadBalancer.NewLoadBalancer("8080", userService)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		loadB.ServeProxy(rw, req)
	}

	// register a proxy handler to handle all requests
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at 'localhost:%s'\n", loadB.Port)
	http.ListenAndServe(":"+loadB.Port, nil)
}
