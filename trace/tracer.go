package trace

import (
	"fmt"
	"io"
)

// Tracer is an interface that describes an object capable of tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

type NilTracer struct {}

func (t *NilTracer) Trace(a ...interface{}) {}

func Off() Tracer {
	return &NilTracer{}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}
