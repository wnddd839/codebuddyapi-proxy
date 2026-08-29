package strutil

import "testing"

func TestFirst(t *testing.T) {
	if got := First(" ", "", "<nil>", " a ", "b"); got != "a" {
		t.Fatalf("First=%q", got)
	}
	if got := First("", " "); got != "" {
		t.Fatalf("empty First=%q", got)
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("  abcdef  ", 3); got != "abc" {
		t.Fatalf("Truncate=%q", got)
	}
	if got := Truncate("ab", 5); got != "ab" {
		t.Fatalf("short Truncate=%q", got)
	}
}

func TestMaskSecret(t *testing.T) {
	if got := MaskSecret("abcdefghijkl", 3); got != "abc...jkl" {
		t.Fatalf("MaskSecret=%q", got)
	}
	if got := MaskSecret("ab", 3); got != "ab..." {
		t.Fatalf("short MaskSecret=%q", got)
	}
}

func TestRandomHex(t *testing.T) {
	a := RandomHex(8)
	b := RandomHex(8)
	if len(a) != 16 || len(b) != 16 {
		t.Fatalf("len a=%d b=%d", len(a), len(b))
	}
	if a == b {
		t.Fatal("expected different random hex")
	}
}
