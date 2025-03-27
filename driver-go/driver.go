package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"fmt"
)

var obstruction = false

func RunElevator() {

	numFloors := config.NumFloors

	elevio.Init("localhost:"+config.Port, numFloors) // default 15657

	var d elevio.MotorDirection = elevio.MD_Stop

	elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	fmt.Println("starting floor: ", config.ElevatorInstance.Floor)
	config.ElevatorInstance.Floor = elevio.GetFloor()
	fmt.Println("starting floor: ", config.ElevatorInstance.Floor)
	if config.ElevatorInstance.Floor == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
		for config.ElevatorInstance.Floor == -1 {
			config.ElevatorInstance.Floor = <-drv_floors
		}

		elevio.SetMotorDirection(elevio.MD_Stop)
	}

	// direction := decideDir()
	// setDir(direction)
	decideDir()

	//go listenForQueueChanges()

	for {
		setAllLights()

		if config.ElevatorInstance.State == config.Idle && elevio.GetFloor() == -1 {
			elevio.SetMotorDirection(elevio.MD_Down)
			config.ElevatorInstance.Floor = -1
			for config.ElevatorInstance.Floor == -1 {
				config.ElevatorInstance.Floor = <-drv_floors
			}
			elevio.SetMotorDirection(elevio.MD_Stop)
		}

		select {
		case order := <-drv_buttons:
			fmt.Printf("%+v\n", order)
			if order.Button == elevio.BT_Cab {
				config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Confirmed
				saveCabOrders()
			} else {
				if !config.IsOfflineMode {
					config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Unconfirmed
				}
			}

		case floor := <-drv_floors:
			config.ElevatorInstance.Floor = floor
			fmt.Printf("%+v\n", floor)
			decideDir()

		case obstr := <-drv_obstr:
			obstruction = obstr

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					continue
				}
			}

		case <-config.MyQueue:
			decideDir()
		}
		//hra.HRA()
		/*
			// Check if the queue has changed
			if hasQueueChanged(config.ElevatorInstance.Queue, previousQueue) {
				previousQueue = config.ElevatorInstance.Queue // Update the previous state
				hra.HRA()                                     // Run HRA only when the queue changes
			}
		*/
	}
}
