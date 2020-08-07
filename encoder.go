package arranger

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/io"
	"github.com/golang/protobuf/proto"
	"github.com/kamva/gutil"
	"reflect"
)

// encoding is capable of encoding and decoding objects
type encoding interface {
	IsSupport(objs []interface{}) bool
	Marshal([]interface{}) ([]byte, error)
	Unmarshal([]byte, []interface{}) error
}

// protobufEncoding encapsulates json encoding and decoding
// This encoding encode data to json, not bindary protobuf
// data, this is because we want to see our payloads in
// workflows and activities, but after we experienced
// with cadence, we will change this behaviour to encode
// our data to the protobuf binary data.
type protobufEncoding struct {
}

func (e *protobufEncoding) IsSupport(objs []interface{}) bool {
	if len(objs) == 0 {
		return false
	}

	for _, v := range objs {
		if !e.isProtobuf(v) {
			return false
		}
	}
	return true
}

func (e *protobufEncoding) isProtobuf(v interface{}) bool {
	_, ok := gutil.MustValuePtr(v).(proto.Message)
	return ok
}

// Marshal encodes an array of object into bytes
func (e protobufEncoding) Marshal(objs []interface{}) ([]byte, error) {
	var buf bytes.Buffer
	dw := io.NewDelimitedWriter(&buf)
	for i, obj := range objs {
		if err := dw.WriteMsg(gutil.MustValuePtr(obj).(proto.Message)); err != nil {
			return nil, fmt.Errorf(
				"unable to encode argument: %d, %v, with protobuf writer error: %v", i, reflect.TypeOf(obj), err)
		}
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes a byte array into the passed in objects
func (e protobufEncoding) Unmarshal(b []byte, objs []interface{}) error {
	dr := io.NewDelimitedReader(bytes.NewReader(b), 1024*1024)

	for i, obj := range objs {
		if err := dr.ReadMsg(gutil.MustValuePtr(obj).(proto.Message)); err != nil {
			return fmt.Errorf(
				"unable to decode argument: %d, %v, with protouf reader error: %v", i, reflect.TypeOf(obj), err)
		}
	}
	return nil
}

// Assertion
var _ encoding = &protobufEncoding{}
