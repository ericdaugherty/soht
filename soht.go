package main

import (
	"flag"
	"fmt"
)

var serverflag bool
var clientflag bool
var serverAddr string

func init() {
	flag.BoolVar(&serverflag, "server", false, "Run as server")
	flag.BoolVar(&clientflag, "client", false, "Run as client")
	flag.StringVar(&serverAddr, "addr", ":8080", "Address to listen on")
}

func main() {
	flag.Parse()
	if serverflag && clientflag {
		fmt.Println("Please select only -client or -server, not both.")
		return
	}
	if !serverflag && !clientflag {
		fmt.Println("Please select -client or -server mode.")
		return
	}
	fmt.Printf("Welcome to Socket Over HTTP (SOHT)\n")
	if serverflag {
		fmt.Println("Running in server mode.")
		runServer(serverAddr)
	}
	if clientflag {
		fmt.Println("Running in client mode.")
		client()
	}
}

func client() {
}
