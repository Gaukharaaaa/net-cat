package main

import (
	"fmt"
	"time"
)

type Message struct {
	Author string
	Text   string
	Time   time.Time
}

func (m Message) String() string {
	t := m.Time.Format("2006-01-02 15:04:05 07:00")
	return fmt.Sprintf("[%s][%s]:%s", t, m.Author, m.Text)
}

func (m Message) StatusString() string {
	return fmt.Sprintf("%s %s", m.Author, m.Text)
}
