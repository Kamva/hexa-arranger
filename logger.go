package arranger

import (
	"fmt"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"go.temporal.io/sdk/log"
)

// logger implements temporal logger using hexa logger.
type logger struct {
	logger hexa.Logger
}

func (l *logger) Debug(msg string, keyvals ...interface{}) {
	l.logger.Debug(msg, keyValuesToFields(keyvals)...)
}

func (l *logger) Info(msg string, keyvals ...interface{}) {
	l.logger.Info(msg, keyValuesToFields(keyvals)...)
}

func (l *logger) Warn(msg string, keyvals ...interface{}) {
	l.logger.Warn(msg, keyValuesToFields(keyvals)...)

}

func (l *logger) Error(msg string, keyvals ...interface{}) {
	l.logger.Error(msg, keyValuesToFields(keyvals)...)
}

func keyValuesToFields(keyVals []interface{}) []hexa.LogField {
	if len(keyVals)%2 != 0 {
		lastKey := fmt.Sprint(keyVals[len(keyVals)-1])
		keyVals = append(keyVals, fmt.Sprintf("missed log value for key:%s", lastKey))
	}
	m, err := gutil.KeyValuesToMap(keyVals...)
	if err != nil {
		hlog.Error("can not convert key-values to map", hlog.Err(tracer.Trace(err)))
	}
	return hlog.MapToFields(m)
}

func NewLogger(l hexa.Logger) log.Logger {
	return &logger{logger: l}
}

var _ log.Logger = &logger{}
