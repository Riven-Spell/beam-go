package rpc

import (
	"encoding/json"
)

type Endpoint interface {
	RPCExecute(method string, params interface{}) (res *json.RawMessage, err error)
}

type JsonParams map[string]interface{}
