package hyperliquid

import (
	"encoding/json"
	"strings"
	"testing"
)

// Regression: plain action acks have no response.data; both envelope shapes
// must decode and classify correctly.
func TestActionAckEnvelope(t *testing.T) {
	var ok actionAck
	if err := json.Unmarshal([]byte(`{"status":"ok","response":{"type":"default"}}`), &ok); err != nil {
		t.Fatalf("ok ack decode: %v", err)
	}
	if err := ok.rejection("fallback"); err != nil {
		t.Fatalf("ok ack must not reject: %v", err)
	}

	var rej actionAck
	if err := json.Unmarshal([]byte(`{"status":"err","response":"Insufficient margin to place order."}`), &rej); err != nil {
		t.Fatalf("err ack decode: %v", err)
	}
	err := rej.rejection("fallback")
	if err == nil || !strings.Contains(err.Error(), "Insufficient margin") {
		t.Fatalf("err ack must surface the exchange message, got %v", err)
	}

	var empty actionAck
	if err := json.Unmarshal([]byte(`{"status":"err"}`), &empty); err != nil {
		t.Fatalf("empty err ack decode: %v", err)
	}
	if err := empty.rejection("fallback"); err == nil || err.Error() != "fallback" {
		t.Fatalf("empty err ack must use fallback, got %v", err)
	}
}
