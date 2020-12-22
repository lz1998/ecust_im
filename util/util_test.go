package util

import "testing"

func TestGenerateId(t *testing.T) {
	t.Log(GenerateId())
	t.Log(GenerateId())
	t.Log(GenerateId())
	t.Log(GenerateId())
}
