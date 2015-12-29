package main

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

const (
	connectionsBeforeClose = 1
	testIterations         = 1000
	enableSleep            = false
)

type test struct {
	ln             net.Listener
	addr           string
	gotConnections int32
}

func (t *test) startListener() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	t.ln = ln
	t.addr = ln.Addr().String()
	go t.acceptLoop()

	fmt.Println("TCP Listener started on", t.addr)
}

func (t *test) acceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			return
		}

		go t.handleConn(conn)
	}
}

func (t *test) handleConn(conn net.Conn) {
	if atomic.AddInt32(&t.gotConnections, 1) > connectionsBeforeClose {
		fmt.Println("  got unexpected conn with local addr", conn.LocalAddr(), "remote", conn.RemoteAddr())
	}
	conn.Close()
}

func runTest() error {
	t := &test{}
	t.startListener()

	for i := 0; i < connectionsBeforeClose; i++ {
		if _, _, err := connect(t.addr); err != nil {
			panic(err)
		}
	}

	t.ln.Close()

	if enableSleep {
		time.Sleep(time.Millisecond)
	}

	if laddr, raddr, err := connect(t.addr); err == nil {
		if false {
			return fmt.Errorf("connect succeeded even though it shouldn't have. local: %v remote %v",
				laddr, raddr)
		}
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

func connect(addr string) (localAddr string, remoteAddr string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		localAddr = conn.LocalAddr().String()
		remoteAddr = conn.RemoteAddr().String()
		conn.Close()
	}
	return localAddr, remoteAddr, err
}
