package wallet

import (
	"errors"
	"github.com/BeamMW/beam-go/rpc"
)

type AssetSelectionOption struct {
	// Only one should be specified. Meta can only be used if you own the asset.
	AssetId *uint64
	AssetMeta *AssetDescriptor
}

func (o AssetSelectionOption) rpcPrepare() (rpc.JsonParams, error) {
	if o.AssetId != nil {
		return rpc.JsonParams{ "asset_id": *o.AssetId }, nil
	} else if o.AssetMeta != nil {
		return rpc.JsonParams{ "asset_meta": o.AssetMeta.String() }, nil
	} else {
		return nil, errors.New("AssetSelectionOption must have at least one value specified")
	}
}

type ModifyAssetOptions struct {
	selector AssetSelectionOption
	TransactionSpecificationOptions
}

func (o ModifyAssetOptions) rpcPrepare(groth uint64) (params rpc.JsonParams, err error) {
	if params, err = o.selector.rpcPrepare(); err != nil {
		return
	}

	params["value"] = groth

	return params.Merge(o.TransactionSpecificationOptions.rpcPrepare()), nil
}