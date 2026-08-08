package bitcoin

import (
	"fmt"
	"strings"
	"time"
)

// AnchoringService anchors content hashes to Bitcoin using OP_RETURN.
// When no bitcoind node is reachable it degrades gracefully to a clearly
// labelled simulated anchor so the full platform flow keeps working.
type AnchoringService struct {
	Enabled    bool
	Client     *RPCClient
	recordRepo RecordSaver
	now        func() int64
}

type RecordSaver interface {
	Save(record *AnchorRecord) error
}

type AnchorRecord struct {
	ContentHash string `json:"content_hash"`
	TxID        string `json:"txid"`
	OpReturn    string `json:"op_return"`
	Status      string `json:"status"`
	RecordType  string `json:"record_type"`
	RecordID    int64  `json:"record_id"`
	CreatedAt   int64  `json:"created_at"`
}

func NewAnchoringService(enabled bool, client *RPCClient, saver RecordSaver) *AnchoringService {
	return &AnchoringService{
		Enabled:    enabled,
		Client:     client,
		recordRepo: saver,
		now:        func() int64 { return time.Now().Unix() },
	}
}

const (
	StatusAnchored  = "anchored"
	StatusSimulated = "simulated"
)

// AnchorHash stores the hash locally and attempts on-chain anchoring.
func (s *AnchoringService) AnchorHash(contentHash, recordType string, recordID int64) (*AnchorRecord, error) {
	record := &AnchorRecord{
		ContentHash: contentHash,
		Status:      StatusSimulated,
		RecordType:  recordType,
		RecordID:    recordID,
		CreatedAt:   s.now(),
	}

	anchored, err := s.tryAnchor(contentHash)
	if err == nil {
		record.Status = StatusAnchored
		record.TxID = anchored
		record.OpReturn = BuildOpReturn(contentHash)
	} else {
		record.TxID = SimulatedTxID(contentHash)
		record.OpReturn = BuildOpReturn(contentHash)
	}

	if s.recordRepo != nil {
		if err := s.recordRepo.Save(record); err != nil {
			return nil, fmt.Errorf("persist anchor record: %w", err)
		}
	}
	return record, nil
}

func (s *AnchoringService) tryAnchor(contentHash string) (string, error) {
	if !s.Enabled || s.Client == nil {
		return "", fmt.Errorf("bitcoin anchoring disabled")
	}
	if !s.Client.IsHealthy() {
		return "", fmt.Errorf("bitcoind node unreachable")
	}
	return s.Client.SendOpReturn(contentHash)
}

// SimulatedTxID produces a deterministic, clearly non-real txid for
// simulated anchors so auditors can tell them apart at a glance.
func SimulatedTxID(contentHash string) string {
	if len(contentHash) >= 16 {
		return "sim-" + strings.ToLower(contentHash[:16])
	}
	return "sim-" + strings.ToLower(contentHash)
}
