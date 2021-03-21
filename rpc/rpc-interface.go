package rpc

import (
	"encoding/json"
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
