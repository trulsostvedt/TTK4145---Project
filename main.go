package main

import (
	network "TTK4145---project/Network-go"
	"flag"
	"fmt"
)

func main() {

	flag.Parse()
	fmt.Println("Hello World")
	network.Network()

	select {}
}
