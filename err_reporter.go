package arranger

import (
	"net/http"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

var UnknownHexaErr = hexa.NewError(http.StatusInternalServerError, "lib.unknown_hexa_err", nil)

// ReportErr reports our error:
// supported types:
// - wrapped hexa error.
// - ApplicationError which is result of converting hexa error to ApplicationError.
// - otherwise warp error un UnknownHexaErr err.
//
// if error is hexa error, so report hexa error.
// if error is ApplicationError and if we can convert it to hexa error, so we
// convert it before report.
// otherwise wrap error in UnknownHexaErr error before report.
// TODO: we can also extend reporter to report ApplicationErrors which
// are not hexa error(e.g., CancelActivity,...) properly and in a good
// format.
func ReportErr(e error) {
	if e == nil {
		return
	}

	var hexaErr hexa.Error
	hexaErr, ok := hexa.AsHexaErr(e)
	if !ok {
		if err, ok := HexaErrFromApplicationErrWithOk(e); ok {
			hexaErr = err.(hexa.Error)
		}
	}

	if hexaErr == nil {
		hexaErr = UnknownHexaErr.SetError(tracer.Trace(e))
	}

	hexaErr.ReportIfNeeded(hlog.With(), nil)
}
