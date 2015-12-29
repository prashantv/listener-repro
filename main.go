package main

import (
	"fmt"
	"net"
	"time"
)

const (
	connectionsBeforeClose = 5
	testIterations         = 1000
	enableSleep            = false
)

func runTest() error {
	addr, ln := startListener()

	for i := 0; i < connectionsBeforeClose; i++ {
		if err := connect(addr); err != nil {
			panic(err)
		}
	}
	ln.Close()

	if enableSleep {
		time.Sleep(time.Millisecond)
	}

	if err := connect(addr); err == nil {
		return fmt.Errorf("connect succeeded even though it shouldn't have")
	}
	return nil
}

func main() {
	failures := 0
	for i := 0; i < testIterations; i++ {
		if err := runTest(); err != nil {
			fmt.Println("Iteration", i, "failed with", err)
			failures++
		}
	}
	fmt.Printf("Got %v failures out of %v\n", failures, testIterations)
}

func startListener() (string, net.Listener) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go acceptLoop(ln)

	addr := ln.Addr().String()
	fmt.Println("TCP Listener started on", addr)
	return addr, ln
}

func connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		conn.Close()
	}
	return err
}

func acceptLoop(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		conn.Close()
	}
}
