package mempool

import (
	"testing"
)

func TestPriorityQueueNewEmpty(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Len() != 0 {
		t.Errorf("new queue should be empty, got len %d", pq.Len())
	}
}

func TestPriorityQueuePushAndPop(t *testing.T) {
	pq := NewPriorityQueue()
	tx := NewTransaction("0xAlice", 0, 50, 100)
	pq.Push(tx)

	if pq.Len() != 1 {
		t.Fatalf("expected len 1, got %d", pq.Len())
	}

	popped := pq.Pop()
	if popped == nil {
		t.Fatal("Pop returned nil")
	}
	if popped.Hash != tx.Hash {
		t.Errorf("popped wrong tx: got %s, want %s", popped.Hash, tx.Hash)
	}
	if pq.Len() != 0 {
		t.Errorf("queue should be empty after pop, got len %d", pq.Len())
	}
}

func TestPriorityQueuePopEmpty(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Pop() != nil {
		t.Error("Pop on empty queue should return nil")
	}
}

func TestPriorityQueuePeekEmpty(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.Peek() != nil {
		t.Error("Peek on empty queue should return nil")
	}
}

func TestPriorityQueueRemoveLowestEmpty(t *testing.T) {
	pq := NewPriorityQueue()
	if pq.RemoveLowest() != nil {
		t.Error("RemoveLowest on empty queue should return nil")
	}
}

func TestPriorityQueueOrdering(t *testing.T) {
	pq := NewPriorityQueue()

	// Push transactions with varying gas prices.
	prices := []uint64{10, 50, 30, 70, 20}
	for i, price := range prices {
		tx := NewTransaction("0xAlice", uint64(i), price, 100)
		pq.Push(tx)
	}

	// Pop should return highest gas price first.
	expectedOrder := []uint64{70, 50, 30, 20, 10}
	for i, expected := range expectedOrder {
		tx := pq.Pop()
		if tx == nil {
			t.Fatalf("Pop returned nil at index %d", i)
		}
		if tx.GasPrice != expected {
			t.Errorf("pop %d: got gasPrice %d, want %d", i, tx.GasPrice, expected)
		}
	}
}

func TestPriorityQueuePeekDoesNotRemove(t *testing.T) {
	pq := NewPriorityQueue()
	tx := NewTransaction("0xAlice", 0, 100, 100)
	pq.Push(tx)

	peeked := pq.Peek()
	if peeked == nil {
		t.Fatal("Peek returned nil")
	}
	if peeked.Hash != tx.Hash {
		t.Errorf("Peek returned wrong tx")
	}
	if pq.Len() != 1 {
		t.Errorf("Peek should not remove element, len is %d", pq.Len())
	}
}

func TestPriorityQueuePeekReturnsHighest(t *testing.T) {
	pq := NewPriorityQueue()

	tx1 := NewTransaction("0xAlice", 0, 10, 100)
	tx2 := NewTransaction("0xBob", 0, 50, 100)
	tx3 := NewTransaction("0xCharlie", 0, 30, 100)

	pq.Push(tx1)
	pq.Push(tx2)
	pq.Push(tx3)

	top := pq.Peek()
	if top.GasPrice != 50 {
		t.Errorf("Peek should return highest gas price tx, got %d", top.GasPrice)
	}
}

func TestPriorityQueueRemoveLowest(t *testing.T) {
	pq := NewPriorityQueue()

	tx1 := NewTransaction("0xAlice", 0, 10, 100)
	tx2 := NewTransaction("0xBob", 0, 50, 100)
	tx3 := NewTransaction("0xCharlie", 0, 30, 100)

	pq.Push(tx1)
	pq.Push(tx2)
	pq.Push(tx3)

	lowest := pq.RemoveLowest()
	if lowest == nil {
		t.Fatal("RemoveLowest returned nil")
	}
	if lowest.GasPrice != 10 {
		t.Errorf("RemoveLowest should return gas price 10, got %d", lowest.GasPrice)
	}
	if pq.Len() != 2 {
		t.Errorf("expected len 2 after removal, got %d", pq.Len())
	}
}

func TestPriorityQueueRemoveLowestTieBreaking(t *testing.T) {
	pq := NewPriorityQueue()

	// Two transactions with the same gas price. The one with the later timestamp
	// should be evicted first (newest among equally cheap).
	tx1 := &Transaction{Hash: "0x1", Sender: "0xA", Nonce: 0, GasPrice: 10, Size: 100, Timestamp: 1000}
	tx2 := &Transaction{Hash: "0x2", Sender: "0xB", Nonce: 0, GasPrice: 10, Size: 100, Timestamp: 2000}

	pq.Push(tx1)
	pq.Push(tx2)

	lowest := pq.RemoveLowest()
	if lowest.Timestamp != 2000 {
		t.Errorf("expected later timestamp to be evicted first, got timestamp %d", lowest.Timestamp)
	}
}

func TestPriorityQueueAll(t *testing.T) {
	pq := NewPriorityQueue()

	tx1 := NewTransaction("0xAlice", 0, 10, 100)
	tx2 := NewTransaction("0xBob", 0, 50, 100)
	tx3 := NewTransaction("0xCharlie", 0, 30, 100)

	pq.Push(tx1)
	pq.Push(tx2)
	pq.Push(tx3)

	all := pq.All()
	if len(all) != 3 {
		t.Fatalf("All should return 3 items, got %d", len(all))
	}

	// Verify All returns a copy, not the internal slice.
	all[0] = nil
	if pq.Len() != 3 {
		t.Error("modifying All result should not affect the queue")
	}
	// Verify the queue is still intact.
	if pq.Peek() == nil {
		t.Error("queue should still have elements after modifying All result")
	}
}

func TestPriorityQueueAllEmpty(t *testing.T) {
	pq := NewPriorityQueue()
	all := pq.All()
	if len(all) != 0 {
		t.Errorf("All on empty queue should return empty slice, got %d", len(all))
	}
}

func TestPriorityQueueSingleElement(t *testing.T) {
	pq := NewPriorityQueue()
	tx := NewTransaction("0xAlice", 0, 42, 100)
	pq.Push(tx)

	// RemoveLowest on single element should return that element.
	removed := pq.RemoveLowest()
	if removed == nil {
		t.Fatal("RemoveLowest returned nil")
	}
	if removed.GasPrice != 42 {
		t.Errorf("expected gas price 42, got %d", removed.GasPrice)
	}
	if pq.Len() != 0 {
		t.Errorf("queue should be empty, got len %d", pq.Len())
	}
}

func TestPriorityQueueManyElements(t *testing.T) {
	pq := NewPriorityQueue()

	// Insert 100 elements with varying gas prices.
	for i := 0; i < 100; i++ {
		tx := NewTransaction("0xSender", uint64(i), uint64(i*3+1), 100)
		pq.Push(tx)
	}

	if pq.Len() != 100 {
		t.Fatalf("expected len 100, got %d", pq.Len())
	}

	// Pop all and verify they come out in descending gas price order.
	prev := uint64(999999)
	for i := 0; i < 100; i++ {
		tx := pq.Pop()
		if tx == nil {
			t.Fatalf("Pop returned nil at index %d", i)
		}
		if tx.GasPrice > prev {
			t.Errorf("out of order at index %d: got %d after %d", i, tx.GasPrice, prev)
		}
		prev = tx.GasPrice
	}
}

func TestPriorityQueueTieBreakByTimestamp(t *testing.T) {
	pq := NewPriorityQueue()

	// Two transactions with the same gas price but different timestamps.
	// The earlier timestamp should have higher priority (pop first).
	tx1 := &Transaction{Hash: "0xearly", Sender: "0xA", Nonce: 0, GasPrice: 50, Size: 100, Timestamp: 1000}
	tx2 := &Transaction{Hash: "0xlate", Sender: "0xB", Nonce: 0, GasPrice: 50, Size: 100, Timestamp: 2000}

	pq.Push(tx2) // push later one first
	pq.Push(tx1) // push earlier one second

	popped := pq.Pop()
	if popped.Timestamp != 1000 {
		t.Errorf("expected earlier timestamp to have higher priority, got timestamp %d", popped.Timestamp)
	}
}

func TestPriorityQueueInterleavedPushPop(t *testing.T) {
	pq := NewPriorityQueue()

	tx1 := NewTransaction("0xA", 0, 10, 100)
	tx2 := NewTransaction("0xB", 0, 30, 100)
	tx3 := NewTransaction("0xC", 0, 20, 100)

	pq.Push(tx1)
	pq.Push(tx2)

	// Pop highest (30).
	top := pq.Pop()
	if top.GasPrice != 30 {
		t.Errorf("expected 30, got %d", top.GasPrice)
	}

	// Push 20, now queue has 10 and 20.
	pq.Push(tx3)

	top = pq.Pop()
	if top.GasPrice != 20 {
		t.Errorf("expected 20, got %d", top.GasPrice)
	}

	top = pq.Pop()
	if top.GasPrice != 10 {
		t.Errorf("expected 10, got %d", top.GasPrice)
	}

	if pq.Pop() != nil {
		t.Error("queue should be empty")
	}
}
