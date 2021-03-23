package wallet_test

import (
	"github.com/BeamMW/beam-go/to"
	"github.com/BeamMW/beam-go/wallet"
	chk "gopkg.in/check.v1"
	"time"
)

type addressSuite struct{}

var _ = chk.Suite(&addressSuite{})

func (*addressSuite) TestNewAddress(c *chk.C) {
	client := GetWalletClient(c)

	// Test a no-frills address generation.
	address, err := client.CreateAddress(nil)

	c.Assert(err, chk.IsNil)
	c.Assert(address.Address, chk.Not(chk.Equals), "")

	// Test with a couple special parameters. Try some 6.0 exclusive stuff.
	address, err = client.CreateAddress(&wallet.CreateAddressOptions{
		AddressType:            wallet.AddressTypeNewRegular.ToPtr(),
		Expiration:             wallet.ExpirationAuto.ToPtr(),
		Comment:                to.StringPtr("owo"),
	})

	c.Assert(err, chk.IsNil)
	c.Assert(address.Address, chk.Not(chk.Equals), "")
	c.Assert(address.Duration, chk.Equals, int64(time.Hour * 24 * 61 / time.Second))
}

func (*addressSuite) TestValidateAddress(c *chk.C) {
	// try validating some keyboard spam junk.
	address := wallet.Address{Address: "saoeutsaotnh uaoestnh sanoht snaotheu "}
	client := GetWalletClient(c)

	valid, owner, err := client.ValidateAddress(address)

	c.Assert(err, chk.IsNil)
	c.Assert(valid, chk.Equals, false)
	c.Assert(owner, chk.Equals, false)

	// this is a real mainnet address, should validate, but not be "ours"
	address.Address = "11437de1b63c3db491460a6077423b5add63ba582cb2aa5eb9e157f4d91ddc56ea5"

	valid, owner, err = client.ValidateAddress(address)

	c.Assert(err, chk.IsNil)
	c.Assert(valid, chk.Equals, true)
	c.Assert(owner, chk.Equals, false)

	// Generate an address and use it.
	address, err = client.CreateAddress(nil)
	c.Assert(err, chk.IsNil)

	valid, owner, err = client.ValidateAddress(address)

	c.Assert(err, chk.IsNil)
	c.Assert(valid, chk.Equals, true)
	c.Assert(owner, chk.Equals, true)
}

func (*addressSuite) TestListAddresses(c *chk.C) {
	// Generate an address with some special properties, and check those.
	client := GetWalletClient(c)

	const addressComment = "Hello world!"

	address, err := client.CreateAddress(&wallet.CreateAddressOptions{
		Comment:                to.StringPtr(addressComment),
	})
	c.Assert(err, chk.IsNil)

	addresses, err := client.ListAddresses(true)

	c.Assert(err, chk.IsNil)
	c.Assert(len(addresses) >= 1, chk.Equals, true)

	for _,v := range addresses {
		if v.Address == address.Address {
			//c.Assert(v.Duration, chk.Equals, address.Duration) // don't check expiry
			c.Assert(v.Comment, chk.Equals, address.Comment)
			c.Assert(v.IsOwned, chk.Equals, true)
		}
	}
}
