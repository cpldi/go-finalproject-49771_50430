package main

import (
	msg "bitcoin_miner/message"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	const numArgs = 4
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <host:port> <message> <max>", os.Args[0])
		return
	}
	host := os.Args[1]
	message := os.Args[2]
	max, err := strconv.ParseUint(os.Args[3], 10, 64)
	if err != nil {
		fmt.Printf("%s is not a number.\n", os.Args[3])
		return
	}

	response, err := sendRequest(message, host, max)
	if err != nil {
		fmt.Printf("Error parsing response : %v  ", err)
		return
	}
	printResult(response.Hash, response.Nonce)
}

func sendRequest(message string, host string, max uint64) (*msg.Message, error) {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, err
	}

	req := msg.NewRequest(message, 0, max)
	jsonb, err := req.ToJSON()
	_, err = conn.Write(jsonb)

	if err != nil {
		return nil, err
	}

	response, err := msg.FromJSONReader(conn)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// printResult prints the final result to stdout.
func printResult(hash, nonce uint64) {
	fmt.Println("Result", hash, nonce)
}

// printDisconnected prints a disconnected message to stdout.
func printDisconnected() {
	fmt.Println("Disconnected from the server.")
}
