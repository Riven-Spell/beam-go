package rpc

/*

Disabled for now. Currently, this is only usable on the client RPC.
In addition, this has the overhead of needing a new TCP connection every time.
Furthermore, it doesn't currently support TLS.

Thus, TCPEndpoint is currently defunct and unsupported, until the node explorer gets TCP support.

// why netaddr.IPPort over net.IP?
// more better: https://tailscale.com/blog/netaddr-new-ip-type-for-go/
type TCPEndpoint netaddr.IPPort

func (endpoint *TCPEndpoint) RPCExecute(method string, params interface{}) (res *json.RawMessage, err error) {
	var pbytes []byte
	if pbytes, err = json.Marshal(params); err != nil {
		return
	}

	rawId := json.RawMessage(`"tcp-dummy"`)
	rawParams := json.RawMessage(pbytes)

	request := jsonrpc.RequestHeader{
		Jsonrpc: "2.0",
		Id:      &rawId,
		Method:  method,
		Params:  &rawParams,
	}

	var rbytes []byte
	if rbytes, err = json.Marshal(request); err != nil {
		return
	}

	var c net.Conn
	if c, err = net.Dial("tcp", netaddr.IPPort(*endpoint).String()); err != nil {
		return
	}

	defer c.Close()

	if _, err = c.Write(rbytes); err != nil {
		return
	}

	var body []byte
	if body, err = bufio.NewReader(c).ReadBytes('\n'); err != nil {
		return
	}

	var rpcr jsonrpc.ResponseHeader
	if err = json.Unmarshal(body, &rpcr); err != nil {
		return
	}

	res = rpcr.Result
	return
}

 */