package arranger

import (
	"github.com/kamva/hexa"
	"go.temporal.io/sdk/workflow"
)

type Workflows struct {
}

// Ctx converts workflow context to hexa context.
func (ac Workflows) Ctx(ctx workflow.Context) hexa.Context {
	return HexaCtxFromCadenceCtx(ctx)
}
