package main

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"TTK4145---project/elevatorApp"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// --- Parse flags ---
	flag.StringVar(&config.ElevatorInstance.ID, "id", "", "id of this peer")
	flag.StringVar(&config.Port, "port", "15657", "port to listen on")
	flag.Parse()

	// --- Initialize elevator state ---
	queue := [config.NumFloors][config.NumButtons]config.OrderState{}
	for i := 0; i < config.NumFloors; i++ {
		for j := 0; j < config.NumButtons; j++ {
			queue[i][j] = config.Uninitialized
		}
	}
	config.ElevatorInstance = config.Elevator{
		ID:        config.ElevatorInstance.ID,
		State:     config.Idle,
		Direction: elevio.MD_Stop,
		Floor:     0,
		Queue:     queue,
	}
	config.Elevators = make(map[string]config.Elevator)
	config.Elevators[config.ElevatorInstance.ID] = config.ElevatorInstance

	// --- Set up signal handler ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// --- Create restart channel ---
	restartCh := make(chan struct{})

	for {
		ctx, cancel := context.WithCancel(context.Background())
		app := elevatorApp.New(ctx, restartCh)

		// Handle interrupt (Ctrl+C) or SIGTERM
		go func() {
			select {
			case sig := <-sigChan:
				fmt.Printf("\n[Main] Caught signal: %s. Exiting cleanly.\n", sig)
				cancel()
				os.Exit(0)
			case <-restartCh:
				fmt.Println("[Main] Restart signal received. Restarting app...")
				cancel()
			}
		}()

		app.Start()
	}
}
