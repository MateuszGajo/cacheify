package utils

import "testing"

func TestFindNOccurance(t *testing.T) {
	index := FindNOccurance("abcdefabc", "abc", 2)

	if index != 6 {
		t.Fatalf("Index should be 6 insted we got: %v", index)
	}
}
