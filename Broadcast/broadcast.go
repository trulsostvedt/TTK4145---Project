package broadcast

import (
	cyclicCounter "TTK4145---project/CyclicCounter"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	broadcast_IP      = "10.22.91.255:20008" // alternativt "0.0.0.0:20008"
	broadcastInterval = 1 * time.Second
)

type ElevatorState struct {
	ID     int         `json:"id"`
	Floor  int         `json:"floor"`
	Moving bool        `json:"moving"`
	Orders map[int]int `json:"orders"` // map[floor] = state
}

var (
	worldState = make(map[int]ElevatorState) // Verdensbilde
	stateMutex = sync.Mutex{}                // Lås for synkronisering
)

func StartBroadcast(id int, cyclic *cyclicCounter.CyclicCounter, floor int, moving bool) {
	conn, err := net.Dial("udp", broadcast_IP)
	if err != nil {
		fmt.Println("Broadcast error:", err)
		return
	}
	defer conn.Close()

	for {
		stateMutex.Lock()
		state := ElevatorState{
			ID:     id,
			Floor:  floor,
			Moving: moving,
			Orders: copyOrders(cyclic),
		}
		stateMutex.Unlock()

		data, _ := json.Marshal(state)
		_, err := conn.Write(data)
		if err != nil {
			fmt.Println("Error broadcasting:", err)
		} else {
			fmt.Println("Broadcasted:", string(data))
		}

		time.Sleep(broadcastInterval)
	}
}

func ListenForBroadcast(cyclic *cyclicCounter.CyclicCounter) {
	addr, err := net.ResolveUDPAddr("udp", broadcast_IP)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving broadcast:", err)
			continue
		}

		var incoming ElevatorState
		err = json.Unmarshal(buf[:n], &incoming)
		if err != nil {
			fmt.Println("Error decoding broadcast:", err)
			continue

		}

		stateMutex.Lock()
		worldState[incoming.ID] = incoming
		stateMutex.Unlock()
		fmt.Println("Received:", incoming)
		updateCyclicCounter(cyclic, incoming)
	}
}

func copyOrders(cyclic *cyclicCounter.CyclicCounter) map[int]int {
	cyclic.Mu.Lock()
	defer cyclic.Mu.Unlock()

	ordersCopy := make(map[int]int)
	for floor, order := range cyclic.Orders {
		ordersCopy[floor] = order.State
	}
	return ordersCopy
}

func updateCyclicCounter(cyclic *cyclicCounter.CyclicCounter, newState ElevatorState) {
	cyclic.Mu.Lock()
	defer cyclic.Mu.Unlock()

	for floor, state := range newState.Orders {
		if _, exists := cyclic.Orders[floor]; !exists {
			cyclic.Orders[floor] = cyclicCounter.Orderstatus{
				Floor: floor,
				State: state,
			}
		}
	}
}
