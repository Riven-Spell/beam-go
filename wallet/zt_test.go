package wallet_test

import (
	"github.com/BeamMW/beam-go/rpc"
	"github.com/BeamMW/beam-go/wallet"
	"github.com/dnaeon/go-vcr/recorder"
	chk "gopkg.in/check.v1"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test(t *testing.T) { chk.TestingT(t) }

var (
	walletEndpoint = os.Getenv("MINER_WALLET_ENDPOINT")
	secondaryWalletEndpoint = os.Getenv("SECONDARY_WALLET_ENDPOINT")
)

func TernaryString(_if bool, then, or string) string {
	if _if {
		return then
	}

	return or
}

func GetWalletClient(c *chk.C, r *recorder.Recorder, useMinerWallet bool) *wallet.Client {
	c.Assert(walletEndpoint, chk.Not(chk.Equals), "")
	c.Assert(secondaryWalletEndpoint, chk.Not(chk.Equals), "")

	var transport http.RoundTripper = r
	if r == nil {
		transport = http.DefaultTransport
	}

	return wallet.NewClient(&rpc.HTTPEndpoint{
		Endpoint: TernaryString(useMinerWallet, walletEndpoint, secondaryWalletEndpoint),
		Transport: transport,
	})
}

func AwaitFundsReady(c *chk.C, useMinerWallet bool, recorderState recorder.Mode, minimumFundsInGroth uint64) {
	if recorderState == recorder.ModeRecording || recorderState == recorder.ModeDisabled {
		// actually hold the line until funds are ready, because we're live-testing.
		// supplying no recorder helps reduce logging in the recordings.
		// additionally, it helps keep runtime down when playing back recordings.
		w := GetWalletClient(c, nil, useMinerWallet)

		for {
			basicStatus, err := w.WalletStatus()
			c.Assert(err, chk.IsNil)

			if basicStatus.Available >= minimumFundsInGroth {
				return
			}

			time.Sleep(time.Second)
		}
	}

}

func AwaitTransactionCompleted(c *chk.C, useMinerWallet bool, recorderState recorder.Mode, txId string) {
	if recorderState == recorder.ModeRecording || recorderState == recorder.ModeDisabled {
		// ditto
		w := GetWalletClient(c, nil, useMinerWallet)

		for {
			txInfo, err := w.TransactionStatus(txId)
			c.Assert(err, chk.IsNil)

			if txInfo.Status == wallet.TransactionStatusCompleted {
				return
			} else if  txInfo.Status == wallet.TransactionStatusFailed || txInfo.Status == wallet.TransactionStatusCanceled {
				// The transaction has failed (or been cancelled!)
				c.Log(txInfo.StatusInfo)
				c.Fail()
			}

			time.Sleep(time.Second)
		}
	}
}
