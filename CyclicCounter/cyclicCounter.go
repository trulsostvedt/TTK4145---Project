package cyclicCounter

import (
	"sync"
	"time"
)

// Enum for the states of the orders
const (
	NO_ORDER    = 0
	UNCONFIRMED = 1
	CONFIRMED   = 2
)

// Orderstatus is a struct that holds the status of an order
type Orderstatus struct {
	Floor     int
	State     int
	Timestamp int64
}

// CyclicCounter administers the states of the orders
type CyclicCounter struct {
	Orders map[int]Orderstatus // Map of orders
	mu     sync.Mutex          // Mutex for the map
}

// NewCyclicCounter creates a new CyclicCounter instance
func NewCyclicCounter() *CyclicCounter {
	return &CyclicCounter{Orders: make(map[int]Orderstatus)}
}

// RegisterOrder registers a new unconfirmed order
func (cc *CyclicCounter) RegisterOrder(floor int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// If the order is already registered, update the timestamp
	if _, exists := cc.Orders[floor]; exists {
		cc.Orders[floor] = Orderstatus{
			Floor:     floor,
			State:     UNCONFIRMED,
			Timestamp: time.Now().Unix(),
		}
	}

}

// ConfirmOrder changes the state of an order from uncomfirmed to confirmed
func (cc *CyclicCounter) ConfirmOrder(floor int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// If the order is already registered, update the state
	if _, exists := cc.Orders[floor]; exists && order.State == UNCONFIRMED {
		order.State = CONFIRMED
		cc.Orders[floor] = order
	}
}

// RemoveOrder removes an order from the map
func (cc *CyclicCounter) RemoveOrder(floor int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	delete(cc.Orders, floor)
}
