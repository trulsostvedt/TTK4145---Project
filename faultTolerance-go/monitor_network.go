package faulttolerance

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"

	"TTK4145---project/config"
)

var (
	Interval         = 2 * time.Second
	LastNetworkCheck = time.Now()
	LastPeerMessage  = time.Now()
	LastRestartTime  = time.Now()
)

// MonitorNetwork is a function that monitors the network connection.
// If the elevator has not received a message from another elevator in 10 seconds,
// it will attempt to restart itself.
func MonitorNetwork() {
	for {
		time.Sleep(2 * time.Second)

		if time.Since(LastPeerMessage) < 10*time.Second {
			continue
		}

		if !CheckNetworkStatus() {
			fmt.Println("[MonitorNetwork] Network failure detected. Waiting for recovery...")
			time.Sleep(10 * time.Second)

			if !CheckNetworkStatus() {
				fmt.Println("[MonitorNetwork] Restarting self due to persistent network failure...")
				RestartSelf()
			}
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

// RestartSelf restarts the elevator process.
// It waits 30 seconds between each restart to prevent infinite restart loops.
func RestartSelf() {
	if time.Since(LastRestartTime) < 5*time.Second {
		fmt.Println("Too many restarts. Waiting...")
		return
	}

	fmt.Println("Restarting elevator process...")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" { // Windows
		cmd = exec.Command("cmd.exe", "/C", "start", "cmd.exe", "/K", "go run main.go -id="+config.ElevatorInstance.ID)
	} else { // Linux or Mac
		cmd = exec.Command("gnome-terminal", "--", "go", "run", "main.go", "-id="+config.ElevatorInstance.ID)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to restart elevator:", err)
	} else {
		fmt.Println("Elevator restarted successfully.")
		LastRestartTime = time.Now()
		os.Exit(1)
	}
}
