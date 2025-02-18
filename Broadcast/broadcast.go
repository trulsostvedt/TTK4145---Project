package broadcast

import (
	"fmt"
	"net"
	"time"
)

const (
	broadcast_IP      = "localhost:20008"
	broadcastInterval = 1 * time.Second
	message           = "Jeg broadcaster!"
)

func Startbroadcast() {
	conn, err := net.Dial("udp", broadcast_IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("Broadcasted: ", message)
		}
		time.Sleep(broadcastInterval)
	}
}

func ListenForBroadcast() {
	addr, err := net.ResolveUDPAddr("udp", broadcast_IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)

	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received: ", string(buf[0:n]))
	}
}
