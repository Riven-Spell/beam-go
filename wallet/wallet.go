package wallet

type WalletStatus struct {
	CurrentHeight uint64 `json:"current_height"`
	CurrentStateHash string `json:"current_state_hash"`
	PreviousStateHash string `json:"prev_state_hash"`
	Available uint64 `json:"available"`
	Receiving uint64 `json:"receiving"`
	Sending uint64 `json:"sending"`
	Maturing uint64 `json:"maturing"`
	Locked uint64 `json:"locked"`
	Difficulty float64 `json:"difficulty"`
}

func (c *Client) WalletStatus() (status WalletStatus, err error) {
	err = c.basicRequest("wallet_status", nil, &status)
	return
}
