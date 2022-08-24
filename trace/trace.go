//@File     trace.go
//@Time     2022/08/24
//@Author   #Suyghur,

package trace

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"runtime"
)

func StartFuncSpan(ctx context.Context, name string, callback func(context.Context), kv ...attribute.KeyValue) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, name, oteltrace.WithSpanKind(oteltrace.SpanKindInternal), oteltrace.WithAttributes(kv...))
	defer span.End()
	callback(spanCtx)
}

func RunWithTrace(traceId, spanName string, callback func(ctx context.Context)) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	if spanName == "" {
		spanName = callerFuncName()
	}
	traceIDFromHex, _ := oteltrace.TraceIDFromHex(traceId)
	ctx := oteltrace.ContextWithSpanContext(context.Background(), oteltrace.NewSpanContext(oteltrace.SpanContextConfig{TraceID: traceIDFromHex}))
	spanCtx, span := tracer.Start(ctx, spanName)
	defer span.End()
	callback(spanCtx)
}

func callerFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := oteltrace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		traceId := spanCtx.TraceID().String()
		return traceId
	}
	return ""
}
