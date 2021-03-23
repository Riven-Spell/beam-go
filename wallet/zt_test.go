package wallet_test

import (
	"github.com/BeamMW/beam-go/rpc"
	"github.com/BeamMW/beam-go/wallet"
	chk "gopkg.in/check.v1"
	"os"
	"testing"
)

func Test(t *testing.T) { chk.TestingT(t) }

var (
	walletEndpoint = os.Getenv("WALLET_ENDPOINT")
)

func GetWalletClient(c *chk.C) *wallet.Client {
	c.Assert(walletEndpoint, chk.Not(chk.Equals), "")

	return wallet.NewClient(&rpc.HTTPEndpoint{
		Endpoint: walletEndpoint,
	})
}
