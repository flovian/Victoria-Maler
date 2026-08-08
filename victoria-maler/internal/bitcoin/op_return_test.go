package bitcoin

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestBuildOpReturnScriptFraming(t *testing.T) {
	script, err := BuildOpReturnScript("abc123")
	if err != nil {
		t.Fatal(err)
	}

	raw, err := hex.DecodeString(script)
	if err != nil {
		t.Fatal(err)
	}

	if raw[0] != opReturn {
		t.Fatalf("expected OP_RETURN (0x6a), got %#x", raw[0])
	}
	payload := []byte(Prefix + "abc123")
	if int(raw[1]) != len(payload) {
		t.Fatalf("expected push length %d, got %d", len(payload), raw[1])
	}
	if string(raw[2:]) != string(payload) {
		t.Fatal("payload mismatch")
	}
}

func TestBuildOpReturnIncludesPrefix(t *testing.T) {
	script, _ := BuildOpReturnScript("c0ffee")
	if !strings.Contains(script, hex.EncodeToString([]byte(Prefix))) {
		t.Fatal("OP_RETURN payload must carry the ECOVIC| prefix")
	}
}

func TestSimulatedTxIDPrefix(t *testing.T) {
	id := SimulatedTxID("aabbccdd11223344")
	if !strings.HasPrefix(id, "sim-") {
		t.Fatalf("expected sim- prefix, got %s", id)
	}
}

func TestVerifyHashSimulated(t *testing.T) {
	hash := "deadbeefdeadbeefdeadbeefdeadbeef"
	result := VerifyHash(hash, SimulatedTxID(hash))
	if !result.Verified {
		t.Fatal("simulated anchors should verify against stored hash")
	}
	if result.Method != "simulated" {
		t.Fatalf("expected simulated method, got %s", result.Method)
	}
}
