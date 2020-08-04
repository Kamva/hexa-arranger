package arranger

import (
	"github.com/kamva/tracer"
	"go.uber.org/cadence/client"
)

type (
	// Arranger provide cadence clients and all methods of cadence client for us.
	Arranger interface {
		client.Client
		Factory
	}

	// arranger implements the Arranger interface
	arranger struct {
		client.Client
		Factory
	}
)

// New returns new instance of the Arranger
func New(f Factory) (Arranger, error) {
	cadenceClient, err := f.CadenceClient()
	return &arranger{
		Client:  cadenceClient,
		Factory: f,
	}, tracer.Trace(err)
}
