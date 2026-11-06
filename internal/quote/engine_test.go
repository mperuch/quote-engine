package quote

import (
	"path/filepath"
	"testing"
	"time"
)

func testEngine(t *testing.T) *Engine {
	t.Helper()
	path := filepath.Join("..", "..", "data", "tables.json")
	return New(path, time.Minute)
}

func TestQuoteSouthExpress(t *testing.T) {
	got, err := testEngine(t).Quote(Request{Origin: "SAO", Destination: "CWB", WeightKg: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Quotes) != 1 {
		t.Fatalf("got %d quotes, want 1", len(got.Quotes))
	}
	want := 12.5 + 10*2.4
	if got.Quotes[0].Amount != want {
		t.Fatalf("amount = %v, want %v", got.Quotes[0].Amount, want)
	}
}

func TestQuoteUnknownLane(t *testing.T) {
	_, err := testEngine(t).Quote(Request{Origin: "SAO", Destination: "XXX", WeightKg: 1})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestQuoteRejectsBadWeight(t *testing.T) {
	_, err := testEngine(t).Quote(Request{Origin: "SAO", Destination: "CWB", WeightKg: 0})
	if err == nil {
		t.Fatal("expected error")
	}
}
