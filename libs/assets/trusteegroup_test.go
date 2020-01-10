package assets

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TrusteeGroup(t *testing.T) {

	for i := 0; i < 10; i++ {
		fmt.Println(i)
		store := NewMapstore()
		idInitiator, idT1, idT2, idT3 := SetupIDDocs(store)

		_ = idInitiator

		testName := "ABC!23"
		testDescription := "ZXC#@!"
		i, err := NewIDDoc(testName)
		assert.Nil(t, err, "Error should be nil")
		i.Sign()
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

		q1, _ := w.serializePayload()

		w.Sign(i)
		assert.NotNil(t, w.PBSignedAsset.Signature, "Signature is empty")
		res, err := w.Verify(i)

		q2, _ := w.serializePayload()

		resq := bytes.Compare(q1, q2)
		assert.True(t, resq == 0, "Failed to match data")

		assert.Nil(t, err, "Error should be nil")
		assert.True(t, res, "Verify should be true")
		w.Save()

		retrieved, err := Load(i.store, w.Key())
		retrievedTrusteeGroup := retrieved.Asset.GetTrusteeGroup()

		assert.Equal(t, testDescription, retrievedTrusteeGroup.TrusteeGroup.Description, "Load/Save failed")
	}
}
