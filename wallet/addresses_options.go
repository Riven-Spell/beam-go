package wallet

import (
	"github.com/BeamMW/beam-go/rpc"
)

type CreateAddressOptions struct {
	AddressType *AddressType
	Expiration *Expiration
	Comment *string
	// UseNewRegularStyleOnly refers whether to treat the default AddressTypeRegular value as AddressTypeNewRegular
	UseNewRegularStyleOnly *bool
	// OfflinePaymentCount specifies how many payments a regular offline address can receive -- this is NOT for AddressTypePublicOffline, but AddressTypeOffline
	OfflinePaymentCount *uint64
}

func (o *CreateAddressOptions) rpcPrepare() rpc.JsonParams {
	params := rpc.JsonParams{}

	if o == nil {
		return params
	}

	if o.AddressType != nil {
		params["type"] = *o.AddressType
	}

	if o.Expiration != nil {
		params["expiration"] = *o.Expiration
	}

	if o.Comment != nil {
		params["comment"] = *o.Comment
	}

	if o.UseNewRegularStyleOnly != nil {
		params["new_style_regular"] = *o.UseNewRegularStyleOnly
	}

	if o.OfflinePaymentCount != nil {
		params["offline_payments"] = *o.OfflinePaymentCount
	}

	return params
}

type EditAddressOptions struct {
	Expiration *Expiration
	Comment *string
}

func (o *EditAddressOptions) rpcPrepare(a Address) rpc.JsonParams {
	params := rpc.JsonParams{
		"address": a.Address,
	}

	if o == nil {
		return params
	}

	if o.Expiration != nil {
		params["expiration"] = *o.Expiration
	}

	if o.Comment != nil {
		params["comment"] = *o.Comment
	}

	return params
}