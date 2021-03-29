package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Endpoint interface {
	RPCExecute(method string, params interface{}) (res *json.RawMessage, err error)
}

type JsonParams map[string]interface{}

// the new parameters are treated as overriding the old parameters.
func (j JsonParams) Merge(params JsonParams) JsonParams {
	out := JsonParams{}

	for k,v := range j {
		out[k] = v
	}

	for k,v := range params {
		out[k] = v
	}

	return out
}

func (j JsonParams) Compare(j2 JsonParams) bool {
	return reflect.DeepEqual(j, j2)
}

type RequestHeader struct {
	Jsonrpc string           `json:"jsonrpc"`
	Id      *json.RawMessage `json:"id"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
}

type ResponseHeader struct {
	Jsonrpc string           `json:"jsonrpc"`
	Id      *json.RawMessage `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   ResponseError    `json:"error"`
}

type ResponseError struct {
	Code int64 `json:"code"`
	Data string `json:"data"`
	Message string `json:"message"`
}

func (r ResponseError) Error() string {
	return fmt.Sprintf("Error code %d: %s %s", r.Code, r.Message, r.Data)
}
