package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TraceIDHeader = "x-trace-id"
)

func UnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ctx = propagator.Extract(ctx, metadataCarrier(md))

		ctx, span := tracer.Start(
			ctx,
			info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		ctx = AddTraceIDToReponse(ctx)

		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}

func UnaryClientInterceptor(serviceName string) grpc.UnaryClientInterceptor {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		spanName := formatSpanName(ctx, method)

		ctx, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindClient),
		)
		defer span.End()

		carrier := metadataCarrier(extractOutgoingMetadata(ctx))

		propagator.Inject(ctx, carrier)

		ctx = metadata.NewOutgoingContext(ctx, metadata.MD(carrier))

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			trace.SpanFromContext(ctx).RecordError(err)
		}

		return err
	}
}

func formatSpanName(ctx context.Context, method string) string {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return "client." + method
	}

	return method
}

func extractOutgoingMetadata(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.New(nil)
	}

	return md.Copy()
}

func GetTraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

func AddTraceIDToReponse(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)
	if traceID == "" {
		return ctx
	}

	md := extractOutgoingMetadata(ctx)

	md.Set(TraceIDHeader, traceID)
	return metadata.NewOutgoingContext(ctx, md)
}
