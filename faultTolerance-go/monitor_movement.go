package faulttolerance

import (
	"TTK4145---project/config"
	"fmt"
	"time"
)

const (
	timeoutBetweenFloors = 10 * time.Second
	tickRate             = 100 * time.Millisecond
)

// MonitorMovement is a function that monitors the movement of the elevator.
// If the elevator is stuck between floors, it will attempt to restart itself.
func MonitorMovement() {
	fmt.Println("[Movementmonitor]: Starting elevator movement monitor")
	lastFloor := config.ElevatorInstance.Floor
	lastFloorTime := time.Now()

	for {
		time.Sleep(tickRate)

		currentFloor := config.ElevatorInstance.Floor

		if currentFloor != -1 && currentFloor != lastFloor {
			lastFloor = currentFloor
			lastFloorTime = time.Now()
			fmt.Printf("[Movementmonitor]: Reached floor %d\n", currentFloor)
		}

		if config.ElevatorInstance.State == config.Moving &&
			time.Since(lastFloorTime) > timeoutBetweenFloors {

			fmt.Println("\n [Movementmonitor]: Timeout: Elevator appears stuck between floors.")
			fmt.Println("\n [Movementmonitor]: Attempting self-restart...")
			RestartSelf()
		}
	}
}
