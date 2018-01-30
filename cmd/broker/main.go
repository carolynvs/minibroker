package main

import (
	"fmt"

	"github.com/deis/incluster-broker/broker"
	"github.com/deis/incluster-broker/server"
)

func main() {
	fmt.Println("Hi, I'm an in-cluster broker!")
	inclusterBroker := broker.InClusterBroker{}
	server.Run(inclusterBroker)
}
