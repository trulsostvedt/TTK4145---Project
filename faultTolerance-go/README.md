# Fault Tolerance Module

The `faultTolerance-go` module ensures that the elevator control system remains robust and operational even in the presence of network failures, crashes, or other unexpected issues. This module is designed to meet all fault tolerance requirements specified for the project.

---

## Features

### Offline Mode Initialization
- The elevator system is designed to start in **offline mode** by default. 
- **Reason**: This ensures that if the elevator cannot establish a network connection during startup, it will automatically operate in offline mode and handle local cab orders independently. This design choice provides a seamless fallback mechanism for uninterrupted operation.

### Network Loss Handling
- When the elevator loses its network connection, it performs a **hard restart**.
- **Reason**: A hard restart ensures that the elevator resets its state and attempts to reconnect to the network. If the network remains unavailable, the elevator transitions back to offline mode and continues to handle local cab orders. This guarantees that the elevator remains functional even in isolated conditions.

### Automatic Recovery
- The module continuously monitors the network connection. When the network is restored, the elevator synchronizes its state with other elevators in the system.
- **Synchronization**: Hall orders, cab orders, and elevator states are updated to ensure consistency across the system.

### Fault Tolerance Requirements
This module fulfills all fault tolerance requirements:
1. **Crash Recovery**: The elevator automatically restarts and synchronizes with the network after a crash.
2. **Network Failure**: The elevator operates in offline mode during network outages and handles local cab orders independently.
3. **Order Persistence**: All orders are retained and executed, even after restarts or network reconnections.
4. **Robust Synchronization**: The system ensures that all elevators have a consistent view of hall orders and states when the network is operational.

---

## Why Hard Restarts?
A hard restart is performed when the elevator loses its network connection to:
1. Reset the elevator's state and ensure it is in a known, stable condition.
2. Attempt to reconnect to the network and synchronize with other elevators.
3. Provide a fallback mechanism where the elevator transitions to offline mode if the network remains unavailable.

This approach minimizes the risk of inconsistent states and ensures that the elevator can continue operating independently when necessary.

---

## Implementation Details
- **Offline Mode**: The elevator handles local cab orders and will not receive new hall orders during network outages. However all hall orders that the elevator confirmed it would take will be handled in offline-mode before the elevator tries to restart. 
- **Network Monitoring**: The module uses periodic heartbeats to detect network availability.
- **State Synchronization**: Upon reconnection, the elevator synchronizes its state with the network to ensure consistency.

---

## Usage
To integrate the fault tolerance module into your project, ensure it is initialized during the elevator's startup process. The module will automatically handle network monitoring, offline mode transitions, and state synchronization.

---

## Conclusion
The `faultTolerance-go` module is a critical component of the elevator control system, ensuring robust operation under all conditions. By starting in offline mode, performing hard restarts during network failures, and meeting all fault tolerance requirements, the module guarantees a reliable and fault-tolerant elevator system.

---