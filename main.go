package main

import (
	network "TTK4145---project/Network-go"
	config "TTK4145---project/config"
	driver "TTK4145---project/driver-go"
	elevio "TTK4145---project/driver-go/elevio"
	"flag"
	"fmt"
	"time"
)

func main() {

	go network.Network(&config.ElevatorInstance)
	go driver.RunElevator()
	go monitorSelf() // Monitor the elevator's own state
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

func monitorSelf() {
	for {
		time.Sleep(2 * time.Second) // Sjekk status med lavere frekvens

		// Vi antar at vi har nettverk hvis vi har hørt fra en annen heis nylig.
		if time.Since(network.LastPeerMessage) < 10*time.Second {
			continue // Vi har nettverk, ingen restart nødvendig
		}

		// Hvis vi ikke har hørt noe på 10 sekunder, sjekk DNS for sikkerhets skyld.
		if !network.CheckNetworkStatus() {
			fmt.Println("Network failure detected. Restarting self...")
			time.Sleep(10 * time.Second)       // Gir tid til nettverket å komme tilbake
			if !network.CheckNetworkStatus() { // Hvis nettverket fortsatt er nede, restart
				fmt.Println("Restarting self due to detected network failure...")
				network.RestartSelf()
			}
		}
	}
}
