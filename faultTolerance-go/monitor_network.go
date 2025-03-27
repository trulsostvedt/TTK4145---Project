package faulttolerance

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go"
	"fmt"
	"net"
	"time"
)

var (
	Interval         = 2 * time.Second
	LastNetworkCheck = time.Now()
	LastPeerMessage  = time.Now()
)

// MonitorNetwork is a function that monitors the network connection.
// If the elevator has not received a message from another elevator in 10 seconds,
// it will attempt to restart itself.
// If the elevator is in offline mode and has completed all cab orders, it will restart itself.
func MonitorNetwork() {
	offlineMode := false

	for {
		time.Sleep(2 * time.Second)

		// Skip if a restart is already in progress
		if isRestarting {
			fmt.Println("[MonitorNetwork] Restart in progress. Skipping network check...")
			continue
		}

		// If we have received a message from another elevator in the last 10 seconds, continue
		if time.Since(LastPeerMessage) < 10*time.Second {
			if offlineMode {
				fmt.Println("[MonitorNetwork] Network reconnected. Exiting offline mode.")
				config.IsOfflineMode = false
			}
			offlineMode = false
			continue
		}

		// If we have lost network connection, check if we have active cab orders
		if !CheckNetworkStatus() {
			if hasActiveCabOrders() {
				//peerUpdateCh := make(chan peers.PeerUpdate)
				if !offlineMode {
					fmt.Println("[MonitorNetwork] Network lost, but cab orders remain. Entering local-only mode.")
					config.IsOfflineMode = true
					offlineMode = true
				}
				continue // Continue to check for network connection
			} else {
				fmt.Println("[MonitorNetwork] Network lost and no active cab orders. Restarting to reconnect...")
				RestartSelf()
			}
		}

		// If we are in offline mode and have completed all cab orders, restart to rejoin network
		if offlineMode && !hasActiveCabOrders() {
			fmt.Println("[MonitorNetwork] Cab orders completed in offline mode. Restarting to rejoin network.")
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

// RestartSelf restarts the elevator process by running the main.go file with the current elevator ID
// This is the same function for restarting the elevator process as in monitor_movement.go
// If the elevator is already restarting, it will not attempt to restart again.
// isRestarting is nesessary to avoid restarting the elevator both for network and movement issues
var isRestarting = false
var restartCooldown = 10 * time.Second // Cooldown period after a restart

// RestartSelf restarts the elevator process by running the main.go file with the current elevator ID
// filepath: /home/student/Documents/gr27/TTK4145---Project/faultTolerance-go/monitor_network.go
func RestartSelf() {
	panic("Restart program!")

}

// hasActiveCabOrders() checks if the elevator has any active cab orders.
// If the elevator has active cab orders, it should not restart.
func hasActiveCabOrders() bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if config.ElevatorInstance.Queue[floor][config.ButtonCab] == config.Confirmed {
			return true
		}
	}
	return false
}

func cleanup() {
	// Perform cleanup if necessary
	driver.StopElevator()
}
