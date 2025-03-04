package config

import (
	hra "TTK4145---project/cost_fns/hra"
	"TTK4145---project/driver-go/elevio"
)

var Port = "15657"

var MyQueue = make(chan [][3]bool, 10)

var time_since_order = make(chan int, 10)

const (
	NumFloors    = 4
	NumButtons   = 3
	NumElevators = 3
)

type Button int

const (
	ButtonUp Button = iota
	ButtonDown
	ButtonCab
)

// type MotorDirection int

// const (
// 	MD_Down MotorDirection = iota - 1
// 	MD_Stop
// 	MD_Up
// )

type OrderState int

const (
	NoOrder OrderState = iota
	Unconfirmed
	Confirmed
)

type ElevatorState int

const (
	Idle ElevatorState = iota
	Moving
	DoorOpen
)

type Elevator struct {
	ID        string
	State     ElevatorState
	Direction elevio.MotorDirection
	Floor     int
	Queue     [NumFloors][NumButtons]OrderState
}

func (e *Elevator) UpdateQueue(floor int, button Button, state OrderState) {
	e.Queue[floor][button] = state
	hra.HRA()
}

var ElevatorInstance Elevator

var Elevators map[string]Elevator
