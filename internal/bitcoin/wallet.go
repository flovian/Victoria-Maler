package bitcoin

import "encoding/json"

// Wallet is a thin wrapper around the node wallet RPC endpoints.
type Wallet struct {
	Client *RPCClient
}

func NewWallet(client *RPCClient) *Wallet {
	return &Wallet{Client: client}
}

// NewAddress requests a new receive address from the node wallet.
func (w *Wallet) NewAddress() (string, error) {
	if w.Client == nil {
		return "", nil
	}
	return w.Client.GetNewAddress()
}

// Balance returns the total wallet balance in BTC.
func (w *Wallet) Balance() (float64, error) {
	raw, err := w.Client.Call("getbalance")
	if err != nil {
		return 0, err
	}
	var balance float64
	if err := json.Unmarshal(raw, &balance); err != nil {
		return 0, err
	}
	return balance, nil
}

// BlockchainHeight reports the current chain height (0 when unreachable).
func (w *Wallet) BlockchainHeight() (int64, error) {
	raw, err := w.Client.Call("getblockcount")
	if err != nil {
		return 0, err
	}
	var height int64
	if err := json.Unmarshal(raw, &height); err != nil {
		return 0, err
	}
	return height, nil
}
