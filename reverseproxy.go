package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/gorilla/mux"
)

func NewSimpleReverseProxy(target string) {
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	rproxy := httputil.NewSingleHostReverseProxy(remote)

	r := mux.NewRouter()
	r.HandleFunc("/", simplehandler(rproxy))

	http.Handle("/", r)
	http.ListenAndServe(":8080", r)
}

func simplehandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
