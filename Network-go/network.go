package network

import (
	"TTK4145---project/Network-go/network/bcast"
	"TTK4145---project/Network-go/network/localip"
	"TTK4145---project/Network-go/network/peers"
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	faulttolerance "TTK4145---project/faultTolerance-go"

	// "TTK4145---project/driver-go"

	"fmt"
	"os"
	"time"
)

// Network is the main function for the network module
// It handles the communication between elevators using UDP broadcast
func Network(elevatorInstance *config.Elevator) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var id = elevatorInstance.ID

	// If we don't have an id, we generate one based on the local IP and PID
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// Channels for receiving and sending peer updates
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// Channels for receiving and sending elevator messages
	elevatorTx := make(chan config.Elevator) // Transmitter
	elevatorRx := make(chan config.Elevator) // Receiver
	go bcast.Transmitter(16569, elevatorTx)  // Transmitter
	go bcast.Receiver(16569, elevatorRx)     // Receiver

	go func() {
		for {
			elevatorTx <- *elevatorInstance
			time.Sleep(10 * time.Millisecond)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			for _, lostPeer := range p.Lost {
				delete(config.Elevators, lostPeer)
			}

		case a := <-elevatorRx:
			faulttolerance.LastPeerMessage = time.Now()

			elev := config.Elevator{
				ID:        a.ID,
				State:     a.State,
				Direction: a.Direction,
				Floor:     a.Floor,
				Queue:     a.Queue,
			}

			config.Elevators[a.ID] = elev

			SyncHallRequests() // Synchronize hall requests between elevators

			hra.HRA() // Run the HRA algorithm to update the elevator's state

		}
	}
}

// SyncHallRequests synchronizes the hall requests between all elevators
func SyncHallRequests() {

	for i := 0; i < config.NumFloors; i++ {
		// If this elevator has uninitialized requests, attempt to copy from other elevators
		if config.ElevatorInstance.Queue[i][config.ButtonUp] == config.Uninitialized {
			initialized := false
			for _, elev := range config.Elevators {
				// Check if any other elevator has a valid request
				// If so, copy it to this elevator
				if elev.Queue[i][config.ButtonUp] != config.Uninitialized {
					config.ElevatorInstance.Queue[i][config.ButtonUp] = elev.Queue[i][config.ButtonUp]
					initialized = true
					break
				}
			}
			if !initialized {
				// If no other elevator has a valid request, set this elevator's request to NoOrder
				config.ElevatorInstance.Queue[i][config.ButtonUp] = config.NoOrder
			}
		}
		// If this elevator has uninitialized requests, attempt to copy from other elevators
		if config.ElevatorInstance.Queue[i][config.ButtonDown] == config.Uninitialized {
			initialized := false
			for _, elev := range config.Elevators {
				// Check if any other elevator has a valid request
				if elev.Queue[i][config.ButtonDown] != config.Uninitialized {
					config.ElevatorInstance.Queue[i][config.ButtonDown] = elev.Queue[i][config.ButtonDown]
					initialized = true
					break
				}
			}
			if !initialized {
				config.ElevatorInstance.Queue[i][config.ButtonDown] = config.NoOrder
			}
		}
	}
	// Check if all elevators have confirmed the requests
	for i := 0; i < config.NumFloors; i++ {
		isConfirmedUp := true
		for _, elev := range config.Elevators {
			if elev.Queue[i][config.ButtonUp] != config.Unconfirmed {
				isConfirmedUp = false
				break
			}
		}
		if isConfirmedUp {
			config.ElevatorInstance.Queue[i][config.ButtonUp] = config.Confirmed
		}

		isConfirmedDown := true
		for _, elev := range config.Elevators {
			if elev.Queue[i][config.ButtonDown] != config.Unconfirmed {
				isConfirmedDown = false
				break
			}
		}
		if isConfirmedDown {
			config.ElevatorInstance.Queue[i][config.ButtonDown] = config.Confirmed
		}
	}

	// Check if any elevator has a request that this elevator does not have
	for i := 0; i < config.NumFloors; i++ {
		for _, elev := range config.Elevators {
			up := elev.Queue[i][config.ButtonUp] - config.ElevatorInstance.Queue[i][config.ButtonUp]
			down := elev.Queue[i][config.ButtonDown] - config.ElevatorInstance.Queue[i][config.ButtonDown]

			if up == 1 || up == -2 {
				config.ElevatorInstance.Queue[i][config.ButtonUp] = elev.Queue[i][config.ButtonUp]
			}

			if down == 1 || down == -2 {
				config.ElevatorInstance.Queue[i][config.ButtonDown] = elev.Queue[i][config.ButtonDown]
			}
		}
	}

}

// PeerUpdate is a struct that contains information about the peers in the network
type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}
