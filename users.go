package main

import (
	"net"
	"sync"
)

type Users struct {
	All map[net.Conn]string
	mu  sync.Mutex
}

func (u *Users) Add(name string, conn net.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.All[conn] = name
}

func (u *Users) Del(conn net.Conn) {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.All, conn)
}

// func (u *Users) ContainConn(conn net.Conn) bool {
// 	u.mu.Lock()
// 	defer u.mu.Unlock()

// 	if _, ok := u.All[conn]; ok {
// 		return true
// 	}

// 	return false
// }

func (u *Users) ContainName(name string) bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	for _, val := range u.All {
		if val == name {
			return true
		}
	}

	return false
}

func (u *Users) IsOverflow() bool {
	u.mu.Lock()
	defer u.mu.Unlock()

	if len(u.All) >= 10 {
		return true
	}

	return false
}
