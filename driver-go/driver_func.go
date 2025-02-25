package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"fmt"
	"time"
)

func removeOrder(floor, button int) {
	config.ElevatorInstance.Queue[floor][button] = config.NoOrder
}

func decideDir() {

	queue := <-config.MyQueue

	reachedFloor := false
	for i := 0; i < config.NumButtons; i++ {
		if queue[config.ElevatorInstance.Floor][i] {
			removeOrder(config.ElevatorInstance.Floor, i)
			reachedFloor = true
		}
		if reachedFloor {
			elevio.SetMotorDirection(elevio.MD_Stop)
			fmt.Println("Stopping")
			// openDoor()

		}
	}

	for i := 0; i < config.NumFloors; i++ {
		if queue[i][config.ButtonUp] || queue[i][config.ButtonDown] || queue[i][config.ButtonCab] {
			if i > config.ElevatorInstance.Floor {
				config.ElevatorInstance.State = config.Moving
				config.ElevatorInstance.Direction = elevio.MD_Up
				fmt.Println("Moving up")
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if i < config.ElevatorInstance.Floor {
				config.ElevatorInstance.State = config.Moving
				config.ElevatorInstance.Direction = elevio.MD_Down
				fmt.Println("Moving down")
				elevio.SetMotorDirection(elevio.MD_Down)
			} else {
				// config.ElevatorInstance.Queue[i][config.ButtonUp] = config.NoOrder
				// config.ElevatorInstance.Queue[i][config.ButtonDown] = config.NoOrder
				config.ElevatorInstance.State = config.Idle
				config.ElevatorInstance.Direction = elevio.MD_Stop
				fmt.Println("Stopping")
				elevio.SetMotorDirection(elevio.MD_Stop)

			}
		}
	}
	config.ElevatorInstance.State = config.Idle
	config.ElevatorInstance.Direction = elevio.MD_Stop
	elevio.SetMotorDirection(elevio.MD_Stop)
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
