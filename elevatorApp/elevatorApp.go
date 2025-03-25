package elevatorApp

import (
	network "TTK4145---project/Network-go"
	driver "TTK4145---project/driver-go"
	faultTolerance "TTK4145---project/faultTolerance-go"
	"context"
	"fmt"
	"sync"
)

type App struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	restartCh chan struct{}
}

// New creates a new elevator application instance with context and restart channel.
func New(ctx context.Context, restartCh chan struct{}) *App {
	childCtx, cancel := context.WithCancel(ctx)
	return &App{
		ctx:       childCtx,
		cancel:    cancel,
		restartCh: restartCh,
	}
}

// Start runs all the main elevator components concurrently and waits for stop signal.
func (a *App) Start() {
	fmt.Println("[App] Starting ElevatorApp...")

	a.runWithRecovery("RunElevator", func() {
		driver.RunElevatorWithContext(a.ctx)
	})
	a.runWithRecovery("MonitorNetwork", func() {
		faultTolerance.MonitorNetwork(a.ctx, a.restartCh)
	})
	a.runWithRecovery("MonitorMovement", func() {
		faultTolerance.MonitorMovement(a.ctx, a.restartCh)
	})
	a.runWithRecovery("Network", func() {
		network.Run(a.ctx)
	})
	// Block until context is canceled
	<-a.ctx.Done()
	fmt.Println("[App] Context canceled, waiting for components to stop...")
	a.wg.Wait()
	fmt.Println("[App] All components stopped.")
}

// Kill gracefully stops all components.
func (a *App) Kill() {
	fmt.Println("[App] Kill() called, cancelling context...")
	a.cancel()
}

// runWithRecovery runs a component function with panic recovery and context awareness.
func (a *App) runWithRecovery(name string, fn func()) {
	a.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[App] Recovered from panic in %s: %v\n", name, r)
				a.restartCh <- struct{}{}
			}
			a.wg.Done()
		}()
		// Run the component
		fn()
	}()
}
