package wallet

import "github.com/BeamMW/beam-go/rpc"

type UTXOSortField string

func (u UTXOSortField) ToPtr() *UTXOSortField {
	return &u
}

const (
	UTXOSortFieldID UTXOSortField = "id"
	UTXOSortFieldAssetID UTXOSortField = "asset_id"
	UTXOSortFieldAmount UTXOSortField = "amount"
	UTXOSortFieldType UTXOSortField = "type"
	UTXOSortFieldMaturity UTXOSortField = "maturity"
	UTXOSortFieldCreateTransactionId UTXOSortField = "createTxId"
	UTXOSortFieldSpentTransactionId UTXOSortField = "spentTxId"
	UTXOSortFieldStatus UTXOSortField = "status"
	UTXOSortFieldStatusString UTXOSortField = "status_string"
)

type SortDirection string

func (s SortDirection) ToPtr() *SortDirection {
	return &s
}

const (
	SortAscending SortDirection = "asc"
	SortDescending SortDirection = "desc"
)

type UTXOSortOptions struct {
	Field *UTXOSortField `json:"field"`
	Direction *SortDirection `json:"direction"`
}

type GetUTXOOptions struct {
	Sort *UTXOSortOptions
	Count *uint64
	Skip *uint64
}

func (o *GetUTXOOptions) rpcPrepare() (params rpc.JsonParams) {
	params = rpc.JsonParams{}

	if o == nil {
		return
	}

	if o.Sort != nil {
		params["sort"] = *o.Sort
	}

	if o.Skip != nil {
		params["skip"] = *o.Skip
	}

	if o.Count != nil {
		params["count"] = *o.Count
	}

	return
}
