package wallet

import "github.com/BeamMW/beam-go/rpc"

func (c *Client) ExportPaymentProof(transactionId string) (paymentProof string, err error) {
	var resp struct { PaymentProof string `json:"payment_proof"` }
	err = c.basicRequest("export_payment_proof", rpc.JsonParams{"txId": transactionId}, &resp)
	return resp.PaymentProof, err
}

type VerifyProofResponse struct {
	IsValid bool `json:"is_valid"`
	AssetId uint64 `json:"asset_id"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount uint64 `json:"amount"` // in groth
	Kernel string `json:"kernel"`
}

func (c *Client) VerifyPaymentProof(proof string) (resp VerifyProofResponse, err error) {
	err = c.basicRequest("verify_payment_proof", rpc.JsonParams{"payment_proof": proof}, &resp)
	return
}
