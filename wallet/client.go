package wallet

import (
	"encoding/json"
	"github.com/BeamMW/beam-go/rpc"
)

type Client struct {
	endpoint rpc.Endpoint
}

func NewClient(endpoint rpc.Endpoint) *Client {
	return &Client{endpoint: endpoint}
}

// internal code duplication reducer
func (c *Client) basicRequest(method string, params interface{}, unmarshalTo interface{}) (err error) {
	var rawResp *json.RawMessage
	if rawResp, err = c.endpoint.RPCExecute(method, params); err != nil {
		return
	}

	if rawResp != nil {
		if unmarshalTo != nil {
			err = json.Unmarshal(*rawResp, unmarshalTo)
		}
	}

	return
}