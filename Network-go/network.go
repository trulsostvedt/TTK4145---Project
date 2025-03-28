package network

import (
	"TTK4145---project/Network-go/network/bcast"
	"TTK4145---project/Network-go/network/localip"
	"TTK4145---project/Network-go/network/peers"
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	faulttolerance "TTK4145---project/faultTolerance-go"
	"fmt"
	"os"
	"time"
)



func Network(elevatorInstance *config.Elevator) {

	var id = elevatorInstance.ID

	// If we have not been assigned an ID, we generate one based on our IP and PID
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}


	peerUpdateCh := make(chan peers.PeerUpdate)

	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	elevatorTx := make(chan config.Elevator) // Transmitter
	elevatorRx := make(chan config.Elevator) // Receiver

	go bcast.Transmitter(16569, elevatorTx)
	go bcast.Receiver(16569, elevatorRx)


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
			//fmt.Printf("Received: %#v\n", a)

			elev := config.Elevator{
				ID:        a.ID,
				State:     a.State,
				Direction: a.Direction,
				Floor:     a.Floor,
				Queue:     a.Queue,
			}

			config.Elevators[a.ID] = elev

			SyncHallRequests()

			hra.HRA()

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
				if elev.Queue[i][config.ButtonUp] != config.Uninitialized {
					config.ElevatorInstance.Queue[i][config.ButtonUp] = elev.Queue[i][config.ButtonUp]
					initialized = true
					break
				}
			}
			if !initialized {
				config.ElevatorInstance.Queue[i][config.ButtonUp] = config.NoOrder
			}
		}

		if config.ElevatorInstance.Queue[i][config.ButtonDown] == config.Uninitialized {
			initialized := false
			for _, elev := range config.Elevators {
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

	// if all elevators have uncontested requests, confirm them
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

	// If one elevator is one step ahead in syclic counter they are correct	
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

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}
