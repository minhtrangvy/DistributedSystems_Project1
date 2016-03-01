package main

import (
	"fmt"

	"github.com/minhtrangvy/DistributedSystems_Project1/project1"
)

const defaultPort = 9999

func main() {
	// Initialize the server.
	server := project1.New()
	if server == nil {
		fmt.Println("New() returned a nil server. Exiting...")
		return
	}

	// Start the server and continue listening for client connections in the background.
	if err := server.Start(defaultPort); err != nil {
		fmt.Printf("MultiEchoServer could not be started: %s\n", err)
		return
	}

	fmt.Printf("Started MultiEchoServer on port %d...\n", defaultPort)

	// Block forever.
	select {}
}
