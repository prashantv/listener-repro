package main

import (
	"fmt"
	"net"
	"time"
)

const (
	connectionsBeforeClose = 1
	testIterations         = 10
	enableSleep            = false
)

func runTest() error {
	addr, ln := startListener()

	for i := 0; i < connectionsBeforeClose; i++ {
		if _, _, err := connect(addr); err != nil {
			panic(err)
		}
	}
	ln.Close()

	if enableSleep {
		time.Sleep(time.Millisecond)
	}

	if laddr, raddr, err := connect(addr); err == nil {
		return fmt.Errorf("connect succeeded even though it shouldn't have. local: %v remote %v",
			laddr, raddr)
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

func connect(addr string) (localAddr string, remoteAddr string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		localAddr = conn.LocalAddr().String()
		remoteAddr = conn.RemoteAddr().String()
		conn.Close()
	}
	return localAddr, remoteAddr, err
}

func acceptLoop(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		fmt.Println("  got conn with local addr", conn.LocalAddr(), "remote", conn.RemoteAddr())
		conn.Close()
	}
}
