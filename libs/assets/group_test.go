package assets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/qredo/assets/libs/protobuffer"
	"github.com/stretchr/testify/assert"
)


func Test_Group(t *testing.T) {

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

	w, err := NewGroup(i, protobuffer.PBGroupType_trusteeGroup)
	expression := "t1 + t2 + t3 > 1 "
	participants := &map[string][]byte{
		"t1": idT1.Key(),
		"t2": idT2.Key(),
		"t3": idT3.Key(),
	}

	w.ConfigureGroup(expression, participants, testDescription)

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
	retrievedGroup := retrieved.Asset.GetGroup()

	assert.Equal(t, testDescription, retrievedGroup.Description, "Load/Save failed")
}

//Tests to ensure objects are determinsitically serialized - something protobuffers by default does not guarantee
//This bug manifests when serializing maps, - their order when serialized is not guaranteed.
//Fix is to use a fork of proto buffers https://github.com/gogo/protobuf
func Test_Determinism(t *testing.T) {
	store := NewMapstore()
	_, idT1, idT2, idT3 := SetupIDDocs(store)
	testName := "ABC!23"
	testDescription := "ZXC#@!"
	i, err := NewIDDoc(testName)
	assert.Nil(t, err, "Error should be nil")

	//Non determinsitic serialization is intermittents, 10 times should be enough to show the error
	//However occasionaly it may not occur on the run of every test, increase iterations if debuggging as issue.
	for j := 0; j < 10; j++ {
		w, _ := NewGroup(i, protobuffer.PBGroupType_trusteeGroup)
		expression := "t1 + t2 + t3 > 1 "
		participants := &map[string][]byte{
			"t1": idT1.Key(),
			"t2": idT2.Key(),
			"t3": idT3.Key(),
		}
		w.ConfigureGroup(expression, participants, testDescription)
		q1, _ := w.serializePayload()
		q2, _ := w.serializePayload()

		resq := bytes.Compare(q1, q2)
		assert.True(t, resq == 0, "Two consecutives serializations yield different results")
	}
}
