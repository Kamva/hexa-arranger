package main

import (
	"context"
	"fmt"
	"github.com/kamva/hexa-arranger/examples/protobufmessage/hello"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

func registerWorkflow(worker worker.Worker) {
	worker.RegisterWorkflow(PrintProtobufMessageWorkflow)
	worker.RegisterActivity(PrintMessageActivity)
}

// PrintProtobufMessageWorkflow is printer workflow
func PrintProtobufMessageWorkflow(ctx workflow.Context, message hello.Message) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)

	var result string
	f := workflow.ExecuteActivity(ctx, PrintMessageActivity, hello.ActivityMessage{Msg: message.Msg,})
	if err := f.Get(ctx, &result); err != nil {
		logger.Error("Error in execution of activity", zap.Error(err))
		return err
	}

	logger.Info("workflow completed :)")
	return nil
}

// PrintMessageActivity is the activity to print a message.
func PrintMessageActivity(ctx context.Context, message hello.ActivityMessage) (string, error) {
	s := fmt.Sprintf("msg: %s", message.Msg)
	fmt.Println(s)

	return s, nil
}
