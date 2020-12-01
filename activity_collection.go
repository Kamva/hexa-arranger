package arranger

import (
	"context"
	"github.com/kamva/hexa"
)

type ActivityCollection struct {
}

func (ac ActivityCollection) Ctx(ctx context.Context) hexa.Context {
	return HexaCtx(ctx)
}
