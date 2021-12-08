package tests

import (
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/examples/protos"
	"testing"
)

func TestSimple(t *testing.T) {
	s := &protos.ABitOfScalars{}
	fmt.Println(s.GetResourceGroup())
}
