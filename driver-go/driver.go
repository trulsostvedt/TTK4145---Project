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

	// direction := decideDir()
	// setDir(direction)
	decideDir()

	for {
		setAllLights()
		// direction = decideDir()
		// setDir(direction) 
		select {
		case order := <-drv_buttons:
			fmt.Printf("%+v\n", order)
			if order.Button == elevio.BT_Cab {
				config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Confirmed
				saveCabOrders()
			} else {
				config.ElevatorInstance.Queue[order.Floor][order.Button] = config.Unconfirmed
			}
			decideDir()

		case floor := <-drv_floors:
			config.ElevatorInstance.Floor = floor
			fmt.Printf("%+v\n", floor)
			decideDir()

		case obstr := <-drv_obstr:
			// fmt.Printf("%+v\n", a)
			// if a {
			// 	elevio.SetMotorDirection(elevio.MD_Stop)
			// 	config.ElevatorInstance.Direction = elevio.MD_Stop
			// } else {
			// 	elevio.SetMotorDirection(d)
			// 	config.ElevatorInstance.Direction = d
			// }
			obstruction = obstr

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					// elevio.SetButtonLamp(b, f, false)
					continue
				}
			}
		case <-config.MyQueue:
			// direction = decideDir()
			// setDir(direction)
			decideDir()

		}

	}
}
