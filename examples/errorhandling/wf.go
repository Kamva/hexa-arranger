package main

import (
	"context"
	"fmt"
	"time"

	"github.com/kamva/hexa"
	arranger "github.com/kamva/hexa-arranger"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func registerWorkflow(worker worker.Worker) {
	worker.RegisterWorkflow(PrintMessageWorkflow)
	worker.RegisterActivity(PrintMessageActivity)
}

type Message struct {
	Msg string `json:"msg"`
}

// PrintMessageWorkflow is printer workflow
func PrintMessageWorkflow(ctx workflow.Context, message Message) error {
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
	f := workflow.ExecuteActivity(ctx, PrintMessageActivity, message)
	if err := f.Get(ctx, &result); err != nil {
		logger.Error("Error in execution of activity",
			"err", err,
			"original_error_type", fmt.Sprintf("%T", err),
		)
		if err, ok := arranger.HexaErrFromApplicationErrWithOk(err); ok {
			hErr := err.(hexa.Error)
			fmt.Println("hexa error in activity result: ", hErr.ID(), hErr.HTTPStatus())
		} else {
			fmt.Println("this is not a hexa error (:")
		}

		return err
	}

	logger.Info("workflow completed :)")
	return nil
}

// PrintMessageActivity is the activity to print a message.
func PrintMessageActivity(ctx context.Context, message Message) (string, error) {
	s := fmt.Sprintf("msg: %s", message.Msg)
	fmt.Println(s)

	fmt.Println("Now I want to return an Example error named: hi_err")
	return "", arranger.HexaToApplicationErr(SampleHexaErr, t)
}
