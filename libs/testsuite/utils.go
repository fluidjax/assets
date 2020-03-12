package testsuite

import (
	"encoding/hex"
	"testing"

	"github.com/qredo/assets/libs/assets"
	"github.com/qredo/assets/libs/cryptowallet"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)

func SetupIDDocs(t *testing.T) (*assets.IDDoc, *assets.IDDoc, *assets.IDDoc, *assets.IDDoc) {
	go StartWait(4)
	idP := BuildTestIDDoc(t)
	txid, chainErr := nc.PostTx(idP)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	idT1 := BuildTestIDDoc(t)
	txid, chainErr = nc.PostTx(idT1)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	idT2 := BuildTestIDDoc(t)
	txid, chainErr = nc.PostTx(idT2)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	idT3 := BuildTestIDDoc(t)
	txid, chainErr = nc.PostTx(idT3)
	assert.Nil(t, chainErr, "Error should be nil")
	assert.NotEqual(t, txid, "", "TxID should not be empty")

	wg.Wait()

	return idP, idT1, idT2, idT3
}

func BuildTestIDDoc(t *testing.T) *assets.IDDoc {
	randBytes, err := cryptowallet.RandomBytes(32)
	if err != nil {
		panic("Fail to create random auth string")
	}
	i, err := assets.NewIDDoc(hex.EncodeToString(randBytes))
	i.DataStore = nc
	assert.Nil(t, err, "Error should be nil")
	err = i.Sign(i)
	return i
}

func buildTestCreateWallet(t *testing.T, idP *assets.IDDoc, idT1 *assets.IDDoc, idT2 *assets.IDDoc, idT3 *assets.IDDoc) *assets.Wallet {

	//Standard Wallet build
	wallet, err := assets.NewWallet(idP, protobuffer.PBCryptoCurrency_BTC)
	assert.Nil(t, err, "Truth table return should be nil")
	assert.NotNil(t, wallet, "Wallet should not be nil")

	//Error no Transfer Participants
	txid, err := nc.PostTx(wallet)
	assetError, _ := err.(*assets.AssetsError)
	assert.NotNil(t, assetError, "Error should not be nil")
	assert.True(t, txid == "", "TXID should be empty")
	assert.True(t, assetError.Code() == assets.CodeConsensusWalletNoTransferRules, "Incorrect Error code")

	//Add Transfers  - it would now work, but we will break it for testing
	expression := "t1 + t2 + t3 > 1 & p"
	participants := &map[string][]byte{
		"p":  idP.Key(),
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}
	wallet.AddTransfer(protobuffer.PBTransferType_SettlePush, expression, participants, "description")

	sigP, _ := wallet.SignAsset(idP)
	err = wallet.AddSigner(idP, "p", sigP)

	//No Error
	txid, err = nc.PostTx(wallet)
	assert.Nil(t, err, "Error should not nil")
	assert.True(t, txid != "", "TXID should not be empty")

	return wallet
}
