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

// MonitorNetwork is a function that monitors the network connection.
// If the elevator has not received a message from another elevator in 10 seconds,
// it will attempt to restart itself.
// If the elevator is in offline mode and has completed all orders, it will restart itself.
// The elevator starts in offline mode and will exit offline mode if it receives a message from another elevator. 
func MonitorNetwork() {
	offlineMode := true // Start in offline mode
	config.IsOfflineMode = true // Start in offline mode

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

		// If we have lost network connection, check if we have active orders
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

		// If we are in offline mode and have completed all orders, restart to rejoin network
		if offlineMode && !hasActiveOrders() {
			fmt.Println("[MonitorNetwork] Orders completed in offline mode. Restarting to rejoin network.")
			RestartSelf() 
		}
	}
}

// CheckNetworkStatus checks if the network is up by pinging Google's DNS server.
// This function is nessesary so that an elevator can find out if he is the one without network
// or if the network is down on all elevators
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

// RestartSelf restarts the elevator process by running the start-bash file with the current elevator ID
func RestartSelf() {
	panic("Restart program!")

}

// hasActiveCabOrders() checks if the elevator has any active orders.
// If the elevator has active orders, it should not restart.
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
