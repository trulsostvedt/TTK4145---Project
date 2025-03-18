package driver

import (
	"TTK4145---project/config"
	hra "TTK4145---project/cost_fns"
	"TTK4145---project/driver-go/elevio"
	"fmt"
)

func RunElevator(elevatorInstance chan config.Elevator, elevators *map[string]chan config.Elevator, myQueue chan [][3]bool) {

	numFloors := config.NumFloors

	elevio.Init("localhost:"+config.Port, numFloors) // default 15657

	var d elevio.MotorDirection = elevio.MD_Stop

	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			elevator := <-elevatorInstance
			fmt.Printf("%+v\n", a)
			if a.Button == elevio.BT_Cab {
				elevator.Queue[a.Floor][config.ButtonCab] = config.Confirmed
				elevatorInstance <- elevator
				elevio.SetButtonLamp(a.Button, a.Floor, true)
				hra.HRA(elevatorInstance, elevators, myQueue)
			} else {
				elevator.Queue[a.Floor][int(a.Button)] = config.Unconfirmed
				elevatorInstance <- elevator
			}

			decideDir(elevatorInstance, myQueue)

		case a := <-drv_floors:
			elevator := <-elevatorInstance
			elevator.Floor = a
			elevatorInstance <- elevator
			fmt.Printf("%+v\n", a)
			elevio.SetFloorIndicator(a)

			decideDir(elevatorInstance, myQueue)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			elevator := <-elevatorInstance
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevator.Direction = elevio.MD_Stop
				elevatorInstance <- elevator
			} else {
				elevio.SetMotorDirection(d)
				elevator.Direction = d
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		case <-myQueue:
			decideDir(elevatorInstance, myQueue)

		}

	}
}
