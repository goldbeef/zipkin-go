package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"log"
	"time"
)

func main() {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter("http://zipkinhost:9411/api/v2/spans")
	defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("myService", "myservice.mydomain.com:80")
	if err != nil {
	log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
	log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	// tracer can now be used to create spans.
	span := opentracing.StartSpan("some_operation")
	time.Sleep(time.Second)

	ctx1 := span.Context()
	span1, ctx2:= opentracing.StartSpanFromContext(ctx1, "some_operation1")
	time.Sleep(time.Second)
	span1.Finish()

	span2, _:= tracer.StartSpanFromContext(ctx2, "some_operation2")
	time.Sleep(time.Second)
	span2.Finish()

	span.Finish()

}