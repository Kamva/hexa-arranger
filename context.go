package arranger

import (
	"context"
	"github.com/kamva/hexa"
	"go.uber.org/cadence/workflow"
)

// Ctx gets hexa context and returns a go context which embed
// our hexa context.
func Ctx(ctx hexa.Context) context.Context {
	return context.WithValue(context.Background(), hexaCtxKey, ctx)
}

// HexaCtx extracts hexa context from regular go context
func HexaCtx(ctx context.Context) hexa.Context {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		return nil
	}

	return hexaCtx.(hexa.Context)
}

// HexaCtxFromCadenceCtx extracts hexa context from Cadence context.
func HexaCtxFromCadenceCtx(ctx workflow.Context) hexa.Context {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		return nil
	}
	return hexaCtx.(hexa.Context)
}
