package wallet

import "github.com/BeamMW/beam-go/rpc"

type TransactionSpecificationOptions struct {
	TransactionId *string // can be issued via Client.GenerateTransactionId
	AssetId *uint64
	Fee *uint64
}

func (o *TransactionSpecificationOptions) rpcPrepare() rpc.JsonParams {
	out := rpc.JsonParams{}
	if o == nil {
		return out
	}

	if o.TransactionId != nil {
		out["txId"] = *o.TransactionId
	}

	if o.AssetId != nil {
		out["asset_id"] = *o.AssetId
	}

	if o.Fee != nil {
		out["fee"] = *o.Fee
	}

	return out
}

type SendOptions struct {
	SpecificationOptions *TransactionSpecificationOptions
	From    *Address
	Comment *string
}

func (o *SendOptions) rpcPrepare(to Address, groth uint64) rpc.JsonParams {
	params := rpc.JsonParams{
		"address": to.Address,
		"value": groth,
	}

	if o == nil {
		return params
	}

	if o.From != nil {
		params["from"] = o.From.Address
	}

	if o.Comment != nil {
		params["comment"] = *o.Comment
	}

	return params.Merge(o.SpecificationOptions.rpcPrepare())
}

type SplitOptions struct {
	TransactionSpecificationOptions
}

func (o *SplitOptions) rpcPrepare(splits []uint64) rpc.JsonParams {
	params := rpc.JsonParams{
		"coins": splits,
	}

	if o == nil {
		return params
	}

	return params.Merge(o.TransactionSpecificationOptions.rpcPrepare())
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