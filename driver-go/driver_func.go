package driver

import (
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	"TTK4145---project/driver-go/elevio"
	"fmt"
	"os"
	"time"
)


func removeOrder(floor, button int) {
	config.ElevatorInstance.Queue[floor][button] = config.NoOrder

}
func removeOrders(floor int) {
	if floor == -1 {
		return
	}

	queue := <-config.MyQueue
	if queue[floor][int(config.ButtonUp)] {
		removeOrder(floor, int(config.ButtonUp))
	} else if queue[floor][int(config.ButtonDown)] {
		removeOrder(floor, int(config.ButtonDown))
	}
	if queue[floor][int(config.ButtonCab)] {
		removeOrder(floor, int(config.ButtonCab))
	}
}

// Decide and set direction
func decideDir() {
	// dont do something illegal
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
	if elevio.GetFloor() == -1 {
		return
	}

	// decide direction
	var direction elevio.MotorDirection
	if isOrderAbove() {
		direction = elevio.MD_Up
		config.ElevatorInstance.Direction = direction
	} else if isOrderBelow() {
		direction = elevio.MD_Down
		config.ElevatorInstance.Direction = direction
	} else {
		direction = elevio.MD_Stop
		config.ElevatorInstance.Direction = direction
	}

	// open the door if the elevator is in the correct floor
	if reachedFloor(elevio.BT_HallUp) && direction != elevio.MD_Down {
		elevio.SetMotorDirection(elevio.MD_Stop)
		go openDoor(elevio.GetFloor(), int(direction))
		return
	} else if reachedFloor(elevio.BT_HallDown) && direction != elevio.MD_Up {
		elevio.SetMotorDirection(elevio.MD_Stop)
		go openDoor(elevio.GetFloor(), int(direction))
		return
	} else if reachedFloor(elevio.BT_Cab) {
		elevio.SetMotorDirection(elevio.MD_Stop)
		go openDoor(elevio.GetFloor(), int(direction))
		return
	}

	// set the motor direction
	if direction != elevio.MD_Stop {
		config.ElevatorInstance.State = config.Moving
		elevio.SetMotorDirection(direction)

		return
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	config.ElevatorInstance.State = config.Idle

}

func reachedFloor(button elevio.ButtonType) bool {
	if config.ElevatorInstance.Floor == -1 {
		return false
	}
	queue := <-config.MyQueue

	return queue[config.ElevatorInstance.Floor][int(button)]

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

func openDoor(floor, button int) {

	elevio.SetMotorDirection(elevio.MD_Stop)
	fmt.Println("Door open in floor", floor)
	config.ElevatorInstance.State = config.DoorOpen

	time1 := time.Now()
	for {
		if time.Since(time1) > 3*time.Second {
			break
		}
	}
	fmt.Println("Door closing in floor", floor)
	if obstruction {
		go openDoor(floor, button)
		return
	}

	removeOrders(elevio.GetFloor())
	saveCabOrders()

	config.ElevatorInstance.State = config.Idle
	decideDir()
}

func saveCabOrders() {
	filename := "cabOrders" + config.ElevatorInstance.ID + ".txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for i := 0; i < config.NumFloors; i++ {
		file.WriteString(fmt.Sprintf("%d ", config.ElevatorInstance.Queue[i][2]))
	}
}

func ReadCabOrders() {
	// read cab orders from a file
	filename := "cabOrders" + config.ElevatorInstance.ID + ".txt"
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, no problem, just return
			return
		}
		fmt.Println(err)
		return
	}
	defer file.Close()
	// read the cab orders from the file
	var order int
	for i := 0; i < config.NumFloors; i++ {
		_, err := fmt.Fscanf(file, "%d", &order)
		fmt.Println("order", i, order)
		if err != nil {
			fmt.Println(err)
			return
		}

		config.ElevatorInstance.Queue[i][2] = config.OrderState(order)
	}
}

func setAllLights() {
	for i := 0; i < config.NumFloors; i++ {
		for j := 0; j < config.NumButtons; j++ {
			if config.ElevatorInstance.Queue[i][j] == config.Confirmed {
				elevio.SetButtonLamp(elevio.ButtonType(j), i, true)
			} else {
				elevio.SetButtonLamp(elevio.ButtonType(j), i, false)
			}
		}
	}

	elevio.SetFloorIndicator(config.ElevatorInstance.Floor)

	if config.ElevatorInstance.State == config.DoorOpen {
		elevio.SetDoorOpenLamp(true)
	} else {
		elevio.SetDoorOpenLamp(false)
	}
}

func offlineMode() {
	for {
		if config.IsOfflineMode {
			hra.HRA()
			time.Sleep(10 * time.Millisecond)
		}
	}

}

func StopElevator() {
	elevio.SetMotorDirection(elevio.MD_Stop)
}
