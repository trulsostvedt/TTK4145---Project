package countConfirmed

import (
	cyclicCounter "TTK4145---project/CyclicCounter"
	"sync"
	"time"
)

var (
	checkInterval = 500 * time.Millisecond
	worldState    map[int]map[int]int
	stateMutex    = sync.Mutex{}
)

func StartChecker(cyclic *cyclicCounter.CyclicCounter) {
	for {
		stateMutex.Lock()
		checkAndConfirmOrders(cyclic)
		stateMutex.Unlock()

		time.Sleep(checkInterval)
	}
}

func checkAndConfirmOrders(cyclic *cyclicCounter.CyclicCounter) {
	cyclic.Mu.Lock()
	defer cyclic.Mu.Unlock()

	orderCounts := make(map[int]int)

	for _, elevator := range worldState {
		for floor, state := range elevator {
			if state == cyclicCounter.UNCONFIRMED {
				orderCounts[floor]++
			}
		}
	}

	for floor, count := range orderCounts {
		if count == len(worldState) {
			if order, exists := cyclic.Orders[floor]; exists && order.State == cyclicCounter.UNCONFIRMED {
				order.State = cyclicCounter.CONFIRMED
				cyclic.Orders[floor] = order
			}
		}
	}
}

func UpdateWorldState(elevatorID int, orders map[int]int) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	worldState[elevatorID] = orders
}
