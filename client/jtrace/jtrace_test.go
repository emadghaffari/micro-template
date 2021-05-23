package jtrace

import (
	"context"
	"fmt"
	"micro/config"
	"sync"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {

	tests := []struct {
		step string
		conf config.Config
		err  error
	}{
		{
			step: "A",
			conf: config.Config{
				Jaeger: config.Jaeger{},
			},
			err: fmt.Errorf("no service name provided"),
		},
		{
			step: "B",
			conf: config.Config{
				Jaeger: config.Jaeger{
					HostPort: "127.0.0.1",
					LogSpans: false,
				},
			},
			err: fmt.Errorf("no service name provided"),
		},
		{
			step: "C",
			conf: config.Config{
				Service: config.Service{
					Name: "TEST",
				},
				Jaeger: config.Jaeger{
					HostPort: "127.0.0.1",
					LogSpans: false,
				},
			},
			err: fmt.Errorf("address 127.0.0.1: missing port in address"),
		},
		{
			step: "D",
			conf: config.Config{
				Service: config.Service{
					Name: "TEST",
				},
				Jaeger: config.Jaeger{
					HostPort: "127.0.0.1:6831",
					LogSpans: false,
				},
			},
			err: fmt.Errorf("address 127.0.0.1: missing port in address"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			once = sync.Once{}
			_, err := Tracer.Connect(tc.conf)
			if err != nil {
				assert.Equal(t, tc.err.Error(), err.Error())
				return
			}
			// closer.Close()

		})
	}
}

func TestGetTracer(t *testing.T) {
	Tracer.GetTracer()
}

func TestFromContext(t *testing.T) {
	tracer := Tracer.GetTracer()
	span := tracer.StartSpan("testing")

	// context with span
	// create a child span
	Tracer.FromContext(opentracing.ContextWithSpan(context.Background(), span), "testing-2")

	// empty context
	// create new span
	Tracer.FromContext(context.Background(), "testing-3")
}

func TestStartSpan(t *testing.T) {
	span := Tracer.StartSpan("test")
	span.SetTag("tag", "value")
	span.Finish()

}

func TestContextWithSpan(t *testing.T) {
	span := Tracer.StartSpan("test")

	// should not use built-in type string as key for value; define your own type to avoid collisions
	ctx := context.WithValue(context.Background(), Span, span)
	Tracer.ContextWithSpan(ctx, span)

	Tracer.ContextWithSpan(context.Background(), span)
}

func TestSpanFromContext(t *testing.T) {

	span, _ := Tracer.SpanFromContext(context.Background(), "test")
	span.SetTag("tag", "value")
	span.Finish()
}

func TestChildOf(t *testing.T) {
	span, _ := Tracer.SpanFromContext(context.Background(), "test")
	span.SetTag("tag", "value")
	span.Finish()

	child := Tracer.ChildOf(span, "child-test")
	child.SetTag("tag", "value")
	child.Finish()
}
