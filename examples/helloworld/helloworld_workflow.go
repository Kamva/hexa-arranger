package main

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func registerWorkflowsAndActivities(worker worker.Worker) {
	worker.RegisterWorkflow(HelloWorldWorkflow)
	worker.RegisterActivity(printHelloActivity)
}

// HelloWorldWorkflow is helloWorld workflow.
func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 20 * time.Minute,
		StartToCloseTimeout:    20 * time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)

	var result string
	f := workflow.ExecuteActivity(ctx, printHelloActivity, name)
	if err := f.Get(ctx, &result); err != nil {
		logger.Error("Error in execution of hello world activity", "err", err)
		return "", err
	}

	logger.Info("helloworld workflow completed :)")
	return result, nil
}

// printHelloActivity is the activity to print hello message.
func printHelloActivity(ctx context.Context, name string) (string, error) {
	s := fmt.Sprintf("Hello %s", name)
	fmt.Println(s)

	return s, nil
}
