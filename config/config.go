package config

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

type MotorDirection int

const (
	MD_Down MotorDirection = iota - 1
	MD_Stop
	MD_Up
)

type orderState int

const (
	NoOrder orderState = iota
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
	State     ElevatorState
	Direction MotorDirection
	Floor     int
	Queue     [NumFloors][NumButtons]orderState
}
