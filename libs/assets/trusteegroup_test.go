package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func Test_TrusteeGroup(t *testing.T) {

	for i := 0; i < 40; i++ {
		fmt.Println(i)
		store := NewMapstore()
		idInitiator, idT1, idT2, idT3 := SetupIDDocs(store)

		_ = idInitiator

		testName := "ABC!23"
		testDescription := "ZXC#@!"
		i, err := NewIDDoc(testName)
		assert.Nil(t, err, "Error should be nil")
		i.Sign(i)
		i.store = store
		i.Save()

		w, err := NewTrusteeGroup(i)
		expression := "t1 + t2 + t3 > 1 "
		participants := &map[string][]byte{
			"t1": idT1.Key(),
			"t2": idT2.Key(),
			"t3": idT3.Key(),
		}

		w.ConfigureTrusteeGroup(expression, participants, testDescription)

		b := proto.NewBuffer(nil)
		b.SetDeterministic(true)
		b.Marshal(w.Asset)
		msg1, _ := b.DecodeRawBytes(true)
		res1 := sha256.Sum256(msg1)
		fmt.Println(hex.EncodeToString(res1[:]))

		w.Sign(i)

		assert.NotNil(t, w.PBSignedAsset.Signature, "Signature is empty")
		res, err := w.Verify(i)

		c := proto.NewBuffer(nil)
		c.SetDeterministic(true)
		c.Marshal(w.Asset)
		msg2, _ := c.DecodeRawBytes(true)
		res2 := sha256.Sum256(msg2)
		fmt.Println(hex.EncodeToString(res2[:]))

		assert.Nil(t, err, "Error should be nil")
		assert.True(t, res, "Verify should be true")
		w.Save()

		retrieved, err := Load(i.store, w.Key())
		retrievedTrusteeGroup := retrieved.Asset.GetTrusteeGroup()

		assert.Equal(t, testDescription, retrievedTrusteeGroup.TrusteeGroup.Description, "Load/Save failed")
	}
}
