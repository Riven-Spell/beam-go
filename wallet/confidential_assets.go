package wallet

import (
	"errors"
	"fmt"
	"github.com/BeamMW/beam-go/rpc"
	"github.com/BeamMW/beam-go/to"
	"strconv"
	"strings"
)

type AssetDescriptor struct {
	// metadata schema version 1

	Name, ShortName, UnitName, NthUnitName string // standard properties, mandatory

	NthRatio *uint64 // standard properties, non-mandatory

	ShortDescription, LongDescription, SiteURL, PDFURL, FaviconURL, LogoURL *string // optional properties
}

func (a *AssetDescriptor) String() string {
	output := fmt.Sprintf(
		"STD:SCH_VER=1;N=%s;SN=%s;UN=%s;NTHUN=%s",
		a.Name, a.ShortName, a.UnitName, a.NthUnitName,
	)

	if a.NthRatio != nil {
		output += fmt.Sprintf(";%d", *a.NthRatio)
	}

	if a.ShortDescription != nil {
		output += ";OPT_SHORT_DESC=" + *a.ShortDescription
	}

	if a.LongDescription != nil {
		output += ";OPT_LONG_DESC=" + *a.LongDescription
	}

	if a.SiteURL != nil {
		output += ";OPT_SITE_URL=" + *a.SiteURL
	}

	if a.PDFURL != nil {
		output += ";OPT_PDF_URL=" + *a.PDFURL
	}

	if a.FaviconURL != nil {
		output += ";OPT_FAVICON_URL=" + *a.FaviconURL
	}

	if a.LogoURL != nil {
		output += ";OPT_LOGO_URL=" + *a.LogoURL
	}

	return output
}

// Custom json marshaling handlers, for ease of use
func (a *AssetDescriptor) MarshalJSON() ([]byte, error) {
	return []byte("\"" + a.String() + "\""), nil
}

func (a *AssetDescriptor) UnmarshalJSON(data []byte) (err error) {
	*a, err = ParseAssetDescriptor(strings.Trim(string(data), "\""))

	return
}

func ParseAssetDescriptor(input string) (AssetDescriptor, error) {
	var desc AssetDescriptor

	segments := strings.Split(strings.TrimPrefix(input, "STD:"), ";")

	for _,v := range segments {
		splits := strings.SplitN(v, "=", 1)

		switch splits[0] {
		default:
			return AssetDescriptor{}, errors.New("invalid descriptor property name " + splits[0])
		case "SCH_VER": // do not handle
		case "N":
			desc.Name = splits[1]
		case "SN":
			desc.ShortName = splits[1]
		case "UN":
			desc.UnitName = splits[1]
		case "NTHUN":
			desc.NthUnitName = splits[1]
		case "NTH_RATIO":
			rat, err := strconv.ParseUint(splits[1], 10, 64)

			if err != nil {
				return AssetDescriptor{}, err
			}

			desc.NthRatio = &rat
		case "OPT_SHORT_DESC": // out of an abundance of caution, we put these through to.StringPtr. Probably not necessary. Probably the same thing. Still alleviates gut feel.
			desc.ShortDescription = to.StringPtr(splits[1])
		case "OPT_LONG_DESC":
			desc.LongDescription = to.StringPtr(splits[1])
		case "OPT_SITE_URL":
			desc.SiteURL = to.StringPtr(splits[1])
		case "OPT_PDF_URL":
			desc.PDFURL = to.StringPtr(splits[1])
		case "OPT_FAVICON_URL":
			desc.FaviconURL = to.StringPtr(splits[1])
		case "OPT_LOGO_URL":
			desc.LogoURL = to.StringPtr(splits[1])
		}
	}

	return desc, nil
}

type GetAssetInfoResponse struct {
	AssetId uint64 `json:"asset_id"`
	Emission uint64 `json:"emission"`
	IsOwned bool `json:"isOwned"`
	LockHeight uint64 `json:"lockHeight"`
	Metadata AssetDescriptor `json:"metadata"`
	OwnerId string `json:"owner_id"`
	RefreshHeight uint64 `json:"refreshHeight"`
}

// GetAssetInfo obtains the latest asset info and provides it in the form of a AssetDescriptor.
// TODO for reviewers: Should these calls _be_ merged? I feel like it's good for usability.
func (c *Client) GetAssetInfo(selector AssetSelectionOption) (resp GetAssetInfoResponse, err error) {
	var params rpc.JsonParams

	if params, err = selector.rpcPrepare(); err != nil {
		return
	}

	if err = c.basicRequest("tx_asset_info", params, nil); err != nil {
		return
	}

	err = c.basicRequest("get_asset_info", params, &resp)
	return
}

// IssueAsset mints new tokens. You must be the owner of the asset, and the asset should be in the local database.
// Minting is free, and only costs the regular transaction fee.
func (c *Client) IssueAsset(groth uint64, options ModifyAssetOptions) (transactionId string, err error) {
	var resp struct { TransactionId string `json:"txId"` }
	var params rpc.JsonParams

	if params, err = options.rpcPrepare(groth); err != nil {
		return
	}

	err = c.basicRequest("tx_asset_issue", params, &resp)

	return resp.TransactionId, err
}

// BurnAsset burns existing asset tokens. You must own the asset, and you must own the coins you wish to burn.
// Burning is free, and only costs the regular transaction fee.
func (c *Client) BurnAsset(groth uint64, options ModifyAssetOptions) (transactionId string, err error) {
	var resp struct { TransactionId string `json:"txId"` }
	var params rpc.JsonParams

	if params, err = options.rpcPrepare(groth); err != nil {
		return
	}

	err = c.basicRequest("tx_asset_consume", params, &resp)

	return resp.TransactionId, err
}
