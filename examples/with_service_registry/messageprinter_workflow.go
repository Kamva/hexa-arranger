package main

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"time"
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
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var result string
	f := workflow.ExecuteActivity(ctx, PrintMessageActivity, message)
	if err := f.Get(ctx, &result); err != nil {
		logger.Error("Error in execution of activity", zap.Error(err))
		return err
	}

	logger.Info("workflow completed :)")
	return nil
}

// PrintMessageActivity is the activity to print a message.
func PrintMessageActivity(ctx context.Context, message Message) (string, error) {
	s := fmt.Sprintf("msg: %s", message.Msg)
	fmt.Println(s)

	return s, nil
}
