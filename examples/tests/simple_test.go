package tests

import (
	"github.com/dgodyna/protoc-gen-resource/examples/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/testing/protocmp"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestSimpleDeepCopy(t *testing.T) {
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

func TestOptionalsDeepCopy(t *testing.T) {

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

func TestEnumsDeepCopy(t *testing.T) {

	original := &protos.ABitOfEnums{
		EngineType:  protos.ABitOfEnums_ENGINE_TYPE_GAS,
		VehicleType: protos.VehicleType_VEHICLE_TYPE_CAR,
	}
	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	original.EngineType = protos.ABitOfEnums_ENGINE_TYPE_DIESEL
	original.VehicleType = protos.VehicleType_VEHICLE_TYPE_PLANE

	assert.Equal(t, protos.ABitOfEnums_ENGINE_TYPE_GAS, doppelganger.EngineType)
	assert.Equal(t, protos.VehicleType_VEHICLE_TYPE_CAR, doppelganger.VehicleType)
}

func TestMessagesDeepCopy(t *testing.T) {
	original := &protos.ABitOfMessages{
		First: &protos.AnotherM{
			F1: "first",
			F2: "second",
		},
		Second: &protos.ABitOfMessages_Sub{
			I1: 42,
			I2: 42,
		},
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	original.First = &protos.AnotherM{
		F1: "third",
		F2: "fourth",
	}
	original.Second = nil

	assert.Equal(t, &protos.AnotherM{
		F1: "first",
		F2: "second",
	}, doppelganger.First, protocmp.Transform())
	assert.Equal(t, &protos.ABitOfMessages_Sub{
		I1: 42,
		I2: 42,
	}, doppelganger.Second, protocmp.Transform())
}

func TestRepeatedScalarsDeepCopy(t *testing.T) {
	original := &protos.ABitOfRepeatedScalars{
		DoubleType:   []float64{42, 42},
		FloatType:    []float32{42, 42},
		Int32Type:    []int32{42, 42},
		Int64Type:    []int64{42, 42},
		Uint32Type:   []uint32{42, 42},
		Uint64Type:   []uint64{42, 42},
		Sint32Type:   []int32{42, 42},
		Sint64Type:   []int64{42, 42},
		Fixed32Type:  []uint32{42, 42},
		Fixed64Type:  []uint64{42, 42},
		Sfixed32Type: []int32{42, 42},
		Sfixed64Type: []int64{42, 42},
		BoolType:     []bool{true, true, false},
		StringType:   []string{"a", "bit", "of", "every", "thing"},
		BytesType:    [][]byte{[]byte("a bit of"), []byte("everything")},
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	original.Reset()
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())

	// and check that doppelganger was unchanged
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())
	assert.Equal(t, []float64{42, 42}, doppelganger.DoubleType)
	assert.Equal(t, []float32{42, 42}, doppelganger.FloatType)
	assert.Equal(t, []int32{42, 42}, doppelganger.Int32Type)
	assert.Equal(t, []int64{42, 42}, doppelganger.Int64Type)
	assert.Equal(t, []uint32{42, 42}, doppelganger.Uint32Type)
	assert.Equal(t, []uint64{42, 42}, doppelganger.Uint64Type)
	assert.Equal(t, []int32{42, 42}, doppelganger.Sint32Type)
	assert.Equal(t, []int64{42, 42}, doppelganger.Sint64Type)
	assert.Equal(t, []uint32{42, 42}, doppelganger.Fixed32Type)
	assert.Equal(t, []uint64{42, 42}, doppelganger.Fixed64Type)
	assert.Equal(t, []int32{42, 42}, doppelganger.Sfixed32Type)
	assert.Equal(t, []int64{42, 42}, doppelganger.Sfixed64Type)
	assert.Equal(t, []bool{true, true, false}, doppelganger.BoolType)
	assert.Equal(t, []string{"a", "bit", "of", "every", "thing"}, doppelganger.StringType)
	assert.Equal(t, [][]byte{[]byte("a bit of"), []byte("everything")}, doppelganger.BytesType)
}

func TestRepeatedEnumsDeepCopy(t *testing.T) {
	original := &protos.ABitOfRepeatedEnums{
		EngineType: []protos.ABitOfRepeatedEnums_EngineType{protos.ABitOfRepeatedEnums_ENGINE_TYPE_GAS, protos.ABitOfRepeatedEnums_ENGINE_TYPE_DIESEL},
	}
	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	original.Reset()
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())

	assert.Equal(t, []protos.ABitOfRepeatedEnums_EngineType{protos.ABitOfRepeatedEnums_ENGINE_TYPE_GAS, protos.ABitOfRepeatedEnums_ENGINE_TYPE_DIESEL}, doppelganger.EngineType)

}

func TestRepeatedMessagesDeepCopy(t *testing.T) {
	original := &protos.ABitOfRepeatedMessages{
		First: []*protos.ABitOfRepeatedMessages_RepeatedSub{
			{
				I1: 42,
				I2: 42,
			},
			{
				I1: 21,
				I2: 21,
			},
		},
	}

	// check deepcopy itself
	doppelganger := original.DeepCopy()
	assert.Equal(t, original, doppelganger, protocmp.Transform())

	original.Reset()
	assert.NotEqual(t, original, doppelganger, protocmp.Transform())

	assert.Equal(t, []*protos.ABitOfRepeatedMessages_RepeatedSub{
		{
			I1: 42,
			I2: 42,
		},
		{
			I1: 21,
			I2: 21,
		},
	}, doppelganger.First)

}
