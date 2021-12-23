package tests

import (
	"github.com/dgodyna/protoc-gen-resource/examples/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestSimpleDeepcopy(t *testing.T) {
	original := &protos.ABitOfScalars{
		DoubleType:   1,
		FloatType:    2,
		Int32Type:    3,
		Int64Type:    4,
		Uint32Type:   5,
		Uint64Type:   6,
		Sint32Type:   7,
		Sint64Type:   8,
		Fixed32Type:  9,
		Fixed64Type:  10,
		Sfixed32Type: 11,
		Sfixed64Type: 12,
		BoolType:     true,
		StringType:   "string",
		BytesType:    []byte("bytes"),
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	// now change the original
	original.Int32Type = 1
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())

}
