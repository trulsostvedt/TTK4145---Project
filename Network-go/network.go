package network

import (
	"TTK4145---project/Network-go/network/bcast"
	"TTK4145---project/Network-go/network/conn"
	"TTK4145---project/Network-go/network/localip"
	"TTK4145---project/Network-go/network/peers"
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	"TTK4145---project/driver-go/elevio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//
//	will be received as zero-values.

func Network(elevatorInstance *config.Elevator) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var id = elevatorInstance.ID

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
			elevatorTx <- *elevatorInstance
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
				delete(config.Elevators, lostPeer)
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

			config.Elevators[a.ID] = elev

			SyncHallRequests()

			hra.HRA()

		}
	}
}

func SyncHallRequests() {

	// if all elevators have the same unconfirmed request, make the request confirmed
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
			elevio.SetButtonLamp(elevio.BT_HallUp, i, true)
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
			elevio.SetButtonLamp(elevio.BT_HallDown, i, true)
		}
	}

	// if one elevator is one step ahead, make the request the same as the one step ahead
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

var (
	interval         = 2 * time.Second
	LastNetworkCheck = time.Now()
	LastPeerMessage  = time.Now()
	LastRestartTime  = time.Now() // Hindrer for hyppige restarter
)

// Mottar meldinger fra andre heiser og oppdaterer siste mottakstidspunkt
func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {
	var buf [1024]byte
	conn := conn.DialBroadcastUDP(port)
	if conn == nil {
		fmt.Println("Kunne ikke opprette UDP-forbindelse!")
		return
	}

	for {
		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, err := conn.ReadFrom(buf[0:])
		if err == nil && n > 0 {
			LastPeerMessage = time.Now() // Oppdater tid for siste mottatte melding
		}
	}
}

// Sjekker om nettverket fungerer
func CheckNetworkStatus() bool {
	// Hvis vi har fått en melding fra en annen heis de siste 10 sekundene, antar vi at vi har nettverk.
	if time.Since(LastPeerMessage) < 10*time.Second {
		return true
	}

	// Hvis det har gått mer enn 10 sekunder siden forrige DNS-sjekk, gjør en ny sjekk.
	if time.Since(LastNetworkCheck) >= 10*time.Second {
		conn, err := net.Dial("udp", "8.8.8.8:80") // Google DNS som nettverkstest
		if err != nil {
			return false // Nettverket er nede
		}
		conn.Close()
		LastNetworkCheck = time.Now() // Kun oppdater hvis sjekken var vellykket
	}

	return true
}

// Sjekker om heisen bør restarte seg selv
func SelfCheck() bool {
	if !CheckNetworkStatus() {
		fmt.Println("Nettverksfeil! Venter på gjenoppretting...")
		return false // Nettverket er nede
	}
	return true // Alt fungerer
}

// Restarter heisprosessen, men begrenser hvor ofte det kan skje
func RestartSelf() {
	if time.Since(LastRestartTime) < 30*time.Second { // Hindrer for mange restarter
		fmt.Println("Too many restarts. Waiting...")
		return
	}

	fmt.Println("Restarting elevator process...")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", "go run main.go -id="+config.ElevatorInstance.ID)
	} else {
		cmd = exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-id="+config.ElevatorInstance.ID)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to restart elevator:", err)
	} else {
		fmt.Println("Elevator restarted successfully.")
		LastRestartTime = time.Now()
		os.Exit(1) // Avslutter den gamle prosessen
	}
}
