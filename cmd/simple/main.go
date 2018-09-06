package main

import (
	proxy "github.com/jtaylorcpp/gammacorp-proxy"
)

func main() {
	proxy.NewSimpleReverseProxy("http://defaultbackend:8080")
}
