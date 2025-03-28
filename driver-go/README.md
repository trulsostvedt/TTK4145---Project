Elevator driver for Go
======================

This module contains the elevator driver written in Go, which is responsible for controlling the operation of the elevator. The driver interfaces with the hardware through the `elevio` package, which is provided as part of the course materials.

### Features
- Handles button inputs and floor sensor signals.
- Controls motor direction and stops the elevator at the correct floors.
- Manages lights for buttons, floor indicators, and door operations.

### Dependencies
The project relies on the `elevio` package, which is pre-configured to communicate with the elevator hardware. Ensure that the `elevio` library is correctly set up before running the driver.

### Notes
- The `elevio` package is provided as part of the course and must be used to interact with the elevator hardware.
- Make sure the hardware is connected and powered on before running the driver.

### License
This project is for educational purposes as part of the TTK4145 course.









