package arranger

import (
	"context"

	"github.com/kamva/hexa"
)

// A good type to use in our activity resources
// to include some good methods.
type Activities struct {
}

func (ac Activities) Ctx(ctx context.Context) hexa.Context {
	return HexaCtx(ctx)
}

// ApplicationErr converts error to Application error if provided error is
// Hexa error, otherwise returns error untouched.
func (ac Activities) ApplicationErr(ctx hexa.Context, err error) error {
	return HexaToApplicationErr(err, ctx.Translator())
}
