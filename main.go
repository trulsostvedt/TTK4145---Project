package main

import (
	network "TTK4145---project/Network-go"
	config "TTK4145---project/config"
	driver "TTK4145---project/driver-go"
	elevio "TTK4145---project/driver-go/elevio"
	faultTolerance "TTK4145---project/faultTolerance-go"
	"flag"
	"time"
)

func main() {

	// Start all necessary goroutines
	go faultTolerance.MonitorMovement()
	go faultTolerance.MonitorNetwork()
	time.Sleep(time.Second)
	go network.Network(&config.ElevatorInstance)
	go driver.RunElevator()

	select {}

}

func init() {
	flag.StringVar(&config.ElevatorInstance.ID, "id", "", "id of this peer")
	flag.StringVar(&config.Port, "port", "15657", "port to listen on")
	flag.Parse()

	queue := [config.NumFloors][config.NumButtons]config.OrderState{}
	for i := 0; i < config.NumFloors; i++ {
		for j := 0; j < config.NumButtons; j++ {
			queue[i][j] = config.Uninitialized
		}
	}

	config.ElevatorInstance = config.Elevator{
		ID:        config.ElevatorInstance.ID,
		State:     config.Idle,
		Direction: elevio.MD_Stop,
		Floor:     0,
		Queue:     queue,
	}
	driver.ReadCabOrders()

	config.Elevators = make(map[string]config.Elevator)
	config.Elevators[config.ElevatorInstance.ID] = config.ElevatorInstance
}
