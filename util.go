package arranger

import (
	"reflect"
)

var typeOfByteSlice = reflect.TypeOf(([]byte)(nil))

// isTypeByteSlice checks whether the type passed in is a ByteSlice type
func isTypeByteSlice(inType reflect.Type) bool {
	return inType == typeOfByteSlice || inType == reflect.PtrTo(typeOfByteSlice)
}

// TODO: add this function to gutil.
func valPtr(obj interface{}) interface{} {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		return obj
	}
	vp := reflect.New(val.Type())
	vp.Elem().Set(val)
	return vp.Interface()
}
