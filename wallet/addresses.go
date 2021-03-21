package wallet

import (
	"github.com/BeamMW/beam-go/rpc"
	"time"
)

type Address struct {
	Address string `json:"address"`
	Comment string `json:"comment"`

	Creation int64 `json:"create_time"`
	Duration int64 `json:"duration"` // 0 if auto or never
	Expired bool `json:"expired"` // assumed not expired at creation

	IsOwned bool `json:"own"`
	OwnId uint64 `json:"own_id"` // not exposed during creation

	Identity string `json:"identity"` // not exposed during creation
}

// ImportAddress does not connect to the RPC.
// It merely creates a new Address struct, as a usability function.
func ImportAddress(addr string, owned bool) Address {
	return Address{ Address: addr, IsOwned: owned }
}

func (c *Client) CreateAddress(createAddressOptions *CreateAddressOptions) (a Address, err error) {
	var address string
	if err = c.basicRequest("create_address", createAddressOptions.rpcPrepare(), &address); err != nil {
		return
	}

	a = Address{
		Address: address,
		Creation: time.Now().UTC().Unix(),
		IsOwned: true,
	}

	if createAddressOptions.Comment != nil {
		a.Comment = *createAddressOptions.Comment
	}

	if createAddressOptions.Expiration != nil {
		switch *createAddressOptions.Expiration {
		case Expiration24h:
			a.Duration = int64(time.Hour * 24 / time.Second)
		case ExpirationAuto:
			a.Duration = int64(time.Hour * 24 * 61)
		}
	} else {
		a.Duration = int64(time.Hour * 24 / time.Second)
	}

	return
}

func (c *Client) ValidateAddress(a Address) (valid, owner bool, err error) {
	var output struct {
		valid bool `json:"is_valid"`
		owned bool `json:"is_mine"`
	}
	if err = c.basicRequest("validate_address", rpc.JsonParams{ "address" : a.Address }, &output); err != nil {
		return
	}

	return output.valid, output.owned, err
}

func (c *Client) ListAddresses(ownedOnly bool) (addresses []Address, err error) {
	err = c.basicRequest("addr_list", rpc.JsonParams{ "own" : ownedOnly }, &addresses)
	return
}

func (c *Client) DeleteAddress(a Address) error {
	return c.basicRequest("delete_address", rpc.JsonParams{"address": a.Address}, nil)
}

func (c *Client) EditAddress(a Address, options EditAddressOptions) error {
	return c.basicRequest("edit_address", options.rpcPrepare(a), nil)
}
