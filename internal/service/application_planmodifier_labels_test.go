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

func TestNormalize_whitespaceOnly(t *testing.T) {
	got := normalize("  \n  \n")
	if len(got) != 0 {
		t.Errorf("normalize(whitespace) = %v, want empty slice", got)
	}
}

func TestNormalize_plaintextSorts(t *testing.T) {
	got := normalize("k2=v2\nk1=v1")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != len(want) {
		t.Fatalf("normalize len = %d, want %d (got=%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("normalize[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNormalize_plaintextTrimsPerLine(t *testing.T) {
	got := normalize("  k1=v1  \n\t k2=v2\t \n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v", got, want)
	}
}

func TestNormalize_dropsEmptyLines(t *testing.T) {
	got := normalize("k1=v1\n\n\nk2=v2\n\n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v", got, want)
	}
}

func TestNormalize_mixedLineEndings(t *testing.T) {
	got := normalize("k1=v1\r\nk2=v2\r\n")
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize = %v, want %v (handles CRLF)", got, want)
	}
}
