package network

import (
	"TTK4145---project/Network-go/network/bcast"
	"TTK4145---project/Network-go/network/localip"
	"TTK4145---project/Network-go/network/peers"
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"

	"context"
	"fmt"
	"os"
	"time"
)

// Run is the entry point for the network system.
// It listens and broadcasts elevator state, and tracks peers.
// It stops cleanly when the context is cancelled.
func Run(ctx context.Context) {
	var id = config.ElevatorInstance.ID

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
	go peers.Transmitter(ctx, 15647, id, peerTxEnable)
	go peers.Receiver(ctx, 15647, peerUpdateCh)

	elevatorTx := make(chan config.Elevator)
	elevatorRx := make(chan config.Elevator)
	go bcast.Transmitter(ctx, 16569, elevatorTx)
	go bcast.Receiver(ctx, 16569, elevatorRx)

	// Goroutine to send elevator state at regular intervals
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[Network] Shutting down transmitter goroutine.")
				return
			case <-ticker.C:
				elevatorTx <- config.ElevatorInstance
			}
		}
	}()

	fmt.Println("[Network] Started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Network] Context cancelled, shutting down network listener.")
			return

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			for _, lostPeer := range p.Lost {
				delete(config.Elevators, lostPeer)
			}

		case a := <-elevatorRx:
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
		if config.ElevatorInstance.Queue[i][config.ButtonUp] == config.Uninitialized {
			for _, elev := range config.Elevators {
				if elev.Queue[i][config.ButtonUp] != config.Uninitialized {
					config.ElevatorInstance.Queue[i][config.ButtonUp] = elev.Queue[i][config.ButtonUp]
					break
				}
			}
		}
		if config.ElevatorInstance.Queue[i][config.ButtonDown] == config.Uninitialized {
			for _, elev := range config.Elevators {
				if elev.Queue[i][config.ButtonDown] != config.Uninitialized {
					config.ElevatorInstance.Queue[i][config.ButtonDown] = elev.Queue[i][config.ButtonDown]
					break
				}
			}
		}
	}

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
