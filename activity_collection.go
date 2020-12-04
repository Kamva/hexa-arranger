package arranger

import (
	"context"

	"github.com/kamva/hexa"
)

type Activities struct {
}

func (ac Activities) Ctx(ctx context.Context) hexa.Context {
	return HexaCtx(ctx)
}
