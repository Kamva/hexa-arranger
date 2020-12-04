package arranger

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kamva/hexa"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/workflow"
)

// hexaCtxKey is the which we use set exported hexa context
// in cadence context propagator.
const hexaCtxKey = "_hexa_ctx"

// hexaContextPropagator propagate hexa context.
type hexaContextPropagator struct {
	// if set strict flag to true, you must set hexa
	// context on all calls to workflow,activity,....
	strict bool
	cei    hexa.ContextExporterImporter
}

func (h *hexaContextPropagator) Inject(ctx context.Context, hw workflow.HeaderWriter) error {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		if h.strict {
			return errors.New("you must provide hexa context when strict mode is enabled")
		}
		return nil
	}
	m, err := h.cei.Export(hexaCtx.(hexa.Context))
	if err != nil {
		return err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	hw.Set(hexaCtxKey, payload(b))
	return nil
}

func (h *hexaContextPropagator) Extract(ctx context.Context, hr workflow.HeaderReader) (context.Context, error) {
	var b []byte
	err := hr.ForEachKey(func(key string, p *common.Payload) error {
		data := p.Data
		if key == hexaCtxKey {
			b = data
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if b == nil {
		if h.strict {
			return nil, errors.New("we can not find propagated hexa context while strict mode is enabled")
		}
		return ctx, nil
	}
	var m hexa.Map
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	newCtx, err := h.cei.Import(m)
	return context.WithValue(ctx, hexaCtxKey, newCtx), err
}

func (h *hexaContextPropagator) InjectFromWorkflow(ctx workflow.Context, hw workflow.HeaderWriter) error {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		if h.strict {
			return errors.New("you must provide hexa context when strict mode is enabled")
		}
		return nil
	}
	m, err := h.cei.Export(hexaCtx.(hexa.Context))
	if err != nil {
		return err
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	hw.Set(hexaCtxKey, payload(b))
	return nil
}

func (h *hexaContextPropagator) ExtractToWorkflow(ctx workflow.Context, hr workflow.HeaderReader) (workflow.Context, error) {
	var b []byte
	err := hr.ForEachKey(func(key string, p *common.Payload) error {
		data := p.Data
		if key == hexaCtxKey {
			b = data
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if b == nil {
		if h.strict {
			return nil, errors.New("we can not find propagated hexa context while strict mode is enabled")
		}
		return ctx, nil
	}
	var m hexa.Map
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	newCtx, err := h.cei.Import(m)
	return workflow.WithValue(ctx, hexaCtxKey, newCtx), err
}

// NewHexaContextPropagator returns new instance of hexa context propagator.
func NewHexaContextPropagator(cei hexa.ContextExporterImporter, strict bool) workflow.ContextPropagator {
	return &hexaContextPropagator{cei: cei, strict: strict}
}

func payload(data []byte) *common.Payload {
	return &common.Payload{Data: data}
}

var _ workflow.ContextPropagator = &hexaContextPropagator{}
