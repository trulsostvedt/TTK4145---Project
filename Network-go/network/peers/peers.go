package peers

import (
	"TTK4145---project/Network-go/network/conn"
	"context"
	"fmt"
	"net"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(ctx context.Context, port int, id string, transmitEnable <-chan bool) {
	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[peers] Transmitter shutting down")
			return
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(ctx context.Context, port int, peerUpdateCh chan<- PeerUpdate) {
	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[peers] Receiver shutting down")
			return
		default:
			updated := false
			conn.SetReadDeadline(time.Now().Add(interval))
			n, _, _ := conn.ReadFrom(buf[0:])

			id := string(buf[:n])
			p.New = ""
			if id != "" {
				if _, exists := lastSeen[id]; !exists {
					p.New = id
					updated = true
				}
				lastSeen[id] = time.Now()
			}

			p.Lost = make([]string, 0)
			for k, v := range lastSeen {
				if time.Since(v) > timeout {
					updated = true
					p.Lost = append(p.Lost, k)
					delete(lastSeen, k)
				}
			}

			if updated {
				p.Peers = make([]string, 0, len(lastSeen))
				for k := range lastSeen {
					p.Peers = append(p.Peers, k)
				}
				sort.Strings(p.Peers)
				sort.Strings(p.Lost)
				peerUpdateCh <- p
			}
		}
	}
}
