package rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// instead of requiring a url.URL, a string is used instead.
// this is because the errors from url.Parse aren't very useful, and url.Parse accepts garbage regularly.
// thus, it's just not worth bothering with.
type HTTPEndpoint struct {
	Endpoint string
	Transport *http.Transport // specify a custom HTTP transport to be used, or leave at nil for default transport
}

func (endpoint *HTTPEndpoint) RPCExecute(method string, params interface{}) (res *json.RawMessage, err error) {
	var pbytes []byte
	if pbytes, err = json.Marshal(params); err != nil {
		return
	}

	rawId := json.RawMessage(`"http-dummy"`)
	rawParams := json.RawMessage(pbytes)

	request := RequestHeader {
		Jsonrpc: "2.0",
		Id:      &rawId,
		Method:  method,
		Params:  &rawParams,
	}

	var rbytes []byte
	if rbytes, err = json.Marshal(request); err != nil {
		return
	}

	var rbuffer = bytes.NewBuffer(rbytes)
	var resp *http.Response

	var client = http.DefaultClient

	if endpoint.Transport != nil {
		client.Transport = endpoint.Transport
	}

	if resp, err = client.Post(endpoint.Endpoint, "application/json", rbuffer); err != nil {
		return
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	var rpcr ResponseHeader
	if err = json.Unmarshal(body, &rpcr); err != nil {
		return
	}

	res = rpcr.Result

	if rpcr.Error != (ResponseError{}) {
		err = rpcr.Error
	}

	return
}
