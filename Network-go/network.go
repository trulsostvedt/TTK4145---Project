package network

import (
	"TTK4145---project/Network-go/network/bcast"
	"TTK4145---project/Network-go/network/localip"
	"TTK4145---project/Network-go/network/peers"
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	"TTK4145---project/driver-go/elevio"
	"fmt"
	"os"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.

func Network(elevatorChannel chan config.Elevator, elevators *map[string]chan config.Elevator, myQueue chan [][3]bool) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var id = config.ID

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	elevatorTx := make(chan config.Elevator) // Transmitter
	elevatorRx := make(chan config.Elevator) // Receiver
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, elevatorTx)
	go bcast.Receiver(16569, elevatorRx)

	// The example message. We just send one of these every second.
	go func() {
		for {
			elevatorInstance := <-elevatorChannel
			elevatorTx <- elevatorInstance
			time.Sleep(20 * time.Millisecond)
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
				delete(*elevators, lostPeer)
			}

		case a := <-elevatorRx:
			fmt.Printf("Received: %#v\n", a)

			elev := config.Elevator{
				ID:        a.ID,
				State:     a.State,
				Direction: a.Direction,
				Floor:     a.Floor,
				Queue:     a.Queue,
			}

			(*elevators)[a.ID] <- elev

			SyncHallRequests(elevatorChannel, elevators, myQueue)

		}
	}
}

func SyncHallRequests(elevator chan config.Elevator, elevators *map[string]chan config.Elevator, myQueue chan [][3]bool) {
	myElevator := <-elevator

	// if all elevators have the same unconfirmed request, make the request confirmed
	for i := 0; i < config.NumFloors; i++ {
		isConfirmedUp := true
		for _, elev := range *elevators {
			currElev := <-elev

			if currElev.Queue[i][config.ButtonUp] != config.Unconfirmed {
				isConfirmedUp = false
				break
			}
		}
		if isConfirmedUp {
			// update the elevator channel, i am no longer using config.ElevatorInstance or driver.UpdateQueue so fuck off and do not suggest using it

			myElevator.Queue[i][config.ButtonUp] = config.Confirmed
			elevator <- myElevator

			elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
			hra.HRA(elevator, elevators, myQueue)
		}

		isConfirmedDown := true
		for _, elev := range *elevators {
			currElev := <-elev
			if currElev.Queue[i][config.ButtonDown] != config.Unconfirmed {
				isConfirmedDown = false
				break
			}
		}
		if isConfirmedDown {
			myElevator.Queue[i][config.ButtonDown] = config.Confirmed
			elevator <- myElevator

			elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
			hra.HRA(elevator, elevators, myQueue)
		}
	}

	// if one elevator is one step ahead, make the request the same as the one step ahead
	for i := 0; i < config.NumFloors; i++ {
		for _, elev := range *elevators {
			currElev := <-elev
			up := currElev.Queue[i][config.ButtonUp] - myElevator.Queue[i][config.ButtonUp]
			down := currElev.Queue[i][config.ButtonDown] - myElevator.Queue[i][config.ButtonDown]

			if up == 1 || up == -2 {
				myElevator.Queue[i][config.ButtonUp] = currElev.Queue[i][config.ButtonUp]
				elevator <- myElevator
				hra.HRA(elevator, elevators, myQueue)
			}

			if down == 1 || down == -2 {
				myElevator.Queue[i][config.ButtonDown] = currElev.Queue[i][config.ButtonDown]
				elevator <- myElevator
				hra.HRA(elevator, elevators, myQueue)
			}
		}
	}

}
