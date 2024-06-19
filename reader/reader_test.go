package reader

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	input := []byte("aaaa")
	res, _ := Read(bytes.NewReader(input))

	t.Error(res)
}
