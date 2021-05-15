package jtrace

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
)

var (
	tracer opentracing.Tracer
	Tracer itracer = &jtracer{}
)

type itracer interface {
	Connect() (io.Closer, error)
	GetTracer() opentracing.Tracer
	FromContext(ctx context.Context, startName string) opentracing.Span
	StartSpan(str string) opentracing.Span
	ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context
	SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
	ChildOf(span opentracing.Span, name string) opentracing.Span
}

type jtracer struct{}
