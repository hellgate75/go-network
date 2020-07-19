package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-api-client" {
			TestApiClient()
		} else if os.Args[1] == "-api-server" {
			TestApiServer()
		} else if os.Args[1] == "-tcp-client" {
			TestTcpClient()
		} else if os.Args[1] == "-tcp-server" {
			TestTcpServer()
		} else if os.Args[1] == "-pipe-in" {
			TestInputPipeNode()
		} else if os.Args[1] == "-pipe-out" {
			TestOutputPipeNode()
		} else {
			fmt.Printf("Unknwon argument: %s, accepted: -api-client or -api-server or -tcp-client or -tcp-server or -pipe-in or -pipe-out\n", os.Args[1])
		}
	} else {
		fmt.Println("Not enough arguments use -client or -server")
	}
}
