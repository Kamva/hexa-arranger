package main

import (
	"context"
	"flag"
	"github.com/kamva/hexa-arranger"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/pborman/uuid"
	"github.com/uber-go/tally"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	clientName   = "arranger-helloworld"
	serviceName  = "cadence-frontend"
	hostAddr     = "127.0.0.1:7933"
	domain       = "arrangerlab"
	taskListName = "arranger-helloworld-tasklist"
)

func main() {
	var mode string
	flag.StringVar(&mode, "m", "trigger", "Mode can be worker to start worker or trigger to rigger workflow")
	flag.Parse()

	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zapcore.InfoLevel)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	factory := arranger.NewFactory(arranger.FactoryOptions{
		ClientName:     clientName,
		ServiceName:    serviceName,
		HostAddr:       hostAddr,
		Domain:         domain,
		MetricsScope:   tally.NoopScope,
		Zap:            logger,
		CtxPropagators: nil,
		DataConverter:  nil,
	})

	myArranger, err := arranger.New(factory)
	if err != nil {
		panic(err)
	}

	switch mode {
	case "worker":
		err := startWorker(myArranger)
		if err != nil {
			panic(err)
		}
	case "trigger":
		err := triggerWorkflow(myArranger)
		if err != nil {
			panic(err)
		}
	default:
		panic("unknown mode, trigger mode can be either worker or trigger.")
	}
}

func startWorker(arranger arranger.Arranger) error {
	workerOptions := worker.Options{
		MetricsScope: tally.NoopScope,
		Logger:       arranger.FactoryOptions().Zap,
	}
	w, err := arranger.Worker(taskListName, workerOptions)
	if err != nil {
		return tracer.Trace(err)
	}
	registerWorkflow(w)
	return w.Run()
}

func triggerWorkflow(arranger arranger.Arranger) error {
	workflowClient, err := arranger.CadenceClient()
	if err != nil {
		return tracer.Trace(err)
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "helloworld_" + uuid.New(),
		TaskList:                        taskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}
	msg := Message{Msg: "Hello from printer :)"}
	e, err := workflowClient.StartWorkflow(context.Background(), workflowOptions, PrintMessageWorkflow, msg)
	if err != nil {
		return tracer.Trace(err)
	}

	hlog.WithFields("WorkflowID", e.ID, "RunID", e.RunID).Info("Start workflow!")

	return nil
}
