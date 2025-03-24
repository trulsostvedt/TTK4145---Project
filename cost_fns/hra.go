package hra

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

//behavior is either "idle", "moving" or "doorOpen"
//direction is either "up", "down" or "stop"
//cabRequests is an array of length 4, where each element is true if the corresponding button is pressed

// create a map from elevator state to behavior string

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func HRA() {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
		err := os.Chmod("./cost_fns/hall_request_assigner/"+hraExecutable, 0755)
		if err != nil {
			fmt.Println("Error setting executable permissions:", err)
			return
		}
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	mapElevStateToBehavior := map[config.ElevatorState]string{
		config.Idle:     "idle",
		config.Moving:   "moving",
		config.DoorOpen: "doorOpen",
	}

	mapQueueToCabRequests := func(queue [config.NumFloors][config.NumButtons]config.OrderState) []bool {
		cabRequests := make([]bool, config.NumFloors)
		for i := 0; i < config.NumFloors; i++ {
			if queue[i][config.ButtonCab] == config.Confirmed {
				cabRequests[i] = true
			}
		}
		return cabRequests
	}

	mapDirectionToString := func(direction elevio.MotorDirection) string {
		switch direction {
		case elevio.MD_Up:
			return "up"
		case elevio.MD_Down:
			return "down"
		case elevio.MD_Stop:
			return "stop"
		default:
			return "unknown"
		}
	}

	var hallRequests [][2]bool
	for i := 0; i < config.NumFloors; i++ {
		hallRequests = append(hallRequests, [2]bool{false, false})
		for j := 0; j < config.NumButtons-1; j++ {
			if config.ElevatorInstance.Queue[i][j] == config.Confirmed {
				hallRequests[i][j] = true
			}
		}
	}

	input := HRAInput{
		HallRequests: hallRequests,
		States:       make(map[string]HRAElevState),
	}

	for id, elev := range config.Elevators {
		input.States[id] = HRAElevState{
			Behavior:    mapElevStateToBehavior[elev.State],
			Floor:       elev.Floor,
			Direction:   mapDirectionToString(elev.Direction),
			CabRequests: mapQueueToCabRequests(elev.Queue),
		}
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		return
	}

	ret, err := exec.Command("./cost_fns/hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes), "--includeCab").CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][3]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	fmt.Println((*output)[config.ElevatorInstance.ID])

	config.MyQueue <- (*output)[config.ElevatorInstance.ID]

}
