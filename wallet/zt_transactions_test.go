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
		s.recorder.SetMatcher(BeamWalletAPIMatcher)

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
	addr := wallet.ImportAddress("51su6ZRjAm2AFwNSXJKyZJAeEhYXrhTXpUFW7EQcgbLJVjRCDQiffM4seDi4fzhNd72V2HeRzLbMmYUqeewRfYzWsk7AheVcCmPAQpT6M8j2VF15Q69dsgW7LNxFmtBQNVdWBftD4atNiFkXFFnRF1AeHeTodYwLgkujLXdaVqKeAJcHQQ8Payah7mcJ51LUgb9pbDJ8iRGjEkShFioVXRuyTZHQLxmWpqwnUG3985eobYfV6Rb325XxDU6jBX2cLHtrkDasbQDM7c67aXfvEUE88asgpBLC4KmBsL7vjVJn5YnYAYqnhdDjpApCqT23rur7aj9ytF1X2bmjcmHgFthfj9oGji65VhPNTaika2rLTqR1anhLoGCunkAkYNPzRVp4FGBDwxmRsRbyT6qfUbSM53xzgCePERuDVfYzqa6B2SAZ5Qn6URCUYngGVN3V1dx1spqbPsPAm8vx6uQj4PNwDjrThg82QN7fnu25npoVUTtv4f3CkpkAJUzYD11qLeCQ9Ladhw2HV1F5kYwcHk6gRBNWBhqNYYQh1Ajvd5Tyx2CbGRLdbJoguuuLJJhmrkSfcFk6czY4NeSFQ7DWem6cngKZcBAGe8t3VazSBg1Rhaacq7aX2GCqpXqduecbdZPxnCzf1dU8XY4Mepd97UfpyQtHdMqrXXbVscrSjs7qsc1tCxRbFSrkGTk2dFzfrZyBcRFHjJT3GqwCh3bjH4NuesmEKw9csHLUDoyQ3WmTumPjWg8gh2eS1xdM41bRdcihwa1nMSBDxVrzTwbKEbVAwXqBWP2CinN6kjLQAWMb7KKssfAnuLjW25a5NKg4qLyZNXVnuPFNrqNmxrKe4nEPnH5bmv2BZYVWfvD39oKEAfk6kjm5mMPJoxyBNXXCNcbKexoyJCLazpj6KqUiD8rZFMadUgc5V4cprjjF1EYcChym6yMm8SV82on2e9Jpp5hRzsCopZMN3Agdte1Wx2WYnn7rQWVcPsVnViJLiEzmEG53kKH15xLvPpV6fAiayrtLRqaTzPEBvar5rXpRc4LHCHNQzWAuUymSypDBgqV7J1nXJhKj8AfjpB1viKvW3vM6rqV4Qh7HNfJErGgC7FzCWLZp97748PjmFT61afhhmbHDhnv9YZWWfYznmpRWBa2TvJXmKDK4XEd6p26YaBvkEmSgLLpEhQiQXWy1NQbDufLFB1ybX4vSWWXto9vdzAcLrXWosW4wA1dGmmz3KQM2XKEJwYj9R7hvkUgfxXdz5LdQY32r5rDCKGgAWNHNRhZrCYtXxc46fksM8BDnHQosvrjL97c56fw2eqDwKPpwKt6p2qsGEsXF1USEADgye3vsmp93a99da8mtuxy5857hped2K1hoT9K2daRq937kkDLpRFcJCDbN6w1BtQkMp9925Xmxf6owS1yqp4EPq2xrNg5ryo2hT9p7bcnrApfaPo6Yx71JEWiG32YBMXjZZ1uN14s1h5ex7hfvPziqbPBzmhkfHpBRHRzP4HfdvH4eePCwWbZsDba66BUuCye9YfsCq72Mbn1kEPQErT3YWmRb9KmekA2dZdkxGYZF2RdDB3RwupGZ84Xc7h1WWWZ56WdwrfN4aFv46mXJodtELvPizWCz6uEi87HVNzeQRx74zbLTpc66zeKGRrMZHea8hCtBBNw33fJqR7zKLyf35e9RARFPbXdGBt8YqJkWpD9cBJVt9LgfX57Pr2qokbyf6c8WSbPzAaSHGQSqLLrvHtsFbxaVjjatSoJvmZc3Qxo3ihZnWuWwCR4B6jD88RQYQAxNtQEp1XrN5eKFhr6tFysGprBehSjqp39rzQe7nYPPc4gTLCDE3bFJk2rg7g5YRfWRwnh54MYqz4qqsbjg7pwvNHVVxQruSEoovfadt1wvjqLX3NpDoaDEdC9UDRxUan8HFp4k37cFFgYLuf2BhMBhJGkcS3hK8jkpRivtWzUCiEXPhg2UdAAJ92ePmxTpSoxQ76pVf6QuzQV4xVcnkaWgSmrq9oUGCLKaazWZDyhUmRyK8GbCYZfsnzL4xnTALv6a7pLyYexeLed92pGypyVcCcDz3LiG6sa1wLwuRo28U77ArMZuJWXmvbCY8a43F1M9pKNCGQ1tGkBuJXy5aGYabdnkXMMyqCYbYmedeZne8T15ErKkpks7bBLrEmnNb2VFbwW7ga8WmU7BE59sjrq1nCsXV8kFzFNrksEbVyqhda8AuEG3CgPqtoPFMowZLowxfEvj2vHEMmc5rDhSksmVhftCYh9oiABjah5R9x4xxNjg6ynLAimXQWspiMbihBsUpDaixWaqRTvfwUiAaJ48hivNWuMudZyuXKYbF1Y8GPCHv3E6q1cpNatY6J29doqgkmVj7UpF4bJKHt5GJsiwmxhnkBH4gNXo6UJkAVQ5qR2Dyc6ZwN515eUpH7aPsnX6k2eiQHCNoQ2VHubBK2Aq75C2bpmyHh5uCjYFRhiSWcfn62ZKWmnL72qsGrKTbxuV7hP7UVX3DnpZyfuj6N5Bnrrs9gzHCQEs7Uoh97M4GpiDQNAA2HNFkxTxFdyhh1Tab88yu7wXigrWAGDNPnqMk9vXWuu5pAk5FRFPndBjqznGLL1VMyrmxeFVejfyDdavwUq9zhSgdw3RvmFP21A16J4RqFWfXijsyGZDyxdwL9pphK5NDBCfSfWFNsFBCUqLpJjWaC6XVE9DKuePuZnKMEVRivzMm5Tz1Mz8GWvx8C9iLGDZ1GEgxm7onqy1nNC5pzR45xgbXYSP1CtwosrXHas875PzVmZEQx5QDy3s5RBeQhf453E2FpXLC6u3NshkKM6yoXvR1qPfexE6Rkhe768SnhCxTCVmNXehDxfm4LjfMga1ZxvyUAp5x32v1iHe7myaxR7b6rksMw1mpZgpW28sNqXigUwb6uqjTAuMQ57nCgvwi53ZXdEcpjwpecJmAGqWSC2PLzboFVR8noE8QQ16Nypgm5W3ovrqqTHJZYW3xnxRF2yYu21PZ2vxgvX6Wtv7sPewwffEJP7YABFvCAwsJA2X8J6hwuJDu3bwJHMgXEWiwaB5czc2cBsjmQBrwgNH7ypE38i6dACtC1gM8cW4FXMMf1Fv9A7rM7mFoArmHNs15hco1btwKBgKmxAt7Ng27DR7cL5F2uSTXgdHwTUHrBKik3LQa3n7rm1jwYHm9BeHeZJUvrpxbEBYJd1G93d7CVw2sBU3Y5tvtMzVkPHUfAVRZ", false)

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
