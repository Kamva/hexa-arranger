package main

import (
	"context"
	"flag"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa-arranger"
	"github.com/kamva/hexa-arranger/examples/protobufmessage/hello"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const (
	hostPort      = "127.0.0.1:7233"
	namespace     = "arrangerlab"
	taskQueueName = "arranger-protobufmessage"
)

func main() {
	var mode string
	flag.StringVar(&mode, "m", "trigger", "Mode can be worker to start worker or trigger to rigger workflow")
	flag.Parse()

	c, err := client.NewClient(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
		Logger:    arranger.NewLogger(hlog.NewPrinterDriver(hlog.DebugLevel)),
	})
	gutil.PanicErr(err)
	arr := arranger.New(c)

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
		EnableLoggingInReplay: true,
	})

	registerWorkflow(w)
	return w.Run(worker.InterruptCh())
}

func triggerWorkflow(arr arranger.Arranger) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:                 "helloworld_" + uuid.New(),
		TaskQueue:          taskQueueName,
		WorkflowRunTimeout: 20 * time.Minute,
	}
	msg := hello.Message{Msg: "Hello from protoub message :)"}
	e, err := arr.ExecuteWorkflow(context.Background(), workflowOptions, PrintProtobufMessageWorkflow, msg)
	if err != nil {
		return tracer.Trace(err)
	}

	hlog.With(hlog.String("WorkflowID", e.GetID()), hlog.String("RunID", e.GetRunID())).Info("Start workflow!")
	return nil
}
