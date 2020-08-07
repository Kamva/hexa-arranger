package arranger

import (
	"go.uber.org/cadence/encoded"
	"reflect"
)

// protobufDataConverter uses protobuf encoder/decoder when possible, otherwise
//  uses default cadences data convert
type protobufDataConverter struct {
	encoding         encoding
	defaultConverter encoded.DataConverter
}

func (c *protobufDataConverter) ToData(objs ...interface{}) ([]byte, error) {
	if len(objs) == 1 && isTypeByteSlice(reflect.TypeOf(objs[0])) {
		return objs[0].([]byte), nil
	}

	if !c.encoding.IsSupport(objs) {
		return c.defaultConverter.ToData(objs...)
	}
	return c.encoding.Marshal(objs)
}

func (c *protobufDataConverter) FromData(data []byte, to ...interface{}) error {
	if len(to) == 1 && isTypeByteSlice(reflect.TypeOf(to[0])) {
		reflect.ValueOf(to[0]).Elem().SetBytes(data)
		return nil
	}

	if !c.encoding.IsSupport(to) {
		return c.defaultConverter.FromData(data, to...)
	}
	return c.encoding.Unmarshal(data, to)
}

var ProtobufDataConverter = newProtobufDataConverter()

func newProtobufDataConverter() encoded.DataConverter {
	return &protobufDataConverter{
		encoding:         &protobufEncoding{},
		defaultConverter: encoded.GetDefaultDataConverter(),
	}
}

var _ encoded.DataConverter = &protobufDataConverter{}
