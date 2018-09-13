package main

import (
	"os"
	proxy "github.com/jtaylorcpp/gammacorp-proxy"
)

func main() {
	remote,ok := os.LookupEnv("SIMPLEBACKEND")
	if !ok {
		remote = "http://default-backend:8080"
	}
	proxy.NewSimpleReverseProxy(remote)
}
