package faulttolerance

import (
	"TTK4145---project/config"
	"fmt"
	"net"
	"time"
)

var (
	Interval         = 2 * time.Second
	LastNetworkCheck = time.Now()
	LastPeerMessage  = time.Now().Add(-10 * time.Second)
)


func MonitorNetwork() {
	offlineMode := true // Start in offline mode
	config.IsOfflineMode = true

	for {
		time.Sleep(5 * time.Second)

		// Check if we have received a message from another elevator in the last 10 seconds
		if time.Since(LastPeerMessage) < 10*time.Second {
			if offlineMode {
				fmt.Println("[MonitorNetwork] Network reconnected. Exiting offline mode.")
				config.IsOfflineMode = false
				offlineMode = false
			}
			continue
		}

		// If we have lost network connection, check if we have active cab orders
		if !CheckNetworkStatus() {
			if hasActiveOrders() {
				if !offlineMode {
					fmt.Println("[MonitorNetwork] Network lost, but active orders remain. Entering offline mode.")
					config.IsOfflineMode = true
					offlineMode = true
				}
				continue // Continue to check for network connection
			} else {
				fmt.Println("[MonitorNetwork] Network lost and no active orders. Restarting to reconnect...")
				RestartSelf()
			}
		}

		// If we are in offline mode and have completed all cab orders, restart to rejoin network
		if offlineMode && !hasActiveOrders() {
			fmt.Println("[MonitorNetwork] Orders completed in offline mode. Restarting to rejoin network.")
			RestartSelf()
		}
	}
}

// CheckNetworkStatus checks if the network is up by pinging Google's DNS server.
func CheckNetworkStatus() bool {
	if time.Since(LastPeerMessage) < 10*time.Second {
		return true
	}

	if time.Since(LastNetworkCheck) >= 10*time.Second {
		conn, err := net.Dial("udp", "8.8.8.8:80") // Google's DNS server, always up
		if err != nil {
			return false
		}
		conn.Close()
		LastNetworkCheck = time.Now()
	}

	return true
}


func RestartSelf() {
	panic("Restart program!")

}

func hasActiveOrders() bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if config.ElevatorInstance.Queue[floor][config.ButtonCab] == config.Confirmed ||
			config.ElevatorInstance.Queue[floor][config.ButtonUp] == config.Confirmed ||
			config.ElevatorInstance.Queue[floor][config.ButtonDown] == config.Confirmed {
			fmt.Println("[MonitorNetwork] Active orders found.")
			return true
		}
	}
	fmt.Println("[MonitorNetwork] No active orders found.")
	return false
}
