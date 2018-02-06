package main

import (
	"fmt"

	"github.com/deis/minibroker/broker"
	"github.com/deis/minibroker/server"
)

func main() {
	fmt.Println("Hi, I'm an in-cluster broker!")
	minibroker := broker.Minibroker{}
	server.Run(minibroker)
}
