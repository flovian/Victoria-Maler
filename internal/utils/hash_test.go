package utils

import "testing"

func TestHashHexLength(t *testing.T) {
	hash := HashHex("victoria-maler")
	if len(hash) != 64 {
		t.Fatalf("expected 64 hex chars, got %d", len(hash))
	}
}

func TestHashHexDeterministic(t *testing.T) {
	a := HashHex("cleanup|1|entebbe")
	b := HashHex("cleanup|1|entebbe")
	if a != b {
		t.Fatalf("hash should be deterministic: %s != %s", a, b)
	}
}

func TestHashHexSensitiveToInput(t *testing.T) {
	a := HashHex("cleanup|1|entebbe")
	b := HashHex("cleanup|1|entebbe ")
	if a == b {
		t.Fatal("different inputs should produce different hashes")
	}
}

func TestHashHexBytes(t *testing.T) {
	if HashHexBytes([]byte("file-content")) != HashHex("file-content") {
		t.Fatal("HashHexBytes should match HashHex for same bytes")
	}
}
