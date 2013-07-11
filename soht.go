package main

import (
	"flag"
	"fmt"
	"log"
)

var serverflag bool
var clientflag bool
var debugFlag bool
var serverAddr string

type NoopWriter struct {
}

func (w NoopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func init() {
	flag.BoolVar(&serverflag, "server", false, "Run as server")
	flag.BoolVar(&clientflag, "client", false, "Run as client")
	flag.BoolVar(&debugFlag, "debug", false, "Output Debug Log to stderr")
	flag.StringVar(&serverAddr, "addr", ":8080", "Address to listen on")
}

func main() {
	flag.Parse()
	if !debugFlag {
		noopWriter := NoopWriter{}
		log.SetOutput(noopWriter)
	}
	log.Println("Debug Output Enabled")
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
