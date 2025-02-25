package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"time"
)

func removeOrder(floor, button int) {
	config.ElevatorInstance.Queue[floor][button] = config.NoOrder
}

func decideDir() elevio.MotorDirection {

	queue := <-config.MyQueue

	reachedFloor := false
	for i := 0; i < config.NumButtons; i++ {
		if queue[config.ElevatorInstance.Floor][i] {
			removeOrder(config.ElevatorInstance.Floor, i)
			reachedFloor = true
		}
		if reachedFloor {
			openDoor()
			return elevio.MD_Stop
		}
	}

	for i := 0; i < config.NumFloors; i++ {
		if queue[i][config.ButtonUp] || queue[i][config.ButtonDown] || queue[i][config.ButtonCab] {
			if i > config.ElevatorInstance.Floor {
				return elevio.MD_Up
			} else if i < config.ElevatorInstance.Floor {
				return elevio.MD_Down
			} else {
				// config.ElevatorInstance.Queue[i][config.ButtonUp] = config.NoOrder
				// config.ElevatorInstance.Queue[i][config.ButtonDown] = config.NoOrder
				return elevio.MD_Stop
			}
		}
	}
	return elevio.MD_Stop
}

func reachedFloor() bool {
	queue := <-config.MyQueue
	for i := 0; i < config.NumButtons; i++ {
		if queue[elevio.GetFloor()][i] {
			return true
		}
	}
	return false
}

func openDoor() {
	elevio.SetDoorOpenLamp(true)
	config.ElevatorInstance.State = config.DoorOpen
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
	config.ElevatorInstance.State = config.Idle
}
