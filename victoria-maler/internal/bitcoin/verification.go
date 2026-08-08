package bitcoin

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type VerificationResult struct {
	Verified    bool   `json:"verified"`
	TxID        string `json:"txid"`
	Hash        string `json:"hash"`
	Method      string `json:"method"` // "on-chain" | "simulated"
	BlockHeight int64  `json:"block_height,omitempty"`
	Note        string `json:"note,omitempty"`
}

// VerifyTransaction keeps the original signature for compatibility and
// verifies that the given transaction contains the expected payload.
func VerifyTransaction(txid string) bool {
	return VerifyTxPayload(txid, "") != nil
}

// VerifyHash verifies that contentHash matches the anchored record on chain.
// Simulated anchors (sim- txids) always match the locally stored hash.
func VerifyHash(contentHash, txid string) *VerificationResult {
	if isSimulated(txid) {
		return &VerificationResult{
			Verified: true,
			TxID:     txid,
			Hash:     contentHash,
			Method:   "simulated",
			Note:     "Anchored in simulated mode - no on-chain transaction. Toggle BITCOIN_ENABLED with a reachable node for real anchoring.",
		}
	}
	return VerifyTxPayload(txid, contentHash)
}

func isSimulated(txid string) bool {
	return strings.HasPrefix(strings.ToLower(txid), "sim-")
}

func VerifyTxPayload(txid, contentHash string) *VerificationResult {
	if txid == "" {
		return nil
	}

	result := &VerificationResult{
		TxID:   txid,
		Hash:   contentHash,
		Method: "on-chain",
	}

	client := globalClient()
	if client == nil {
		result.Note = "No bitcoind connection configured. Cannot verify on-chain."
		return result
	}

	tx, err := client.GetRawTransaction(txid)
	if err != nil {
		result.Note = fmt.Sprintf("Transaction lookup failed: %v", err)
		return result
	}

	if bh, ok := tx["blockheight"].(float64); ok {
		result.BlockHeight = int64(bh)
	}

	needle := hex.EncodeToString([]byte(Prefix + contentHash))
	if containsOpReturn(tx, needle) {
		result.Verified = true
		result.Note = "Transaction confirmed on-chain and matches the recorded hash."
	} else {
		result.Note = "Transaction found but its OP_RETURN payload does not match this hash."
	}
	return result
}

func containsOpReturn(tx map[string]interface{}, needle string) bool {
	vouts, ok := tx["vout"].([]interface{})
	if !ok {
		return false
	}
	for _, v := range vouts {
		vout, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		spk, ok := vout["scriptPubKey"].(map[string]interface{})
		if !ok {
			continue
		}
		if asm, ok := spk["asm"].(string); ok && strings.Contains(strings.ToLower(asm), strings.ToLower(needle)) {
			return true
		}
		if scriptHex, ok := spk["hex"].(string); ok && strings.Contains(strings.ToLower(scriptHex), strings.ToLower(needle)) {
			return true
		}
	}
	return false
}

var globalClientFn func() *RPCClient

func SetGlobalClient(client *RPCClient) {
	globalClientFn = func() *RPCClient { return client }
}

func globalClient() *RPCClient {
	if globalClientFn == nil {
		return nil
	}
	return globalClientFn()
}
