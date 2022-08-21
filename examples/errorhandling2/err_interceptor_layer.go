package main

import (
	"context"

	"github.com/kamva/hexa"
	arranger "github.com/kamva/hexa-arranger"
	"github.com/kamva/hexa/hlog"
	"go.temporal.io/sdk/workflow"
)

type errInterceptor struct {
	h Handlers
	l hlog.Logger
	t hexa.Translator
}

func NewErrInterceptorLayer(h Handlers, l hlog.Logger, t hexa.Translator) Handlers {
	return &errInterceptor{
		h: h,
		l: l,
		t: t,
	}
}

func (e *errInterceptor) PrintMessageActivity(ctx context.Context, message Message) (string, error) {
	r1, err := e.h.PrintMessageActivity(ctx, message)
	return r1, arranger.HandleErr(err, e.l, e.t)
}

func (e *errInterceptor) PrintMessageWorkflow(ctx workflow.Context, message Message) error {
	return arranger.HandleErr(e.h.PrintMessageWorkflow(ctx, message), e.l, e.t)
}
