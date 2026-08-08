package bitcoin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RPCClient struct {
	URL  string
	User string
	Pass string
	HTTP *http.Client
}

func NewRPCClient(url, user, pass string) *RPCClient {
	return &RPCClient{
		URL:  url,
		User: user,
		Pass: pass,
		HTTP: &http.Client{Timeout: 10 * time.Second},
	}
}

type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error"`
}

func (c *RPCClient) Call(method string, params ...interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(rpcRequest{
		JSONRPC: "1.0",
		ID:      "ecovictoria",
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.User, c.Pass)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var out rpcResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	if out.Error != nil {
		return nil, errors.New(out.Error.Message)
	}
	return out.Result, nil
}

// IsHealthy reports whether a bitcoind node is reachable.
func (c *RPCClient) IsHealthy() bool {
	_, err := c.Call("getblockchaininfo")
	return err == nil
}

// GetNewAddress requests a fresh address from the node wallet.
func (c *RPCClient) GetNewAddress() (string, error) {
	raw, err := c.Call("getnewaddress")
	if err != nil {
		return "", err
	}
	var addr string
	if err := json.Unmarshal(raw, &addr); err != nil {
		return "", err
	}
	return addr, nil
}

// SendOpReturn creates, funds, signs and broadcasts an OP_RETURN transaction
// carrying the given payload. It returns the transaction id.
func (c *RPCClient) SendOpReturn(payload string) (string, error) {
	script, err := BuildOpReturnScript(payload)
	if err != nil {
		return "", err
	}

	inputs := []interface{}{}
	outputs := map[string]interface{}{
		"data": script,
	}

	raw, err := c.Call("createrawtransaction", inputs, outputs)
	if err != nil {
		return "", fmt.Errorf("createrawtransaction: %w", err)
	}

	var rawHex string
	if err := json.Unmarshal(raw, &rawHex); err != nil {
		return "", fmt.Errorf("decode raw tx: %w", err)
	}

	funded, err := c.Call("fundrawtransaction", rawHex)
	if err != nil {
		return "", fmt.Errorf("fundrawtransaction: %w", err)
	}
	fundedHex, err := rawHexField(funded)
	if err != nil {
		return "", fmt.Errorf("decode funded tx: %w", err)
	}

	signed, err := c.Call("signrawtransactionwithwallet", fundedHex)
	if err != nil {
		return "", fmt.Errorf("signrawtransactionwithwallet: %w", err)
	}
	signedHex, err := rawHexField(signed)
	if err != nil {
		return "", fmt.Errorf("decode signed tx: %w", err)
	}

	sent, err := c.Call("sendrawtransaction", signedHex)
	if err != nil {
		return "", fmt.Errorf("sendrawtransaction: %w", err)
	}

	var txid string
	if err := json.Unmarshal(sent, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// GetRawTransaction fetches a raw transaction by id and decodes it.
func (c *RPCClient) GetRawTransaction(txid string) (map[string]interface{}, error) {
	raw, err := c.Call("getrawtransaction", txid, true)
	if err != nil {
		return nil, err
	}
	var tx map[string]interface{}
	if err := json.Unmarshal(raw, &tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func rawHexField(raw json.RawMessage) (string, error) {
	var v struct {
		Hex string `json:"hex"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	return v.Hex, nil
}
