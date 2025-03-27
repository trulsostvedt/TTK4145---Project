package hra

import (
	"TTK4145---project/config"
	"TTK4145---project/driver-go/elevio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
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
	if !config.IsOfflineMode {
		for id, elev := range config.Elevators {
			if elev.Floor != -1 {
				input.States[id] = HRAElevState{
					Behavior:    mapElevStateToBehavior[elev.State],
					Floor:       elev.Floor,
					Direction:   mapDirectionToString(elev.Direction),
					CabRequests: mapQueueToCabRequests(elev.Queue),
				}
			}
		}
	} else {
		input.States[config.ElevatorInstance.ID] = HRAElevState{
			Behavior:    mapElevStateToBehavior[config.ElevatorInstance.State],
			Floor:       config.ElevatorInstance.Floor,
			Direction:   mapDirectionToString(config.ElevatorInstance.Direction),
			CabRequests: mapQueueToCabRequests(config.ElevatorInstance.Queue),
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

	// Pretty print the output in the desired format
	fmt.Print("\033[H\033[2J") // Clear the terminal

	// Sort elevator IDs to ensure the lowest ID is first
	elevatorIDs := make([]string, 0, len(*output))
	for id := range *output {
		elevatorIDs = append(elevatorIDs, id)
	}
	sort.Strings(elevatorIDs)

	// Define box dimensions
	boxWidth := 25
	boxContentWidth := boxWidth - 2

	// Print top border for each elevator box
	fmt.Printf("%8s", "")
	for range elevatorIDs {
		fmt.Printf("┌%s┐", strings.Repeat("─", boxContentWidth))
	}
	fmt.Println()

	// Print elevator IDs centered in the box
	fmt.Printf("%8s", "")
	for _, id := range elevatorIDs {
		padding := (boxContentWidth - len(id)) / 2
		fmt.Printf("│%s%s%s│", strings.Repeat(" ", padding), id, strings.Repeat(" ", boxContentWidth-len(id)-padding))
	}
	fmt.Println()

	// Print elevator states centered in the box
	fmt.Printf("%8s", "")
	for _, id := range elevatorIDs {
		state := input.States[id].Behavior
		padding := (boxContentWidth - len(state)) / 2
		fmt.Printf("│%s%s%s│", strings.Repeat(" ", padding), state, strings.Repeat(" ", boxContentWidth-len(state)-padding))
	}
	fmt.Println()

	// Print elevator directions centered in the box
	fmt.Printf("%8s", "")
	for _, id := range elevatorIDs {
		direction := input.States[id].Direction
		padding := (boxContentWidth - len(direction)) / 2
		fmt.Printf("│%s%s%s│", strings.Repeat(" ", padding), direction, strings.Repeat(" ", boxContentWidth-len(direction)-padding))
	}
	fmt.Println()

	// Print middle border for each elevator box
	fmt.Printf("%8s", "")
	for range elevatorIDs {
		fmt.Printf("├%s┤", strings.Repeat("─", boxContentWidth))
	}
	fmt.Println()

	// Print "up", "down", "cab" labels centered in the box
	fmt.Printf("%8s", "")
	for range elevatorIDs {
		fmt.Printf("│%7s %7s %7s│", "up", "down", "cab")
	}
	fmt.Println()

	// Print the states for each floor
	for floor := config.NumFloors - 1; floor >= 0; floor-- {
		fmt.Printf("hall%-2d  ", floor)
		for _, id := range elevatorIDs {
			states := (*output)[id]
			fmt.Printf("│%7v %7v %7v│", states[floor][0], states[floor][1], states[floor][2]) // Box content
		}
		fmt.Println()
	}

	// Print bottom border for each elevator box
	fmt.Printf("%8s", "")
	for range elevatorIDs {
		fmt.Printf("└%s┘", strings.Repeat("─", boxContentWidth))
	}
	fmt.Println()

	config.MyQueue <- (*output)[config.ElevatorInstance.ID]

}
