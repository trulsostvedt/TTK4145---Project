package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
)

func removeOrder(floor, button int, elevatorInstance chan config.Elevator) {
	elevator := <-elevatorInstance
	elevator.Queue[floor][button] = config.NoOrder
	elevatorInstance <- elevator

	elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
}

func decideDir(elevatorInstance chan config.Elevator, myQueue chan [][3]bool) elevio.MotorDirection {
	elevator := <-elevatorInstance

	if elevator.Direction == elevio.MD_Up && elevator.Floor == config.NumFloors-1 {
		return elevio.MD_Stop
	} else if elevator.Direction == elevio.MD_Down && elevator.Floor == 0 {
		return elevio.MD_Stop
	}

	if elevator.State == config.Idle {
		if isOrderAbove(elevatorInstance, myQueue) {
			return elevio.MD_Up
		} else if isOrderBelow(elevatorInstance, myQueue) {
			return elevio.MD_Down
		}
	} else if elevator.State == config.Moving {
		if reachedFloor(elevatorInstance, myQueue) {
			removeOrder(elevator.Floor, int(config.ButtonCab), elevatorInstance)
			removeOrder(elevator.Floor, int(elevator.Direction), elevatorInstance)
			return elevio.MD_Stop
		}
	} else if elevator.State == config.DoorOpen {
		if isOrderAbove(elevatorInstance, myQueue) {
			return elevio.MD_Up
		} else if isOrderBelow(elevatorInstance, myQueue) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	}

	return elevio.MD_Stop

}

func isOrderAbove(elevatorInstance chan config.Elevator, myQueue chan [][3]bool) bool {
	elevator := <-elevatorInstance
	queue := <-myQueue
	for i := elevator.Floor + 1; i < config.NumFloors; i++ {
		if queue[i][config.ButtonCab] {
			return true
		}
	}
	return false
}

func isOrderBelow(elevatorInstance chan config.Elevator, myQueue chan [][3]bool) bool {
	elevator := <-elevatorInstance
	queue := <-myQueue
	for i := elevator.Floor - 1; i >= 0; i-- {
		if queue[i][config.ButtonCab] {
			return true
		}
	}
	return false
}

func reachedFloor(elevatorInstance chan config.Elevator, myQueue chan [][3]bool) bool {
	elevator := <-elevatorInstance
	queue := <-myQueue
	for i := 0; i < config.NumButtons; i++ {
		if queue[elevator.Floor][i] {
			return true
		}
	}
	return false
}
