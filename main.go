package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	go Start()
	listenForInterrupt()
}

func listenForInterrupt() {
	fmt.Println("Listening for interrupt.")
	// Create a channel which transfers objects of type os.Signal
	signalChan := make(chan os.Signal, 1)

	// signal.Notify writes to signalChan when an os.Interrupt signal is detected.
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println("Received interrupt signal. Stopping operation.")
}

// Task represents a concurrent task. It has a channel which it listens on for when to gracefully exit.
type Task struct {
	closed chan struct{}
}
