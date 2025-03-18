package main

import (
	network "TTK4145---project/Network-go"
	config "TTK4145---project/config"
	driver "TTK4145---project/driver-go"
	elevio "TTK4145---project/driver-go/elevio"
	"flag"
	"fmt"
)

func main() {

	myQueue := make(chan [][3]bool, 10)

	ElevatorInstance := config.Elevator{
		ID:        config.ID,
		State:     config.Idle,
		Direction: elevio.MD_Stop,
		Floor:     0,
		Queue:     [config.NumFloors][config.NumButtons]config.OrderState{},
	}
	elevatorChannel := make(chan config.Elevator, 1)

	Elevators := make(map[string]chan config.Elevator)
	Elevators[config.ID] = elevatorChannel

	go network.Network(elevatorChannel, &Elevators, myQueue)
	go driver.RunElevator(elevatorChannel, &Elevators, myQueue)

	select {
	case elevatorChannel <- ElevatorInstance:
		fmt.Println("ElevatorInstance sent to channel")
	default:
		fmt.Println("Channel is full, message not sent")
	}

	select {}
}

func init() {
	flag.StringVar(&config.ID, "id", "", "id of this peer")
	flag.StringVar(&config.Port, "port", "15657", "port to listen on")
	flag.Parse()
}
