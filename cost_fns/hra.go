package hra

import (
	"TTK4145---project/Driver-go/elevio"
	"TTK4145---project/config"
	"encoding/json"
	"fmt"
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

	input := HRAInput{
		HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
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

	ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
}
