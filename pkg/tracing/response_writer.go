package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type traceResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	span        trace.Span
	headerAdded bool
}

func (w *traceResponseWriter) addTraceIDHeader() {
	if !w.headerAdded {
		traceID := w.span.SpanContext().TraceID().String()
		if traceID != "" {
			w.ResponseWriter.Header().Set(HTTPTraceIDHeader, traceID)
		}
		w.headerAdded = true
	}
}

func (w *traceResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.addTraceIDHeader()
	w.ResponseWriter.WriteHeader(code)
}

func (w *traceResponseWriter) Write(b []byte) (int, error) {
	w.addTraceIDHeader()
	return w.ResponseWriter.Write(b)
}
