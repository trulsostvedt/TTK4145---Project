package main

import (
	network "TTK4145---project/Network-go"
	"TTK4145---project/config"
	"flag"
	"fmt"
)

func main() {

	flag.Parse()
	fmt.Println("Hello World")
	network.Network()

	select {}
}

func init() {
	flag.StringVar(&config.ElevatorInstance.ID, "id", "", "id of this peer")
	flag.Parse()

	config.ElevatorInstance = config.Elevator{
		ID:        config.ElevatorInstance.ID,
		State:     config.Idle,
		Direction: config.MD_Stop,
		Floor:     0,
		Queue:     [config.NumFloors][config.NumButtons]config.OrderState{},
	}
}
