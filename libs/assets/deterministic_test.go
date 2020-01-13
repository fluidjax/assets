package assets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		w, _ := NewTrusteeGroup(i)
		expression := "t1 + t2 + t3 > 1 "
		participants := &map[string][]byte{
			"t1": idT1.Key(),
			"t2": idT2.Key(),
			"t3": idT3.Key(),
		}
		w.ConfigureTrusteeGroup(expression, participants, testDescription)
		q1, _ := w.serializePayload()
		q2, _ := w.serializePayload()

		resq := bytes.Compare(q1, q2)
		assert.True(t, resq == 0, "Two consecutives serializations yield different results")

	}

}
