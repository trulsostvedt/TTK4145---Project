package faulttolerance

import (
	"TTK4145---project/config"
	"fmt"
	"time"
)

// timeoutBetweenFloors is the time the elevator waits before attempting to restart itself
const (
	timeoutBetweenFloors = 5 * time.Second
	tickRate             = 100 * time.Millisecond
)

// MonitorMovement is a function that monitors the movement of the elevator.
// If the elevator is stuck between floors, it will attempt to restart itself.
func MonitorMovement() {
	fmt.Println("[Movementmonitor]: Starting elevator movement monitor")
	lastFloor := config.ElevatorInstance.Floor
	lastFloorTime := time.Now()
	lastState := config.ElevatorInstance.State

	for {
		time.Sleep(tickRate)

		// Skip monitoring if a restart is already in progress
		if isRestarting {
			fmt.Println("[Movementmonitor]: Restart in progress. Skipping movement check...")
			continue
		}

		currentFloor := config.ElevatorInstance.Floor
		currentState := config.ElevatorInstance.State

		// Reset timer if the elevator starts moving
		if currentState == config.Moving && lastState != config.Moving {
			lastFloorTime = time.Now()
			fmt.Println("[Movementmonitor]: Elevator started moving. Timer reset.")
		}

		// Update last floor and reset timer when reaching a new floor
		if currentFloor != -1 && currentFloor != lastFloor {
			lastFloor = currentFloor
			lastFloorTime = time.Now()
			fmt.Printf("[Movementmonitor]: Reached floor %d\n", currentFloor)
		}

		// Check for timeout between floors
		if currentState == config.Moving &&
			time.Since(lastFloorTime) > timeoutBetweenFloors {

			fmt.Println("\n[Movementmonitor]: Timeout: Elevator appears stuck between floors.")
			fmt.Println("\n[Movementmonitor]: Attempting self-restart...")
			RestartSelf()
		}

		lastState = currentState
	}
}
