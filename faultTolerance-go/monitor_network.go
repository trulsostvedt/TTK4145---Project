package faulttolerance

import (
	"TTK4145---project/config"
	"context"
	"fmt"
	"net"
	"time"
)

var (
	Interval         = 2 * time.Second
	LastNetworkCheck = time.Now()
	LastPeerMessage  = time.Now()
)

// MonitorNetwork monitors the network state and initiates a restart if:
// - the elevator is disconnected from the network,
// - and has no active confirmed orders (cab or hall).
func MonitorNetwork(ctx context.Context, restartCh chan struct{}) {
	fmt.Println("[MonitorNetwork] Started")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	offlineMode := false

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[MonitorNetwork] Shutting down gracefully")
			return

		case <-ticker.C:
			if time.Since(LastPeerMessage) < 5*time.Second {
				if offlineMode {
					fmt.Println("[MonitorNetwork] Network reconnected. Exiting offline mode.")
					config.IsOfflineMode = false
					offlineMode = false
				}
				continue
			}

			// If network is lost, check whether we should enter offline mode where we
			// still handle active orders, but do not accept new ones, or restart the elevator
			if !CheckNetworkStatus() {
				if hasActiveOrders() {
					if !offlineMode {
						fmt.Println("[MonitorNetwork] Network lost, but active orders remain. Entering local-only mode.")
						config.IsOfflineMode = true
						offlineMode = true
					}
					continue
				} else {
					fmt.Println("[MonitorNetwork] Network lost and no active orders. Requesting restart...")
					select {
					case restartCh <- struct{}{}:
					default:
						fmt.Println("[MonitorNetwork] Restart already requested.")
					}
					return
				}
			}

			// Offline mode completed orders -> restart
			if offlineMode && !hasActiveOrders() {
				fmt.Println("[MonitorNetwork] Completed all orders in offline mode. Requesting restart...")
				select {
				case restartCh <- struct{}{}:
				default:
					fmt.Println("[MonitorNetwork] Restart already requested.")
				}
				return
			}
		}
	}
}

// CheckNetworkStatus attempts to determine if the elevator has internet access.
// Used to distinguish between "I lost network" and "the entire network is gone".
func CheckNetworkStatus() bool {
	if time.Since(LastPeerMessage) < 5*time.Second {
		return true
	}

	if time.Since(LastNetworkCheck) >= 5*time.Second {
		conn, err := net.Dial("udp", "8.8.8.8:80") // Google's DNS server, always up
		if err != nil {
			return false
		}
		conn.Close()
		LastNetworkCheck = time.Now()
	}

	return true
}

// hasActiveOrders returns true if the elevator has any confirmed cab or hall calls
func hasActiveOrders() bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if config.ElevatorInstance.Queue[floor][config.ButtonCab] == config.Confirmed ||
			config.ElevatorInstance.Queue[floor][config.ButtonUp] == config.Confirmed ||
			config.ElevatorInstance.Queue[floor][config.ButtonDown] == config.Confirmed {
			return true
		}
	}
	return false
}
