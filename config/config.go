package config

import (
	"TTK4145---project/Driver-go/elevio"
)

const (
	NumFloors    = 4
	NumButtons   = 3
	NumElevators = 3
)

type Button int

const (
	BunnonUp Button = iota
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

var ElevatorInstance Elevator

var Elevators map[string]Elevator
