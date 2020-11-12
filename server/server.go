package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func startServer(port int) {
	// TODO: implement this!

	return //Should probably return something actually
}

var LOGF *log.Logger

func main() {
	// You may need a logger for debugging
	const (
		name = "log.txt"
		flag = os.O_RDWR | os.O_CREATE
		perm = os.FileMode(0666)
	)

	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}
	defer file.Close()

	LOGF = log.New(file, "", log.Lshortfile|log.Lmicroseconds)
	// Usage: LOGF.Println() or LOGF.Printf()

	const numArgs = 2
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <port>", os.Args[0])
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Port must be a number:", err)
		return
	}

	startServer(port)
	fmt.Println("Server listening on port", port)

	// TODO: implement this!
}
