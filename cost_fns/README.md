# Cost Function Module

This module provides a mechanism for distributing elevator requests in a distributed elevator system. It leverages the precompiled `hall_request_assigner` binary provided by the course, which implements a robust request distribution algorithm. The module is designed to integrate seamlessly into our system while allowing for customization and adaptation to our specific requirements.

## Rationale

The focus of TTK4145 is on distributed systems, fault tolerance, and related concepts, rather than on designing optimal elevator algorithms. To save time and effort, we use the provided `hall_request_assigner` binary for request distribution. This allows us to concentrate on the core learning objectives of the course while still ensuring a functional and efficient elevator system.

## Usage

The `hall_request_assigner` binary is used to assign hall requests to elevators. It takes the current state of the system as input and outputs the optimal assignment of requests. Below is an example of how to use the binary in our system:

### Input Format

The binary expects the following input data:

- **Hall Requests**: A list of hall requests.
- **Elevator States**: The current state of each elevator, including:
    - Current floor
    - Direction (up, down, or idle)
    - Behavior (moving, door open, or idle)
    - Cab requests 

### Integration

To integrate the binary into our system:

1. Collect the necessary data from the elevator controllers.
2. Format the data as required by the binary.
3. Execute the binary and capture its output.
4. Update the system's request assignments based on the binary's output.

### Example Workflow

Below is an example workflow for integrating the binary:

1. **Collect Data**: Gather the current state of all elevators and unassigned requests.
2. **Format Data**: Convert the collected data into the input format required by the binary.
3. **Run Binary**: Execute the `hall_request_assigner` binary with the formatted input.
4. **Parse Output**: Read the output and update the system's request assignments accordingly.

# Example input data
input_data = {
    "hallRequests": [
        [False, False],  # Floor 0: No hall requests
        [True, False],   # Floor 1: Down request
        [False, True],   # Floor 2: Up request
        [False, False]   # Floor 3: No hall requests
    ],
    "states": {
        "elevator_1": {
            "behaviour": "idle",
            "floor": 0,
            "direction": "stop",
            "cabRequests": [False, False, False, False]
        },
        "elevator_2": {
            "behaviour": "moving",
            "floor": 3,
            "direction": "down",
            "cabRequests": [False, False, False, False]
        }
    }
}

## Notes

- Ensure the `hall_request_assigner` binary is executable. If not, run:
    ```bash
    chmod a+x hall_request_assigner
    ```
- The binary must be placed in the same directory as the script or its path must be included in the system's `PATH` variable.

By using the `hall_request_assigner` binary, we can efficiently handle request distribution while focusing on the core aspects of distributed systems.

