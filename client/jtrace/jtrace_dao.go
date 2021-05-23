package jtrace

import (
	"context"
	"io"
	"micro/config"
	"sync"

	"github.com/opentracing/opentracing-go"
)

var (
	tracer opentracing.Tracer
	Tracer itracer = &jtracer{}
	once   sync.Once
)

type itracer interface {
	Connect(config.Config) (io.Closer, error)
	GetTracer() opentracing.Tracer
	FromContext(ctx context.Context, startName string) opentracing.Span
	StartSpan(str string) opentracing.Span
	ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context
	SpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context)
	ChildOf(span opentracing.Span, name string) opentracing.Span
}

type jtracer struct{}
