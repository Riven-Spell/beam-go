package wallet

import "github.com/BeamMW/beam-go/rpc"

type SendOptions struct {
	Fee     *uint64 // In groth
	From    *Address
	Comment *string
	TransactionId *string
	AssetId *uint64 // beam = 0
}

func (o *SendOptions) rpcPrepare(to Address, groth uint64) rpc.JsonParams {
	params := rpc.JsonParams{
		"address": to.Address,
		"value": groth,
	}

	if o == nil {
		return params
	}

	if o.Fee != nil {
		params["fee"] = *o.Fee
	}

	if o.From != nil {
		params["from"] = o.From.Address
	}

	if o.Comment != nil {
		params["comment"] = *o.Comment
	}

	if o.AssetId != nil {
		params["asset_id"] = *o.AssetId
	}

	if o.TransactionId != nil {
		params["txId"] = *o.TransactionId
	}

	return params
}

type SplitOptions struct {
	Fee *uint64 // In groth
	TransactionId *string
	AssetId *uint64 // beam = 0
}

func (o *SplitOptions) rpcPrepare(splits []uint64) rpc.JsonParams {
	params := rpc.JsonParams{
		"splits": splits,
	}

	if o == nil {
		return params
	}

	if o.TransactionId != nil {
		params["txId"] = *o.TransactionId
	}

	if o.AssetId != nil {
		params["asset_id"] = *o.AssetId
	}

	if o.Fee != nil {
		params["fee"] = *o.Fee
	}

	return params
}

type TransactionListOptions struct {
	Filter *TransactionListFilter
	Skip *uint64
	Count *uint64
	Assets *bool
}

type TransactionListFilter struct {
	Status  *uint64 `json:"status,omitempty"`
	Height  *uint64 `json:"height,omitempty"`
	AssetId *uint64 `json:"asset_id,omitempty"`
}

func (o *TransactionListOptions) rpcPrepare() rpc.JsonParams {
	params := rpc.JsonParams{}

	if o == nil {
		return params
	}

	if o.Assets != nil {
		params["assets"] = *o.Assets
	}

	if o.Count != nil {
		params["count"] = *o.Count
	}

	if o.Skip != nil {
		params["skip"] = *o.Skip
	}

	if o.Filter != nil {
		params["filter"] = *o.Filter
	}

	return params
}