package main

//Importing broadcast and time
import (
	"TTK4145---project/Broadcast"
)



type Elevator struct {
	ID int
	State int
	Direction int
	Floor int
	OrderMatrix [4][3]int
}




func main() {
	go broadcast.Startbroadcast()
	go broadcast.ListenForBroadcast()
	select {}
}
