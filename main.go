package main

import (
	network "TTK4145---project/Network-go"
	config "TTK4145---project/config"
	driver "TTK4145---project/driver-go"
	"flag"
)

func main() {

	go network.Network(&config.ElevatorInstance)
	go driver.RunElevator(&config.ElevatorInstance)

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
