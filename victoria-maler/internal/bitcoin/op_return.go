package bitcoin

import (
	"encoding/hex"
	"fmt"
)

const (
	opReturn    = 0x6a
	opPushData1 = 0x4c
	opPushData2 = 0x4d
)

// Prefix marks EcoChain Victoria records on-chain.
const Prefix = "ECOVIC|"

// BuildOpReturnScript frames data into a standard OP_RETURN script.
// The returned value is the hex-encoded scriptSig/scriptPubKey content
// (opcode + push + payload) accepted by bitcoind's createrawtransaction.
func BuildOpReturnScript(data string) (string, error) {
	payload := []byte(Prefix + data)
	script := []byte{opReturn}

	switch {
	case len(payload) <= 75:
		script = append(script, byte(len(payload)))
	case len(payload) <= 255:
		script = append(script, opPushData1, byte(len(payload)))
	case len(payload) <= 65535:
		script = append(script, opPushData2, byte(len(payload)&0xff), byte(len(payload)>>8))
	default:
		return "", fmt.Errorf("op_return payload too large: %d bytes", len(payload))
	}

	script = append(script, payload...)
	return hex.EncodeToString(script), nil
}

// BuildOpReturn keeps backwards compatibility as a simple framed payload builder.
func BuildOpReturn(data string) string {
	script, _ := BuildOpReturnScript(data)
	return script
}
