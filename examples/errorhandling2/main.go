package main

import (
	"context"
	"flag"
	"fmt"
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
)

const (
	hostPort      = "127.0.0.1:7233"
	namespace     = "arrangerlab"
	taskQueueName = "arranger-msgprinter"
)

var t = hexatranslator.NewEmptyDriver()

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

	handlers := NewErrInterceptorLayer(NewHandlers(), hlog.With(), t)

	registerWorkflow(w, handlers)
	return w.Run(worker.InterruptCh())
}

func triggerWorkflow(arr arranger.Arranger) error {
	handlers := NewErrInterceptorLayer(NewHandlers(), hlog.With(), t)

	workflowOptions := client.StartWorkflowOptions{
		ID:                 "helloworld_" + uuid.New(),
		TaskQueue:          taskQueueName,
		WorkflowRunTimeout: 20 * time.Minute,
	}

	msg := Message{Msg: "Hello from printer activity :)"}
	e, err := arr.ExecuteWorkflow(context.Background(), workflowOptions, handlers.PrintMessageWorkflow, msg)
	if err != nil {
		return tracer.Trace(err)
	}

	hlog.With(hlog.String("WorkflowID", e.GetID()), hlog.String("RunID", e.GetRunID())).Info("Start workflow!")
	if err := e.Get(context.Background(), nil); err != nil {
		if err, ok := arranger.HexaErrFromApplicationErrWithOk(err); ok {
			hErr := err.(hexa.Error)
			fmt.Println("hexa error in workflow result: ", hErr.HTTPStatus(), hErr.ID())
		} else {
			fmt.Println("returned error from workflow is not a hexa error (:")
		}
	}

	return nil
}
