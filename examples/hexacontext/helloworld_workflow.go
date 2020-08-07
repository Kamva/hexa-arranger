package main

import (
	"context"
	"fmt"
	arranger "github.com/kamva/hexa-arranger"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
	"time"
)

func registerWorkflowsAndActivities(worker worker.Worker) {
	worker.RegisterWorkflow(HelloWorldWorkflow)
	worker.RegisterActivity(printHelloActivity)
}

// HelloWorldWorkflow is helloWorld workflow.
func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
	hexaCtx := arranger.HexaCtxFromCadenceCtx(ctx)
	fmt.Printf("hexa correlation id in workflow: %s\n", hexaCtx.CorrelationID())
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    time.Minute,
		HeartbeatTimeout:       time.Second * 20,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)

	var result string
	f := workflow.ExecuteActivity(ctx, printHelloActivity, name)
	if err := f.Get(ctx, &result); err != nil {
		logger.Error("Error in execution of hello world activity", zap.Error(err))
		return "", err
	}

	logger.Info("helloworld workflow completed :)")
	return result, nil
}

// printHelloActivity is the activity to print hello message.
func printHelloActivity(ctx context.Context, name string) (string, error) {
	hexaCtx := arranger.HexaCtx(ctx)
	fmt.Printf("hexa correlation id in activity: %s\n", hexaCtx.CorrelationID())
	s := fmt.Sprintf("Hello %s", name)
	fmt.Println(s)

	return s, nil
}
