package arranger

import (
	"errors"
	"fmt"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/encoded"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
)

// FactoryOptions contains options which we need to create a new
// Factory
type FactoryOptions struct {
	ClientName     string
	ServiceName    string
	HostAddr       string
	Domain         string
	ClientIdentity string
	MetricsScope   tally.Scope
	Zap            *zap.Logger
	CtxPropagators []workflow.ContextPropagator
	DataConverter  encoded.DataConverter
}

func (o FactoryOptions) Validate() error {
	if len(o.ClientName) == 0 {
		return tracer.Trace(fmt.Errorf("invalid client name: %s to connect to cadence", o.ClientName))
	}
	if len(o.ServiceName) == 0 {
		return tracer.Trace(fmt.Errorf("invalid service name: %s to connect to cadence", o.ServiceName))
	}
	if len(o.HostAddr) == 0 {
		return tracer.Trace(fmt.Errorf("invalid host address: %s to connect to cadence", o.HostAddr))
	}
	if len(o.Domain) == 0 {
		return tracer.Trace(fmt.Errorf("invalid domain: %s to connect to cadence", o.Domain))
	}
	if o.MetricsScope == nil {
		return tracer.Trace(errors.New("metrics scope can not be nil in connection to the Cadence"))
	}
	if o.Zap == nil {
		return tracer.Trace(errors.New("zap can not be nil in connection to the Cadence"))
	}
	return nil
}

// Factory is The Cadence client factory
type (
	Factory interface {
		FactoryOptions() FactoryOptions
		Worker(taskList string, options worker.Options) (worker.Worker, error)
		CadenceClient() (client.Client, error)
		CadenceDomainClient() (client.DomainClient, error)
		WorkflowServiceClient() (workflowserviceclient.Interface, error)
	}

	factory struct {
		o          FactoryOptions
		dispatcher *yarpc.Dispatcher
	}
)

func (f *factory) FactoryOptions() FactoryOptions {
	return f.o
}

func (f *factory) Worker(taskList string, options worker.Options) (worker.Worker, error) {
	service, err := f.WorkflowServiceClient()
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return worker.New(service, f.o.Domain, taskList, options), nil
}

func (f *factory) CadenceClient() (client.Client, error) {
	service, err := f.WorkflowServiceClient()
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return client.NewClient(service, f.o.Domain, f.cadenceClientOptions()), nil
}

func (f *factory) CadenceDomainClient() (client.DomainClient, error) {
	service, err := f.WorkflowServiceClient()
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return client.NewDomainClient(service, f.cadenceClientOptions()), nil
}

func (f *factory) WorkflowServiceClient() (workflowserviceclient.Interface, error) {
	if err := f.setup(); err != nil {
		return nil, tracer.Trace(err)
	}

	if f.dispatcher == nil {
		return nil, tracer.Trace(errors.New("no RPC dispatcher provided to create a connection to Cadence Service"))
	}

	return workflowserviceclient.New(f.dispatcher.ClientConfig(f.o.ServiceName)), nil
}

func (f *factory) setup() error {
	o := f.o
	if f.dispatcher != nil {
		return nil
	}

	if err := o.Validate(); err != nil {
		return err
	}

	ch, err := tchannel.NewChannelTransport(
		tchannel.ServiceName(o.ClientName))
	if err != nil {
		hlog.WithFields("err", tracer.Trace(err)).Error("Failed to create transport channel")
		return tracer.Trace(err)
	}

	hlog.WithFields("ServiceName", o.ServiceName, "HostPort", o.HostAddr).
		Debug("Creating RPC dispatcher outbound")

	f.dispatcher = yarpc.NewDispatcher(yarpc.Config{
		Name: f.o.ClientName,
		Outbounds: yarpc.Outbounds{
			f.o.ServiceName: {Unary: ch.NewSingleOutbound(o.HostAddr)},
		},
	})

	if f.dispatcher != nil {
		if err := f.dispatcher.Start(); err != nil {
			hlog.WithFields("err", tracer.Trace(err)).Error("Failed to create outbound transport channel")
		}
	}

	return nil
}

func (f *factory) cadenceClientOptions() *client.Options {
	o := f.o
	return &client.Options{
		Identity:           o.ClientIdentity,
		MetricsScope:       o.MetricsScope,
		DataConverter:      o.DataConverter,
		ContextPropagators: o.CtxPropagators,
	}
}

// NewFactory returns new instance of the factory
func NewFactory(o FactoryOptions) Factory {
	return &factory{o: o}
}
