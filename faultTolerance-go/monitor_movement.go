package faulttolerance

import (
	"TTK4145---project/config"
	"context"
	"fmt"
	"time"
)

// timeoutBetweenFloors is the time the elevator waits before attempting to restart itself
const (
	timeoutBetweenFloors = 5 * time.Second
	tickRate             = 100 * time.Millisecond
)

// MonitorMovement detects if the elevator is stuck between floors for too long,
// and signals a restart through restartCh if necessary.
func MonitorMovement(ctx context.Context, restartCh chan struct{}) {
	fmt.Println("[Movementmonitor]: Started")
	lastFloor := config.ElevatorInstance.Floor
	lastFloorTime := time.Now()
	lastState := config.ElevatorInstance.State

	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[MovementMonitor] Shutting down gracefully")
			return

		case <-ticker.C:
			currentFloor := config.ElevatorInstance.Floor
			currentState := config.ElevatorInstance.State

			// Reset timer if elevator started moving
			if currentState == config.Moving && lastState != config.Moving {
				lastFloorTime = time.Now()
				fmt.Println("[MovementMonitor] Elevator started moving. Timer reset.")
			}

			// If reached a new floor, reset timer
			if currentFloor != -1 && currentFloor != lastFloor {
				lastFloor = currentFloor
				lastFloorTime = time.Now()
				fmt.Printf("[MovementMonitor] Reached floor %d\n", currentFloor)
			}
			// Timeout check
			if currentState == config.Moving && time.Since(lastFloorTime) > timeoutBetweenFloors {
				fmt.Println("[MovementMonitor] Timeout: Elevator stuck between floors.")
				fmt.Println("[MovementMonitor] Requesting restart...")
				select {
				case restartCh <- struct{}{}:
				default:
					fmt.Println("[MovementMonitor] Restart already requested.")
				}
				return
			}

			lastState = currentState
		}
	}
}
