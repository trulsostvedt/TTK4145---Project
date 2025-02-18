package cyclicCounter

import (
	"sync"
	"time"
)

const (
	NO_ORDER    = 0
	UNCONFIRMED = 1
	CONFIRMED   = 2
)

type Orderstatus struct {
	Floor     int
	State     int
	Timestamp int64
}

type CyclicCounter struct {
	Orders map[int]Orderstatus
	Mu     sync.Mutex
}

func NewCyclicCounter() *CyclicCounter {
	return &CyclicCounter{Orders: make(map[int]Orderstatus)}
}

func (cc *CyclicCounter) RegisterOrder(floor int) {
	cc.Mu.Lock()
	defer cc.Mu.Unlock()

	if _, exists := cc.Orders[floor]; !exists {
		cc.Orders[floor] = Orderstatus{
			Floor:     floor,
			State:     UNCONFIRMED,
			Timestamp: time.Now().Unix(),
		}
	}
}

func (cc *CyclicCounter) ConfirmOrder(floor int) {
	cc.Mu.Lock()
	defer cc.Mu.Unlock()

	if order, exists := cc.Orders[floor]; exists && order.State == UNCONFIRMED {
		order.State = CONFIRMED
		cc.Orders[floor] = order
	}
}

func (cc *CyclicCounter) RemoveOrder(floor int) {
	cc.Mu.Lock()
	defer cc.Mu.Unlock()

	delete(cc.Orders, floor)
}
