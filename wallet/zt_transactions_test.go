package wallet_test

import (
	"github.com/BeamMW/beam-go/wallet"
	"github.com/dnaeon/go-vcr/recorder"
	chk "gopkg.in/check.v1"
	"os"
	"time"
)

type transactionSuite struct{
	recorder *recorder.Recorder
}

var _ = chk.Suite(&transactionSuite{})

func (s *transactionSuite) SetUpSuite(c *chk.C) {
	if os.Getenv("NORECORD") == "" {
		var err error
		s.recorder, err = recorder.New("recordings/transactions")

		c.Assert(err, chk.IsNil)
	}
}

func (s *transactionSuite) TearDownSuite(c *chk.C) {
	if os.Getenv("NORECORD") == "" {
		c.Assert(s.recorder.Stop(), chk.IsNil)
	}
}

func (s *transactionSuite) GetRecordingState() recorder.Mode {
	if os.Getenv("NORECORD") == "" {
		return s.recorder.Mode()
	}

	return recorder.ModeDisabled
}

func (s *transactionSuite) TestWalletStatus(c *chk.C) {
	w := GetWalletClient(c, s.recorder, true)

	stat, err := w.WalletStatus()
	c.Assert(err, chk.IsNil)
	c.Assert(stat.PreviousStateHash, chk.Not(chk.Equals), "")
}

func (s *transactionSuite) TestSendTx(c *chk.C) {
	// await 1 beam from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)

	mw := GetWalletClient(c, s.recorder, true)
	rw := GetWalletClient(c, s.recorder, false)

	// create an address to receive
	rwAddress, err := rw.CreateAddress(&wallet.CreateAddressOptions{})
	c.Assert(err, chk.IsNil)

	// send some funds, wait for their arrival off-recorder
	txid, err := mw.Send(rwAddress, wallet.Beam, nil)
	c.Assert(err, chk.IsNil)
	AwaitTransactionCompleted(c, true, s.GetRecordingState(), txid)

	// check that funds have been received.
	txStatus, err := rw.TransactionStatus(txid)
	c.Assert(err, chk.IsNil)
	c.Assert(txStatus.Status, chk.Equals, wallet.TransactionStatusCompleted)
}

func (s *transactionSuite) TestSplitTx(c *chk.C) {
	// Await 2 beams from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam * 2)

	// Split the two beams into separate outputs
	w := GetWalletClient(c, s.recorder, true)
	txId, err := w.Split([]uint64{wallet.Beam, wallet.Beam}, nil)
	c.Assert(err, chk.IsNil)

	// Beam wallet RPC responds with EOF if we don't wait
	time.Sleep(time.Second)

	// Await the transaction's completion
	AwaitTransactionCompleted(c, true, s.GetRecordingState(), txId)
}

func (s *transactionSuite) TestCancelTx(c *chk.C) {
	// Await 1 beam from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)

	// get a destination address
	mw := GetWalletClient(c, s.recorder, true)
	rw := GetWalletClient(c, s.recorder, false)
	addr, err := rw.CreateAddress(nil)
	c.Assert(err, chk.IsNil)

	// create a send transaction
	txId, err := mw.Send(addr, wallet.Beam, nil)
	c.Assert(err, chk.IsNil)

	time.Sleep(time.Second)

	// cancel it
	err = mw.CancelTransaction(txId)
	c.Assert(err, chk.IsNil)

	time.Sleep(time.Second)

	// ensure it's cancelled
	status, err := mw.TransactionStatus(txId)
	c.Assert(err, chk.IsNil)
	c.Assert(status.Status, chk.Equals, wallet.TransactionStatusCanceled)
}

func (s *transactionSuite) TestDeleteTx(c *chk.C) {
	// Await 1 beam from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)

	// get a destination address
	mw := GetWalletClient(c, s.recorder, true)
	rw := GetWalletClient(c, s.recorder, false)
	addr, err := rw.CreateAddress(nil)
	c.Assert(err, chk.IsNil)

	// create a send transaction
	txId, err := mw.Send(addr, wallet.Beam, nil)
	c.Assert(err, chk.IsNil)

	time.Sleep(time.Second)

	// cancel it
	err = mw.CancelTransaction(txId)
	c.Assert(err, chk.IsNil)

	time.Sleep(time.Second)

	err = mw.DeleteTransaction(txId)
	c.Assert(err, chk.IsNil)

	time.Sleep(time.Second)

	// ensure it's nonexistant
	_, err = mw.TransactionStatus(txId)
	c.Assert(err, chk.NotNil)
}

func (s *transactionSuite) TestGeneratedTxId(c *chk.C) {
	// Await 1 beam from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)

	// get a destination address
	mw := GetWalletClient(c, s.recorder, true)
	rw := GetWalletClient(c, s.recorder, false)
	addr, err := rw.CreateAddress(nil)
	c.Assert(err, chk.IsNil)

	// create a send transaction w/ a pre-grabbed txId
	genTxId, err := mw.GenerateTransactionId()
	c.Assert(err, chk.IsNil)
	txId, err := mw.Send(addr, wallet.Beam, &wallet.SendOptions{SpecificationOptions: &wallet.TransactionSpecificationOptions{ TransactionId: &genTxId }})

	// check that the transaction has the ID we specified
	c.Assert(err, chk.IsNil)
	c.Assert(txId, chk.Equals, genTxId)

	// ensure that it completes
	AwaitTransactionCompleted(c, true, s.GetRecordingState(), txId)
}

func (s *transactionSuite) TestCalculateChange(c *chk.C) {
	// Await 1 beam from "mining"
	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)

	mw := GetWalletClient(c, s.recorder, true)
	_, beamChange, fee, err := mw.CalculateChange(wallet.Beam / 2, 0, 100, true)
	c.Assert(err, chk.IsNil)
	c.Log("beam change: ", beamChange, " fee: ", fee)
}

// There is currently a bug where payment proofs don't appear to be parsed correctly.
// For now, this test is disabled, and a github issue will be created regarding it on the core beam repo.
//func (s *transactionSuite) TestPaymentProofs(c *chk.C) {
//	// Await 1 beam from "mining"
//	AwaitFundsReady(c, true, s.GetRecordingState(), wallet.Beam)
//
//	// get wallets and an address
//	mw := GetWalletClient(c, s.recorder, true)
//	rw := GetWalletClient(c, s.recorder, false)
//	addr, err := rw.CreateAddress(nil)
//	c.Assert(err, chk.IsNil)
//
//	txId, err := mw.Send(addr, wallet.Beam, nil)
//	c.Assert(err, chk.IsNil)
//
//	time.Sleep(time.Second)
//	AwaitTransactionCompleted(c, true, s.GetRecordingState(), txId)
//
//	proof, err := mw.ExportPaymentProof(txId)
//	c.Assert(err, chk.IsNil)
//	c.Assert(proof, chk.Not(chk.Equals), "")
//
//	c.Log(proof)
//
//	validityResp, err := mw.VerifyPaymentProof(txId)
//	c.Assert(err, chk.IsNil)
//	c.Assert(validityResp.Amount, chk.Equals, wallet.Beam)
//	c.Assert(validityResp.IsValid, chk.Equals, true)
//	c.Assert(validityResp.Receiver, chk.Equals, addr.Address)
//}
