package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kamva/hexa"
	arranger "github.com/kamva/hexa-arranger"
	"github.com/kamva/tracer"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var SampleHexaErr = hexa.NewError(http.StatusInternalServerError, "arranger.examples.fake_internal_err")

func registerWorkflow(worker worker.Worker, h Handlers) {
	worker.RegisterWorkflow(h.PrintMessageWorkflow)
	worker.RegisterActivity(h.PrintMessageActivity)
}

type Message struct {
	Msg string `json:"msg"`
}

type handlers struct {
}

func NewHandlers() Handlers {
	return &handlers{}
}

func (h *handlers) PrintMessageActivity(ctx context.Context, message Message) (string, error) {
	s := fmt.Sprintf("msg: %s", message.Msg)
	fmt.Println(s)

	fmt.Println("Now I want to return an Example error named: hi_err")

	return "", SampleHexaErr.SetError(tracer.Trace(errors.New("hey, this is just a sample error :)")))
	//return "", tracer.Trace(errors.New("hey, what happens (:"))
	//return "", arranger.HexaToApplicationErr(SampleHexaErr, t)

	// When temporal implements activity interceptor, we use workflow and
	// activity interceptor to recover that errors also, for now we let
	// temporal itself handle panics.
	// panic(errors.New("fake panic to see what happen :)"))
}

func (h *handlers) PrintMessageWorkflow(ctx workflow.Context, message Message) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var result string
	f := workflow.ExecuteActivity(ctx, h.PrintMessageActivity, message)
	if err := f.Get(ctx, &result); err != nil {
		if err, ok := arranger.HexaErrFromApplicationErrWithOk(err); ok {
			hErr := err.(hexa.Error)
			fmt.Println("hexa error in activity result: ", hErr.ID(), hErr.HTTPStatus())
		} else {
			fmt.Println("returned activity error is not a hexa error (:")
		}

		return err
	}

	logger.Info("workflow completed :)")
	return nil
}

var _ Handlers = &handlers{}
