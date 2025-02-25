package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"time"
)

func removeOrder(floor, button int) {
	config.ElevatorInstance.Queue[floor][button] = config.NoOrder
	elevio.SetButtonLamp(elevio.ButtonType(button), floor, false)
}

func decideDir() {

	if config.ElevatorInstance.Floor == 0 && config.ElevatorInstance.Direction == elevio.MD_Down {
		config.ElevatorInstance.Direction = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
	if config.ElevatorInstance.Floor == config.NumFloors-1 && config.ElevatorInstance.Direction == elevio.MD_Up {
		config.ElevatorInstance.Direction = elevio.MD_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
	}

	queue := <-config.MyQueue

	reachedFloor := false
	for i := 0; i < config.NumButtons; i++ {
		if queue[config.ElevatorInstance.Floor][i] {
			removeOrder(config.ElevatorInstance.Floor, i)
			reachedFloor = true
		}
		if reachedFloor {
			elevio.SetMotorDirection(elevio.MD_Stop)
			openDoor()

		}
	}

	for i := 0; i < config.NumFloors; i++ {
		if queue[i][config.ButtonUp] || queue[i][config.ButtonDown] || queue[i][config.ButtonCab] {
			if i > config.ElevatorInstance.Floor {
				config.ElevatorInstance.State = config.Moving
				config.ElevatorInstance.Direction = elevio.MD_Up
				elevio.SetMotorDirection(elevio.MD_Up)
				break
			} else if i < config.ElevatorInstance.Floor {
				config.ElevatorInstance.State = config.Moving
				config.ElevatorInstance.Direction = elevio.MD_Down
				elevio.SetMotorDirection(elevio.MD_Down)
				break
			} else {
				// config.ElevatorInstance.Queue[i][config.ButtonUp] = config.NoOrder
				// config.ElevatorInstance.Queue[i][config.ButtonDown] = config.NoOrder
				config.ElevatorInstance.State = config.Idle
				config.ElevatorInstance.Direction = elevio.MD_Stop
				elevio.SetMotorDirection(elevio.MD_Stop)

			}
		}
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

func openDoor() {
	elevio.SetDoorOpenLamp(true)
	config.ElevatorInstance.State = config.DoorOpen
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
	config.ElevatorInstance.State = config.Idle
}
