package main

import (
	"net"
	"sync"
)

// NewSaneListener returns a new SaneListener around the provided net.Listener.
func NewSaneListener(l net.Listener) net.Listener {
	return &SaneListener{l: l, c: sync.NewCond(&sync.Mutex{})}
}

// SaneListener wraps a net.Listener and ensures that once SaneListener.Close
// returns the underlying socket has been closed.
type SaneListener struct {
	l        net.Listener
	c        *sync.Cond
	refCount int
}

func (s *SaneListener) incRef() {
	s.c.L.Lock()
	s.refCount++
	s.c.L.Unlock()
}

func (s *SaneListener) decRef() {
	s.c.L.Lock()
	s.refCount--
	s.c.Broadcast()
	s.c.L.Unlock()
}

// Accept waits for and returns the next connection to the listener.
func (s *SaneListener) Accept() (net.Conn, error) {
	s.incRef()
	defer s.decRef()
	return s.l.Accept()
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (s *SaneListener) Close() error {
	err := s.l.Close()
	if err == nil {
		s.c.L.Lock()
		for s.refCount > 0 {
			s.c.Wait()
		}
		s.c.L.Unlock()
	}
	return err
}

// Addr returns the listener's network address.
func (s *SaneListener) Addr() net.Addr {
	return s.l.Addr()
}
