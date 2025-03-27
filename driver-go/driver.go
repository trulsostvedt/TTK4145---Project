package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"context"
	"fmt"
)

var obstruction = false

func RunElevatorWithContext(ctx context.Context) {
	numFloors := config.NumFloors
	elevio.Init("localhost:"+config.Port, numFloors) // default 15657

	var d elevio.MotorDirection = elevio.MD_Stop
	elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(ctx, drv_buttons)
	go elevio.PollFloorSensor(ctx, drv_floors)
	go elevio.PollObstructionSwitch(ctx, drv_obstr)
	go elevio.PollStopButton(ctx, drv_stop)
	go setAllLightsLoop(ctx)

	initDone := make(chan struct{})
	go func() {
		fmt.Println("starting floor: ", config.ElevatorInstance.Floor)
		config.ElevatorInstance.Floor = elevio.GetFloor()
		if config.ElevatorInstance.Floor == -1 {
			elevio.SetMotorDirection(elevio.MD_Down)
			for config.ElevatorInstance.Floor == -1 {
				select {
				case <-ctx.Done():
					fmt.Println("[Driver] Shutdown during floor initialization")
					return
				case floor := <-drv_floors:
					config.ElevatorInstance.Floor = floor
				}
			}
			elevio.SetMotorDirection(elevio.MD_Stop)
		}
		fmt.Println("Init complete. Current floor:", config.ElevatorInstance.Floor)
		close(initDone)
	}()

	select {
	case <-ctx.Done():
		fmt.Println("[Driver] Shutdown before init completed.")
		return
	case <-initDone:
		// Continue normally
	}

	decideDir()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Driver] Shutting down cleanly...")
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(false)
			return

		case order := <-drv_buttons:
			fmt.Printf("%+v\n", order)
			if order.Button == elevio.BT_Cab {
				config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Confirmed
				saveCabOrders()
			} else {
				if config.IsOfflineMode {
					// If offline, we can't confirm new orders but we should still
					// handle the ones we have already confirmed while online (see faultTolerance-go/monitor_network.go)
					config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Confirmed
				} else {
					config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Unconfirmed
				}
			}

		case floor := <-drv_floors:
			config.ElevatorInstance.Floor = floor
			fmt.Printf("%+v\n", floor)
			decideDir()

		case obstr := <-drv_obstr:
			obstruction = obstr

		case <-drv_stop:
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					continue
				}
			}
		case <-config.MyQueue:
			decideDir()
		}
	}
}
