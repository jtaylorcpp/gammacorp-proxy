package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	zipkin "github.com/openzipkin/zipkin-go"
	middleware "github.com/openzipkin/zipkin-go/middleware/http"
	reporter "github.com/openzipkin/zipkin-go/reporter/http"
)

type ReverseProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	tracer *zipkin.Tracer
}

func NewSimpleReverseProxy(target string) {
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	var zipkinAddr string
	zipkinAddr, ok := os.LookupEnv("ZIPKIN")
	if !ok {
		zipkinAddr = "http://zipkin:9411/api/v2/spans"
	}

	fmt.Println("ataching to zipkin at: ", zipkinAddr)

	zipkinReporter := reporter.NewReporter(zipkinAddr)
	zipkingEndpoint, err := zipkin.NewEndpoint(
		"simple-reverseproxy",
		"127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	ZipkinSampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		panic(err)
	}

	tracer, err := zipkin.NewTracer(
		zipkinReporter,
		zipkin.WithSampler(ZipkinSampler),
		zipkin.WithLocalEndpoint(zipkingEndpoint),
		zipkin.WithSharedSpans(true),
		zipkin.WithTraceID128Bit(true),
	)
	if err != nil {
		panic(err)
	}

	zipkinMiddleware := middleware.NewServerMiddleware(
		tracer,
		middleware.SpanName("simple-reverseproxy"),
		middleware.TagResponseSize(true),
	)

	rp := &ReverseProxy{
		target: remote,
		proxy:  httputil.NewSingleHostReverseProxy(remote),
		tracer: tracer,
	}

	rpTransport, err := middleware.NewTransport(
		tracer,
		middleware.TransportTrace(true),
	)
	if err != nil {
		panic(err)
	}

	rp.proxy.Transport = rpTransport

	r := mux.NewRouter()
	r.HandleFunc("/", rp.proxy.ServeHTTP)
	r.Use(zipkinMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
