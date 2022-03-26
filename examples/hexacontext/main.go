package main

import (
	"context"
	"flag"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa-arranger"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	hostPort      = "127.0.0.1:7233"
	namespace     = "arrangerlab"
	taskQueueName = "arranger-ctx-propagation"
)

var logger = hlog.NewPrinterDriver(hlog.DebugLevel)
var translator = hexatranslator.NewEmptyDriver()
var p = hexa.NewContextPropagator(logger, translator)

func boot() arranger.Arranger {
	c, err := client.NewClient(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		Logger:    arranger.NewLogger(hlog.NewPrinterDriver(hlog.DebugLevel)),
		ContextPropagators: []workflow.ContextPropagator{
			arranger.NewHexaContextPropagator(p),
		},
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
	workerOptions := worker.Options{
		EnableLoggingInReplay: true,
	}
	w := worker.New(arr.Client(), taskQueueName, workerOptions)

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
	ctx := hexa.NewContext(context.Background(),hexa.ContextParams{
		CorrelationId: "my_correlation_id",
		Locale:        "en",
		User:          hexa.NewGuest(),
		BaseLogger:        logger,
		BaseTranslator:    translator,
	})
	e, err := arr.ExecuteWorkflow(ctx, workflowOptions, HelloWorldWorkflow, "Mehran")
	if err != nil {
		return tracer.Trace(err)
	}

	hlog.With(
		hlog.String("WorkflowID", e.GetRunID()),
		hlog.String("RunID", e.GetRunID()),
	).Info("Start workflow!!")

	return nil
}
