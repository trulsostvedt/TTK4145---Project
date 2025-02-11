package main

//Importing broadcast and time
import (
	"TTK4145---project/Broadcast"
)

func main() {
	go broadcast.Startbroadcast()
	go broadcast.ListenForBroadcast()
	select {}
}
