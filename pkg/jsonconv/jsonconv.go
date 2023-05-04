package jsonconv

import (
	"bytes"
	"encoding/json"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"
)

// JSONConverter is responsible for handling jsonpb marshalling and defaults to regular json marshaling
// if the object is not of type proto.Message.
type JSONConverter struct {
	Marshaler *jsonpb.Marshaler
}

// NewJSONConverter is a constructor function for JSONConverter.
func NewJSONConverter() *JSONConverter {
	return &JSONConverter{
		Marshaler: &jsonpb.Marshaler{},
	}
}

// Marshal handles jsonpb marshalling when necessary and defaults to regular json marshaling otherwise.
func (j *JSONConverter) Marshal(obj any) ([]byte, error) {
	var (
		jsonBytes []byte
		err       error
		buf       bytes.Buffer
	)

	protoReq, ok := obj.(proto.Message)
	if ok {
		err = j.Marshaler.Marshal(&buf, protoReq)
		if err == nil {
			jsonBytes = buf.Bytes()
		}
	} else {
		jsonBytes, err = json.Marshal(obj)
	}

	return jsonBytes, err
}
