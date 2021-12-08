package tests

import (
	"fmt"
	"github.com/dgodyna/protoc-gen-resource/examples/protos"
	"testing"
)

func TestSimple(t *testing.T) {
	var s protos.ABitOfScalars
	fmt.Println(s)
}
