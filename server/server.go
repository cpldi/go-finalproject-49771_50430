package main

import (
	msg "bitcoin_miner/message"
	"bitcoin_miner/server/miner"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func startServer(port int) error {
	address := fmt.Sprintf(":%v", port)
	LOGF.Printf("Starting server at %v\n", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	m := miner.NewMiner(10, 10, 10, 10)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, m)
	}
}

func handleConnection(conn net.Conn, m *miner.Miner) {
	defer conn.Close()

	req, err := msg.FromJSONReader(conn)

	if err != nil && isCorrect(req) {
		fmt.Println(err)
		return
	}

	LOGF.Printf("Submiting Job %v \n", req)
	resp := m.SubmitJob(req)

	respJson, err := resp.ToJSON()

	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Write(respJson)
}

func isCorrect(req *msg.Message) bool {
	return req.Type == msg.Request &&
		req.Upper >= req.Lower &&
		len(req.Data) > 0
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
}
