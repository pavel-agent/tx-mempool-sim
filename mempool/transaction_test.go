package mempool

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewTransactionFields(t *testing.T) {
	tx := NewTransaction("0xAlice", 7, 100, 256)

	if tx.Sender != "0xAlice" {
		t.Errorf("sender: got %q, want %q", tx.Sender, "0xAlice")
	}
	if tx.Nonce != 7 {
		t.Errorf("nonce: got %d, want 7", tx.Nonce)
	}
	if tx.GasPrice != 100 {
		t.Errorf("gasPrice: got %d, want 100", tx.GasPrice)
	}
	if tx.Size != 256 {
		t.Errorf("size: got %d, want 256", tx.Size)
	}
	if tx.Timestamp == 0 {
		t.Error("timestamp should be set to a non-zero value")
	}
}

func TestNewTransactionHashFormat(t *testing.T) {
	tx := NewTransaction("0xAlice", 0, 50, 200)

	if !strings.HasPrefix(tx.Hash, "0x") {
		t.Errorf("hash should start with 0x, got %q", tx.Hash)
	}
	// SHA3-256 produces 32 bytes = 64 hex chars, plus "0x" prefix = 66 chars.
	if len(tx.Hash) != 66 {
		t.Errorf("hash length: got %d, want 66", len(tx.Hash))
	}
}

func TestNewTransactionHashDeterministic(t *testing.T) {
	// Two transactions with the same fields and same timestamp should have the same hash.
	tx := NewTransaction("0xAlice", 0, 50, 200)
	// Manually recompute to verify determinism.
	recomputed := tx.computeHash()
	if tx.Hash != recomputed {
		t.Errorf("hash is not deterministic: %s vs %s", tx.Hash, recomputed)
	}
}

func TestNewTransactionHashUniqueness(t *testing.T) {
	// Different inputs should produce different hashes.
	tx1 := NewTransaction("0xAlice", 0, 50, 200)
	tx2 := NewTransaction("0xBob", 0, 50, 200)

	if tx1.Hash == tx2.Hash {
		t.Error("different senders should produce different hashes")
	}

	tx3 := NewTransaction("0xAlice", 1, 50, 200)
	if tx1.Hash == tx3.Hash {
		t.Error("different nonces should produce different hashes")
	}

	tx4 := NewTransaction("0xAlice", 0, 100, 200)
	if tx1.Hash == tx4.Hash {
		t.Error("different gas prices should produce different hashes")
	}

	tx5 := NewTransaction("0xAlice", 0, 50, 300)
	if tx1.Hash == tx5.Hash {
		t.Error("different sizes should produce different hashes")
	}
}

func TestTransactionMarshalJSON(t *testing.T) {
	tx := NewTransaction("0xAlice", 3, 75, 128)

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify all expected fields are present.
	expectedFields := []string{"hash", "sender", "nonce", "gasPrice", "size", "timestamp"}
	for _, field := range expectedFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("missing field %q in JSON output", field)
		}
	}

	if parsed["sender"] != "0xAlice" {
		t.Errorf("sender in JSON: got %v, want 0xAlice", parsed["sender"])
	}
	if parsed["nonce"].(float64) != 3 {
		t.Errorf("nonce in JSON: got %v, want 3", parsed["nonce"])
	}
	if parsed["gasPrice"].(float64) != 75 {
		t.Errorf("gasPrice in JSON: got %v, want 75", parsed["gasPrice"])
	}
}

func TestTransactionMarshalJSONRoundTrip(t *testing.T) {
	tx := NewTransaction("0xAlice", 5, 200, 512)

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var decoded Transaction
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if decoded.Hash != tx.Hash {
		t.Errorf("hash mismatch after round-trip: got %s, want %s", decoded.Hash, tx.Hash)
	}
	if decoded.Sender != tx.Sender {
		t.Errorf("sender mismatch after round-trip")
	}
	if decoded.Nonce != tx.Nonce {
		t.Errorf("nonce mismatch after round-trip")
	}
	if decoded.GasPrice != tx.GasPrice {
		t.Errorf("gasPrice mismatch after round-trip")
	}
	if decoded.Size != tx.Size {
		t.Errorf("size mismatch after round-trip")
	}
	if decoded.Timestamp != tx.Timestamp {
		t.Errorf("timestamp mismatch after round-trip")
	}
}

func TestNewTransactionZeroValues(t *testing.T) {
	tx := NewTransaction("", 0, 0, 0)

	if tx.Hash == "" {
		t.Error("hash should be computed even with zero-value fields")
	}
	if tx.Timestamp == 0 {
		t.Error("timestamp should always be set")
	}
}
