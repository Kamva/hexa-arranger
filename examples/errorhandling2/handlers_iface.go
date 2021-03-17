package main

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type Handlers interface {
	PrintMessageActivity(ctx context.Context, message Message) (string, error)
	PrintMessageWorkflow(ctx workflow.Context, message Message) error
}


