package wallet

type UTXO struct {
	Id uint64 `json:"id"`
	AssetId uint64 `json:"asset_id"` // 0 = beam
	Amount uint64 `json:"amount"` // in groth
	Maturity uint64 `json:"maturity"` // in confirmations
	Type string `json:"type"`
	CreateTransactionId string `json:"createTxId"`
	SpentTransactionId string `json:"spentTxId"`
	Status uint64 `json:"status"`
	StatusString string `json:"status_string"`
}

func (c *Client) GetUnlockedUTXOs(options *GetUTXOOptions) (utxos []UTXO, err error) {
	err = c.basicRequest("get_utxo", options.rpcPrepare(), &utxos)
	return
}
