package ctxext

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceIDFromCtx(ctx context.Context) string {
	span := trace.SpanFromContext(ctx).SpanContext()
	if span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}
func GetSpanIDFromCtx(ctx context.Context) string {
	span := trace.SpanFromContext(ctx).SpanContext()
	if span.HasSpanID() {
		return span.SpanID().String()
	}
	return ""
}
func GetTraceIAndSpanIDFromCtx(ctx context.Context) (string, string) {
	span := trace.SpanFromContext(ctx).SpanContext()
	traceID := ""
	spanID := ""
	if span.HasSpanID() {
		spanID = span.SpanID().String()
	}
	if span.HasTraceID() {
		traceID = span.TraceID().String()
	}
	return traceID, spanID
}
func GetNewContextWithTraceIDAndSpanIDStr(ctx context.Context, traceID string, spanID string) context.Context {

	traceIDObj, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return ctx
	}
	spanIDObj, err := trace.SpanIDFromHex(spanID)
	if err != nil {
		return ctx
	}
	return GetNewContextWithTraceIDAndSpanID(ctx, traceIDObj, spanIDObj)
}
func GetNewContextWithTraceIDAndSpanID(ctx context.Context, traceID trace.TraceID, spanID trace.SpanID) context.Context {

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		TraceState: trace.TraceState{},
		Remote:     false,
	})
	return trace.ContextWithSpanContext(ctx, spanCtx)
}
