package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"time"
)

//TODO: Decide direction only decides what direction it should go next, but do not set the motordirection.

func removeOrder(floor, button int) {
	config.ElevatorInstance.Queue[floor][button] = config.NoOrder
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
}

func decideDir() {

	if config.ElevatorInstance.Floor == 0 && config.ElevatorInstance.Direction == elevio.MD_Down {
		config.ElevatorInstance.Direction = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		return
	}
	if config.ElevatorInstance.Floor == config.NumFloors-1 && config.ElevatorInstance.Direction == elevio.MD_Up {
		config.ElevatorInstance.Direction = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		return
	}

	if config.ElevatorInstance.State == config.DoorOpen {
		elevio.SetMotorDirection(elevio.MD_Stop)
		return
	}

	queue := <-config.MyQueue

	// Check if there is an order at the current floor and stop
	for i := 0; i < config.NumButtons; i++ {
		if queue[config.ElevatorInstance.Floor][i] {

			elevio.SetMotorDirection(elevio.MD_Stop)
			go openDoor(config.ElevatorInstance.Floor, i)
			break
		}

	}

	// if reachedFloor() {

	// }

	if isOrderAbove() {
		config.ElevatorInstance.State = config.Moving
		config.ElevatorInstance.Direction = elevio.MD_Up
		elevio.SetMotorDirection(elevio.MD_Up)
		return
	}
	if isOrderBelow() {
		config.ElevatorInstance.State = config.Moving
		config.ElevatorInstance.Direction = elevio.MD_Down
		elevio.SetMotorDirection(elevio.MD_Down)
		return
	}

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

func isOrderAbove() bool {
	queue := <-config.MyQueue
	for i := config.ElevatorInstance.Floor + 1; i < config.NumFloors; i++ {
		for j := 0; j < config.NumButtons; j++ {
			if queue[i][j] {
				return true
			}
		}
	}
	return false
}

func isOrderBelow() bool {
	queue := <-config.MyQueue
	for i := 0; i < config.ElevatorInstance.Floor; i++ {
		for j := 0; j < config.NumButtons; j++ {
			if queue[i][j] {
				return true
			}
		}
	}
	return false
}

func openDoor(floor, i int) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	config.ElevatorInstance.State = config.DoorOpen
	removeOrder(floor, int(config.ButtonCab))
	removeOrder(floor, int(config.ElevatorInstance.Direction))
	time1 := time.Now()
	for {
		if time.Since(time1) > 3*time.Second {
			break
		}
	}
	elevio.SetDoorOpenLamp(false)
	config.ElevatorInstance.State = config.Idle
	decideDir()
}
