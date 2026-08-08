package handlers

import (
	"net/http"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) BitcoinVerifyPage(w http.ResponseWriter, r *http.Request) {
	records, err := h.BitcoinRecords.List()
	if err != nil {
		utils.InternalError(w, "could not load verification records")
		return
	}
	h.Render.Page(w, r, "bitcoin_verify", map[string]interface{}{
		"Records": records,
	})
}

type verifyRequest struct {
	Hash string `json:"hash"`
	TxID string `json:"txid"`
}

func (h *HandlerSet) VerifyHash(w http.ResponseWriter, r *http.Request) {
	var req verifyRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}
	if req.Hash == "" {
		utils.BadRequest(w, "hash is required")
		return
	}

	result := bitcoin.VerifyHash(req.Hash, req.TxID)
	utils.OK(w, result)
}

func (h *HandlerSet) VerifyByHashOnly(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		utils.BadRequest(w, "hash parameter required")
		return
	}

	record, err := h.BitcoinRecords.FindByHash(hash)
	if err != nil {
		utils.NotFound(w, "no anchored record found for this hash")
		return
	}

	result := bitcoin.VerifyHash(record.ContentHash, record.TxID)
	utils.OK(w, result)
}

func (h *HandlerSet) VerificationRecords(w http.ResponseWriter, r *http.Request) {
	records, err := h.BitcoinRecords.List()
	if err != nil {
		utils.InternalError(w, "could not load verification records")
		return
	}
	utils.OK(w, records)
}

func (h *HandlerSet) NodeStatus(w http.ResponseWriter, r *http.Request) {
	enabled := h.Config.Bitcoin.Enabled
	healthy := false
	if h.Wallet != nil && h.Wallet.Client != nil {
		healthy = h.Wallet.Client.IsHealthy()
	}
	height, _ := h.Wallet.BlockchainHeight()
	balance, _ := h.Wallet.Balance()
	utils.OK(w, map[string]interface{}{
		"enabled": enabled,
		"healthy": healthy,
		"network": h.Config.Bitcoin.Network,
		"height":  height,
		"balance": balance,
	})
}
