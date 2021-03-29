package wallet

import "github.com/BeamMW/beam-go/rpc"

type StatusBasic struct {
	CurrentHeight uint64 `json:"current_height"`
	CurrentStateHash string `json:"current_state_hash"`
	PreviousStateHash string `json:"prev_state_hash"`
	CurrencyStatus
	Difficulty float64 `json:"difficulty"`
}

type CurrencyStatus struct {
	Available uint64 `json:"available"`
	Receiving uint64 `json:"receiving"`
	Sending uint64 `json:"sending"`
	Maturing uint64 `json:"maturing"`
	Locked uint64 `json:"locked"`
}

func (c *Client) WalletStatus() (status StatusBasic, err error) {
	err = c.basicRequest("wallet_status", nil, &status)
	return
}

type StatusByAsset struct {
	CurrentHeight uint64 `json:"current_height"`
	CurrentStateHash string `json:"current_state_hash"`
	PreviousStateHash string `json:"prev_state_hash"`
	Totals []CurrencyStatus `json:"totals"`
	Difficulty float64 `json:"difficulty"`
}

func (c *Client) WalletStatusByAsset() (status StatusByAsset, err error) {
	err = c.basicRequest("wallet_status", rpc.JsonParams{"assets": true}, &status)
	return
}