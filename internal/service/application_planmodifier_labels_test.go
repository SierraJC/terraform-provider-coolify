package service

import (
	"testing"
)

func TestNormalize_empty(t *testing.T) {
	got := normalize("")
	if len(got) != 0 {
		t.Errorf("normalize(\"\") = %v, want empty slice", got)
	}
}
