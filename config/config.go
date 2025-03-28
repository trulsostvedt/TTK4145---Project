package config

import (
	"TTK4145---project/driver-go/elevio"
)

// Constants for system configuration
const (
	NumFloors    = 4 // Number of floors in the building
	NumButtons   = 3 // Number of button types (Up, Down, Cab)
)


// Button represents the type of button in the elevator system
type Button int

const (
	ButtonUp Button = iota // Up button
	ButtonDown             // Down button
	ButtonCab              // Cab button
)

// OrderState represents the state of an elevator order using a cyclic counter
type OrderState int

const (
	NoOrder       OrderState = iota // No order
	Unconfirmed                     // Order placed but not confirmed
	Confirmed                       // Order confirmed
	Uninitialized = -1              // Uninitialized state
)

// ElevatorState represents the state of the elevator
type ElevatorState int

const (
	Idle     ElevatorState = iota // Elevator is idle
	Moving                        // Elevator is moving
	DoorOpen                      // Elevator door is open
)

// Elevator represents the state and properties of an elevator 
// that needs to be shared with other elevators
type Elevator struct {
	ID        string                        	// Unique identifier for the elevator
	State     ElevatorState                 	// Current state of the elevator
	Direction elevio.MotorDirection         	// Current direction of the elevator
	Floor     int                           	// Current floor of the elevator
	Queue     [NumFloors][NumButtons]OrderState // Queue of orders for the elevator
}

// Global variables
var (
	MyQueue         = make(chan [][3]bool, 10) 	// Buffered channel to store elevator orders
	TimeSinceOrder  = make(chan int, 10)       	// Buffered channel to track time since an order was placed
	ElevatorInstance Elevator                  	// Instance of the current elevator
	Elevators        map[string]Elevator       	// Map of all elevators in the system, keyed by their ID
	IsOfflineMode   = false                   	// True if the elevator is running in offline mode
	Port         	= "15657" 					// Port number used for running the elevator elevator or the simulator
)
