package driver

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"fmt"
	"time"
)

var obstruction = false


func RunElevator() {

	numFloors := config.NumFloors

	elevio.Init("localhost:"+config.Port, numFloors)

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

	// If not at a floor, move down until a floor is reached
	config.ElevatorInstance.Floor = elevio.GetFloor()
	if config.ElevatorInstance.Floor == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
		for config.ElevatorInstance.Floor == -1 {
			config.ElevatorInstance.Floor = <-drv_floors
		}

		elevio.SetMotorDirection(elevio.MD_Stop)
	}

	// does nothing when elevator is online
	go offlineMode()
	time.Sleep(time.Second)
	

	decideDir()

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



		case <-config.MyQueue:
			decideDir()
		}
		
	}
}
