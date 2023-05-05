package main

import (
	"fmt"
	"net"
	"sync"
)

type History struct {
	Container []string
	mu        sync.Mutex
}

func (h *History) Add(mess string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Container = append(h.Container, mess)
}

func (h *History) PrintTo(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, msg := range h.Container {
		fmt.Fprintln(conn, msg)
	}
}
