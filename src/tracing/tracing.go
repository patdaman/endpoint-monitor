package tracing

import (
	"context"
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
)

// Tracer instance
var tracer opentracing.Tracer

// SetTracer can be used by unit tests to provide a NoopTracer instance. Real users should always
// use the InitTracing func.
func SetTracer(initializedTracer opentracing.Tracer) {
	tracer = initializedTracer
}

// InitTracing connects the calling service to Zipkin and initializes the tracer.
func InitTracing(zipkinURL string, serviceName string) {
	logrus.Infof("Connecting to zipkin server at %v", zipkinURL)
	// collector, err := zipkin.NewHTTPCollector(
	//	fmt.Sprintf("%s/api/v1/spans", zipkinURL))
	// collector, err := zipkinhttp.NewHTTPCollector(
	//	fmt.Sprintf("%s/api/v1/spans", zipkinURL))

	reporter := zipkinhttp.NewReporter("http://zipkinhost:9411/api/v2/spans")
	defer reporter.Close()
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("myService", "myservice.mydomain.com:80")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	if err != nil {
		logrus.Info("Error connecting to zipkin server at " +
			fmt.Sprintf("%s/api/v1/spans", zipkinURL) + ". Error: " + err.Error())
		logrus.Errorln("Error connecting to zipkin server at " +
			fmt.Sprintf("%s/api/v1/spans", zipkinURL) + ". Error: " + err.Error())
		panic("Error connecting to zipkin server at " +
			fmt.Sprintf("%s/api/v1/spans", zipkinURL) + ". Error: " + err.Error())
	}
	// tracer, err = zipkin.NewTracer(
	// 	zipkin.NewRecorder(collector, false, "127.0.0.1:0", serviceName))

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Errorln("Error starting new zipkin tracer. Error: " + err.Error())
		panic("Error starting new zipkin tracer. Error: " + err.Error())
	}
	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)
	logrus.Infof("Successfully started zipkin tracer for service '%v'", serviceName)
}

// StartHTTPTrace loads tracing information from an INCOMING HTTP request.
func StartHTTPTrace(r *http.Request, opName string) opentracing.Span {
	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		return tracer.StartSpan(
			opName, ext.RPCServerOption(clientContext))
	} else {
		return tracer.StartSpan(opName)
	}
}

// MapToCarrier converts a generic map to opentracing http headers carrier
func MapToCarrier(headers map[string]interface{}) opentracing.HTTPHeadersCarrier {
	carrier := make(opentracing.HTTPHeadersCarrier)
	for k, v := range headers {
		// delivery.Headers
		carrier.Set(k, v.(string))
	}
	return carrier
}

// CarrierToMap converts a TextMapCarrier to the amqp headers format
func CarrierToMap(values map[string]string) map[string]interface{} {
	headers := make(map[string]interface{})
	for k, v := range values {
		headers[k] = v
	}
	return headers
}

// StartTraceFromCarrier extracts tracing info from a generic map and starts a new span.
func StartTraceFromCarrier(carrier map[string]interface{}, spanName string) opentracing.Span {

	clientContext, err := tracer.Extract(opentracing.HTTPHeaders, MapToCarrier(carrier))
	var span opentracing.Span
	if err == nil {
		span = tracer.StartSpan(
			spanName, ext.RPCServerOption(clientContext))
	} else {
		span = tracer.StartSpan(spanName)
	}
	return span
}

// AddTracingToReq adds tracing information to an OUTGOING HTTP request
func AddTracingToReq(req *http.Request, span opentracing.Span) {
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err := tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier)
	if err != nil {
		panic("Unable to inject tracing context: " + err.Error())
	}
}

// AddTracingToReqFromContext adds tracing information to an OUTGOING HTTP request
func AddTracingToReqFromContext(ctx context.Context, req *http.Request) {
	if ctx.Value("opentracing-span") == nil {
		return
	}
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err := tracer.Inject(
		ctx.Value("opentracing-span").(opentracing.Span).Context(),
		opentracing.HTTPHeaders,
		carrier)
	if err != nil {
		panic("Unable to inject tracing context: " + err.Error())
	}
}

func AddTracingToTextMapCarrier(span opentracing.Span, val opentracing.TextMapCarrier) error {
	return tracer.Inject(span.Context(), opentracing.TextMap, val)
}

// StartSpanFromContext starts a span.
func StartSpanFromContext(ctx context.Context, opName string) opentracing.Span {
	span := ctx.Value("opentracing-span").(opentracing.Span)
	child := tracer.StartSpan(opName, ext.RPCServerOption(span.Context()))
	return child
}

// StartChildSpanFromContext starts a child span from span within the supplied context, if available.
func StartChildSpanFromContext(ctx context.Context, opName string) opentracing.Span {
	if ctx.Value("opentracing-span") == nil {
		return tracer.StartSpan(opName, ext.RPCServerOption(nil))
	}
	parent := ctx.Value("opentracing-span").(opentracing.Span)
	child := tracer.StartSpan(opName, opentracing.ChildOf(parent.Context()))
	return child
}

// StartSpanFromContextWithLogEvent starts span from context with logevent
func StartSpanFromContextWithLogEvent(ctx context.Context, opName string, logStatement string) opentracing.Span {
	span := ctx.Value("opentracing-span").(opentracing.Span)
	child := tracer.StartSpan(opName, ext.RPCServerOption(span.Context()))
	child.LogEvent(logStatement)
	return child
}

// CloseSpan logs event finishes span.
func CloseSpan(span opentracing.Span, event string) {
	span.LogEvent(event)
	span.Finish()
}

// LogEventToOngoingSpan extracts span from context and adds LogEvent
func LogEventToOngoingSpan(ctx context.Context, logMessage string) {
	if ctx.Value("opentracing-span") != nil {
		ctx.Value("opentracing-span").(opentracing.Span).LogEvent(logMessage)
	}
}

// UpdateContext updates the supplied context with the supplied span.
func UpdateContext(ctx context.Context, span opentracing.Span) context.Context {
	return context.WithValue(ctx, "opentracing-span", span)
}
