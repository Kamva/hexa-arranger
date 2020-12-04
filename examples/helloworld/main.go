package main

import (
	"context"
	"flag"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa-arranger"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	hostPort      = "localhost:7233"
	namespace     = "arrangerlab"
	taskQueueName = "arranger-helloworld-tasklist"
)

func boot() arranger.Arranger {
	c, err := client.NewClient(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		Logger:    arranger.NewLogger(hlog.NewPrinterDriver(hlog.DebugLevel)),
	})
	gutil.PanicErr(err)
	return arranger.New(c)
}

func main() {
	var mode string
	flag.StringVar(&mode, "m", "trigger", "Mode can be worker to start worker or trigger to rigger workflow")
	flag.Parse()

	arr := boot()

	switch mode {
	case "worker":
		err := startWorker(arr)
		if err != nil {
			panic(err)
		}
	case "trigger":
		err := triggerWorkflow(arr)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown mode, trigger mode can be either worker or trigger.")
	}
}

func startWorker(arr arranger.Arranger) error {
	w := worker.New(arr.Client(), taskQueueName, worker.Options{
		EnableLoggingInReplay: false,
	})

	// Register workflows
	registerWorkflowsAndActivities(w)

	// Run worker
	return w.Run(worker.InterruptCh())
}

func triggerWorkflow(arr arranger.Arranger) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:                 "helloworld_" + uuid.New(),
		TaskQueue:          taskQueueName,
		WorkflowRunTimeout: 20 * time.Minute,
	}
	e, err := arr.ExecuteWorkflow(context.Background(), workflowOptions, HelloWorldWorkflow, "Mehran")
	if err != nil {
		return tracer.Trace(err)
	}

	hlog.Info("Start workflow!",
		hlog.String("WorkflowID", e.GetID()),
		hlog.String("RunID", e.GetRunID()),
	)

	return nil
}
