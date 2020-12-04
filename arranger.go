package arranger

import (
	"go.temporal.io/sdk/client"
)

type temporalClient client.Client

// Arranger is temporal clients wrapper
type Arranger interface {
	temporalClient
	Client() client.Client
}

// arranger implements the Arranger interface
type arranger struct {
	temporalClient
}

func (a *arranger) Client() client.Client {
	return a.temporalClient
}

func New(c client.Client) Arranger {
	return &arranger{
		temporalClient: c,
	}
}
