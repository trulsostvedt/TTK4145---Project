package elevator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

const (
	IDLE = iota
	MOVING
	DOOR_OPEN
)

const (
	UP = iota
	DOWN
	STOP
)

type Elevator struct {
	ID        int
	State     int
	Direction int
	Floor     int
}

func (e *Elevator) Run() {
	//Run the elevator
	// change state and direction and floor from time to time
	var state = IDLE
	var direction = STOP
	var floor = 0
	for {
		state += 1
		if state > 2 {
			state = 0
		}
		e.State = state
		direction += 1
		if direction > 2 {
			direction = 0
		}
		e.Direction = direction
		if direction == UP {
			floor += 1
		} else if direction == DOWN {
			floor -= 1
		}

		time.Sleep(2 * time.Second)

	}

}

func (e *Elevator) StartBroadcast() {
	//Start the broadcast
	conn, err := net.Dial("udp", "localhost:20008")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(e)
		if err != nil {
			return
		}

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			return
		}
		fmt.Println("Broadcasted: ", e)
		time.Sleep(1 * time.Second)
	}

}

func (e *Elevator) ListenForBroadcast() {
	//Listen for broadcast
	conn, err := net.ListenPacket("udp", "localhost:20008")
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var buf [1024]byte
		n, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			return
		}

		dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
		err = dec.Decode(e)
		if err != nil {
			return
		}
		fmt.Println("Received: ", e)
	}
}
