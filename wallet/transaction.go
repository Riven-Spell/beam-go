package wallet

import "github.com/BeamMW/beam-go/rpc"

func (c *Client) Send(to Address, groth uint64, options *SendOptions) (transactionId string, err error) {
	var result struct { TxId string `json:"txId"` }
	err = c.basicRequest("tx_send", options.rpcPrepare(to, groth), &result)
	return result.TxId, err
}

func (c *Client) Split(coins []uint64, options *SplitOptions) (transactionId string, err error) {
	var result struct { TxId string }
	err = c.basicRequest("tx_split", options.rpcPrepare(coins), &result)
	return result.TxId, err
}

func (c *Client) CancelTransaction(transactionId string) error {
	return c.basicRequest("tx_cancel", rpc.JsonParams{ "txId": transactionId }, nil)
}

func (c *Client) DeleteTransaction(transactionId string) error {
	return c.basicRequest("tx_delete", rpc.JsonParams{ "txId": transactionId }, nil)
}

type TransactionStatus uint

const (
	TransactionStatusPending TransactionStatus = iota
	TransactionStatusInProgress
	TransactionStatusCanceled
	TransactionStatusCompleted
	TransactionStatusFailed
	TransactionStatusRegistering
)

type TransactionType uint

const (
	TransactionTypeSimple TransactionType = 0 // Simple send/split transactions.
	TransactionTypeAssetIssue TransactionType = 2 // Issued new tokens on a CA.
	TransactionTypeAssetBurn TransactionType = 3 // Destroyed existing tokens on a CA.
	TransactionTypeAssetInfo TransactionType = 6 // Grabs information about an asset.
)

type TransactionInfo struct {
	TransactionId string `json:"txId"`
	AssetId uint64 `json:"asset_id"`
	Comment string `json:"comment"`
	Fee uint64 `json:"fee"`
	Kernel string `json:"kernel"`
	Receiver string `json:"receiver"`
	Sender string `json:"sender"`
	Status TransactionStatus `json:"status"`
	StatusInfo string `json:"status_string"`
	TransactionType TransactionType `json:"tx_type"`
	FailureReason string `json:"failure_reason"`
	Value uint64 `json:"value"` // in groth
	CreationTime uint64 `json:"create_time"`
	Income bool `json:"income"`
	SenderIdentity string `json:"sender_identity"`
	ReceiverIdentity string `json:"receiver_identity"`
	Token string `json:"token"`
}

func (c *Client) TransactionStatus(transactionId string) (txStatus TransactionInfo, err error) {
	err = c.basicRequest("tx_status", rpc.JsonParams{ "txId": transactionId }, &txStatus)
	return
}

func (c *Client) TransactionList(options *TransactionListOptions) (transactions []TransactionInfo, err error) {
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