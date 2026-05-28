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

func TestIsBase64Labels_validBase64WithEqLines(t *testing.T) {
	b64 := "azE9djEKazI9djI=" // base64("k1=v1\nk2=v2")
	if !isBase64Labels(b64) {
		t.Error("isBase64Labels(b64 of k1=v1\\nk2=v2) = false, want true")
	}
}

func TestIsBase64Labels_validBase64NoEqLines(t *testing.T) {
	b64 := "SGVsbG8=" // base64("Hello") — no '=' lines
	if isBase64Labels(b64) {
		t.Error("isBase64Labels(b64 of Hello) = true, want false (no '=' line)")
	}
}

func TestIsBase64Labels_corrupt(t *testing.T) {
	if isBase64Labels("!!!notb64!!!") {
		t.Error("isBase64Labels(!!!notb64!!!) = true, want false")
	}
}

func TestNormalize_base64InputDecodes(t *testing.T) {
	b64 := "azE9djEKazI9djI=" // base64("k1=v1\nk2=v2")
	got := normalize(b64)
	want := []string{"k1=v1", "k2=v2"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("normalize(b64) = %v, want %v", got, want)
	}
}

func TestNormalize_corruptBase64FallsBackToPlaintext(t *testing.T) {
	got := normalize("!!!notb64!!!")
	want := []string{"!!!notb64!!!"}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("normalize(corrupt) = %v, want %v (plaintext fallback)", got, want)
	}
}

func TestNormalize_base64LookingButNotLabels(t *testing.T) {
	// "SGVsbG8=" = base64("Hello") — decodes OK but no '=' lines
	got := normalize("SGVsbG8=")
	want := []string{"SGVsbG8="}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("normalize(SGVsbG8=) = %v, want %v (treat as plaintext)", got, want)
	}
}
