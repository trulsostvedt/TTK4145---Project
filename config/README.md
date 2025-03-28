# Configuration for Elevator System (`config.go`)

This file defines the core configuration, constants, and data structures used in the elevator system. It also includes global variables that manage the state of the elevator and its interactions with other elevators in the network.

---

## **Contents**

1. [Constants](#constants)
2. [Types](#types)
3. [Structures](#structures)
4. [Global Variables](#global-variables)
5. [Usage](#usage)
6. [Dependencies](#dependencies)

---

## **Constants**

### System Configuration
- `NumFloors`: The number of floors in the building (default: `4`).
- `NumButtons`: The number of button types (Up, Down, Cab).

---

## **Types**

### **Button**
Represents the type of button in the elevator system:
- `ButtonUp`: Up button.
- `ButtonDown`: Down button.
- `ButtonCab`: Cab button.

### **OrderState**
Represents the state of an elevator order:
- `NoOrder`: No order exists.
- `Unconfirmed`: Order placed but not confirmed.
- `Confirmed`: Order confirmed.
- `Uninitialized`: Represents an uninitialized state.

### **ElevatorState**
Represents the current state of the elevator:
- `Idle`: Elevator is idle.
- `Moving`: Elevator is moving.
- `DoorOpen`: Elevator door is open.

---

## **Structures**

### **Elevator**
Represents the state and properties of an elevator:
- `ID`: Unique identifier for the elevator.
- `State`: Current state of the elevator (`Idle`, `Moving`, or `DoorOpen`).
- `Direction`: Current direction of the elevator (`elevio.MotorDirection`).
- `Floor`: Current floor of the elevator.
- `Queue`: A 2D array representing the queue of orders for the elevator. Each floor and button type has an associated `OrderState`.

---

## **Global Variables**

### **Shared Variables**
- `MyQueue`: A buffered channel to store elevator orders.
- `TimeSinceOrder`: A buffered channel to track the time since an order was placed.
- `ElevatorInstance`: Represents the current elevator instance.
- `Elevators`: A map of all elevators in the system, keyed by their unique IDs.
- `IsOfflineMode`: A boolean indicating whether the elevator is running in offline mode.
- `Port`: The port number used for running the elevator or simulator (default: `15657`).

---

## **Usage**

This file is used to:
1. Define the core configuration and constants for the elevator system.
2. Manage the state of the current elevator and its interactions with other elevators.
3. Provide shared global variables for communication and state tracking.

### Example
Here’s how you might use the `Elevator` struct and global variables in your code:
```go
config.ElevatorInstance = config.Elevator{
    ID:        "elevator1",
    State:     config.Idle,
    Direction: elevio.MD_Stop,
    Floor:     0,
    Queue:     [config.NumFloors][config.NumButtons]config.OrderState{},
}

config.IsOfflineMode = true
