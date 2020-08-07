package arranger

import (
	"reflect"
)

var typeOfByteSlice = reflect.TypeOf(([]byte)(nil))

// isTypeByteSlice checks whether the type passed in is a ByteSlice type
func isTypeByteSlice(inType reflect.Type) bool {
	return inType == typeOfByteSlice || inType == reflect.PtrTo(typeOfByteSlice)
}