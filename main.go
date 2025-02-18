package main

import (
	broadcast "TTK4145---project/Broadcast"
	countConfirmed "TTK4145---project/CountConfirmed"
	cyclicCounter "TTK4145---project/CyclicCounter"
	"time"
)

func main() {
	cyclic := cyclicCounter.NewCyclicCounter()
	elevatorID := 1
	floor := 0
	moving := false

	go broadcast.StartBroadcast(elevatorID, cyclic, floor, moving)
	go broadcast.ListenForBroadcast(cyclic)
	go countConfirmed.StartChecker(cyclic)

	time.Sleep(20 * time.Second)
}
