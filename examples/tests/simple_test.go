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
		DoubleType:   42,
		FloatType:    42,
		Int32Type:    42,
		Int64Type:    42,
		Uint32Type:   42,
		Uint64Type:   42,
		Sint32Type:   42,
		Sint64Type:   42,
		Fixed32Type:  42,
		Fixed64Type:  42,
		Sfixed32Type: 42,
		Sfixed64Type: 42,
		BoolType:     true,
		StringType:   "string",
		BytesType:    []byte("bytes"),
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	// now change the original
	original.Reset()

	// and check that doppelganger was unchanged
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())
	assert.Equal(t, float64(42), doppelganger.DoubleType)
	assert.Equal(t, float32(42), doppelganger.FloatType)
	assert.Equal(t, int32(42), doppelganger.Int32Type)
	assert.Equal(t, int64(42), doppelganger.Int64Type)
	assert.Equal(t, uint32(42), doppelganger.Uint32Type)
	assert.Equal(t, uint64(42), doppelganger.Uint64Type)
	assert.Equal(t, int32(42), doppelganger.Sint32Type)
	assert.Equal(t, int64(42), doppelganger.Sint64Type)
	assert.Equal(t, uint32(42), doppelganger.Fixed32Type)
	assert.Equal(t, uint64(42), doppelganger.Fixed64Type)
	assert.Equal(t, int32(42), doppelganger.Sfixed32Type)
	assert.Equal(t, int64(42), doppelganger.Sfixed64Type)
	assert.True(t, doppelganger.BoolType)
	assert.Equal(t, "string", doppelganger.StringType)
	assert.Equal(t, []byte("bytes"), doppelganger.BytesType)
}

func TestOptionalsDeepcopy(t *testing.T) {

	var d = float64(42)
	var f = float32(42)
	var i32 = int32(42)
	var i64 = int64(42)
	var ui32 = uint32(42)
	var ui64 = uint64(42)
	var si32 = int32(42)
	var si64 = int64(42)
	var fi32 = uint32(42)
	var fi64 = uint64(42)
	var sfi32 = int32(42)
	var sfi64 = int64(42)
	var b = true
	s := "the answer to life the universe and everything"
	bts := []byte("the answer")

	original := &protos.ABitOfOptionals{
		DoubleType:   &d,
		FloatType:    &f,
		Int32Type:    &i32,
		Int64Type:    &i64,
		Uint32Type:   &ui32,
		Uint64Type:   &ui64,
		Sint32Type:   &si32,
		Sint64Type:   &si64,
		Fixed32Type:  &fi32,
		Fixed64Type:  &fi64,
		Sfixed32Type: &sfi32,
		Sfixed64Type: &sfi64,
		BoolType:     &b,
		StringType:   &s,
		BytesType:    bts,
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	d++
	f++
	i32++
	i64++
	ui32++
	ui64++
	si32++
	si64++
	fi32++
	fi64++
	sfi32++
	sfi64++
	b = false
	s = "42 " + s
	bts = []byte(s)

	assert.Equal(t, float64(42), *doppelganger.DoubleType)
	assert.Equal(t, float32(42), *doppelganger.FloatType)
	assert.Equal(t, int32(42), *doppelganger.Int32Type)
	assert.Equal(t, int64(42), *doppelganger.Int64Type)
	assert.Equal(t, uint32(42), *doppelganger.Uint32Type)
	assert.Equal(t, uint64(42), *doppelganger.Uint64Type)
	assert.Equal(t, int32(42), *doppelganger.Sint32Type)
	assert.Equal(t, int64(42), *doppelganger.Sint64Type)
	assert.Equal(t, uint32(42), *doppelganger.Fixed32Type)
	assert.Equal(t, uint64(42), *doppelganger.Fixed64Type)
	assert.Equal(t, int32(42), *doppelganger.Sfixed32Type)
	assert.Equal(t, int64(42), *doppelganger.Sfixed64Type)
	assert.True(t, *doppelganger.BoolType)
	assert.Equal(t, "the answer to life the universe and everything", *doppelganger.StringType)
	assert.Equal(t, "the answer", string(doppelganger.BytesType))
}
