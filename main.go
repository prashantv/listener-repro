package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	connectionsBeforeClose = 1
	testIterations         = 10000
	enableSleep            = false
	useSaneListener        = true
)

func runTest() error {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	if useSaneListener {
		ln = NewSaneListener(ln)
	}

	addr := ln.Addr().String()
	fmt.Println("Listener started on", addr)

	waitForListener := make(chan error)
	go func() {
		defer close(waitForListener)

		var connCount int
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			connCount++
			if connCount > connectionsBeforeClose {
				waitForListener <- errors.New("got unexpected conn")
				return
			}
			conn.Close()
		}
	}()

	for i := 0; i < connectionsBeforeClose; i++ {
		if _, _, err := connect(addr); err != nil {
			panic(err)
		}
	}

	ln.Close()

	if enableSleep {
		time.Sleep(time.Millisecond)
	}

	connect(addr)

	err, _ = <-waitForListener
	return err
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

func connect(addr string) (localAddr string, remoteAddr string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		localAddr = conn.LocalAddr().String()
		remoteAddr = conn.RemoteAddr().String()
		conn.Close()
	}
	return localAddr, remoteAddr, err
}
