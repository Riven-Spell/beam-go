package wallet

import "github.com/BeamMW/beam-go/rpc"

func (c *Client) Send(to Address, groth uint64, options *SendOptions) (transactionId string, err error) {
	var result struct { txId string }
	err = c.basicRequest("tx_send", options.rpcPrepare(to, groth), &result)
	return result.txId, err
}

func (c *Client) Split(coins []uint64, options *SplitOptions) (transactionId string, err error) {
	var result struct { txId string }
	err = c.basicRequest("tx_send", options.rpcPrepare(coins), &result)
	return result.txId, err
}

func (c *Client) CancelTransaction(transactionId string) error {
	return c.basicRequest("tx_cancel", rpc.JsonParams{ "txId": transactionId }, nil)
}

func (c *Client) DeleteTransaction(transactionId string) error {
	return c.basicRequest("tx_delete", rpc.JsonParams{ "txId": transactionId }, nil)
}

type TransactionStatus struct {
	TransactionId string `json:"txId"`
	AssetId string `json:"asset_id"`
	Comment string `json:"comment"`
	Fee uint64 `json:"fee"`
	Kernel string `json:"kernel"`
	Receiver string `json:"receiver"`
	Sender string `json:"sender"`
	Status uint64 `json:"status"` // todo: enum instead of status + string
	StatusString string `json:"status_string"`
	TransactionType uint64 `json:"tx_type"`
	TransactionTypeString string `json:"tx_type_string"`
	FailureReason string `json:"failure_reason"`
	Value uint64 `json:"value"` // in groth
	CreationTime uint64 `json:"create_time"`
	Income bool `json:"income"`
	SenderIdentity string `json:"sender_identity"`
	ReceiverIdentity string `json:"receiver_identity"`
	Token string `json:"token"`
}

func (c *Client) TransactionStatus(transactionId string) (txStatus TransactionStatus, err error) {
	err = c.basicRequest("tx_status", rpc.JsonParams{ "txId": transactionId }, &txStatus)
	return
}

func (c *Client) TransactionList(options *TransactionListOptions) (transactions []TransactionStatus, err error) {
	err = c.basicRequest("tx_list", options.rpcPrepare(), &transactions)
	return
}

func (c *Client) GenerateTransactionId() (transactionId string, err error) {
	err = c.basicRequest("generate_tx_id", nil, &transactionId)
	return
}

func (c *Client) CalculateChange(groth uint64, assetId uint64, fee uint64, isPushTransaction bool) (assetChange, beamChange, explicitFee uint64, err error) {
	var resp struct {
		AssetChange uint64 `json:"asset_change"`
		BeamChange uint64 `json:"change"`
		ExplicitFee uint64 `json:"explicit_fee"`
	}

	err = c.basicRequest("calc_change",
		rpc.JsonParams{"amount": groth, "asset_id": assetId, "fee": fee, "is_push_transaction": isPushTransaction},
		&resp)

	return resp.AssetChange, resp.BeamChange, resp.ExplicitFee, err
}