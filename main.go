package main

//Importing broadcast and time
import (
	elevator "TTK4145---project/Elevator"
)

func main() {
	//Creating a new elevator
	elevator := elevator.Elevator{ID: 1, State: 0, Direction: 0, Floor: 0}
	//Starting the elevator
	go elevator.Run()
	//Starting the broadcast
	go elevator.StartBroadcast()
	go elevator.ListenForBroadcast()
	select {}
}
