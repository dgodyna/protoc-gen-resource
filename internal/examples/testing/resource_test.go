package testing

import (
	"fmt"
	"github.com/dgodyna/protoc-gen-deepcopy/pkg/internal/examples"
	"gotest.tools/assert"
	"testing"
)

func TestResource(t *testing.T) {

	assert.Equal(t, 5, 5)
	var s examples.ABitOfScalars

	fmt.Println(s)
}
