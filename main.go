package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	port     string
	logo     string
	messChan chan Message
	statChan chan Message
)

const (
	yourName = "[ENTER YOUR NAME]:"
	joinChat = "has joined our chat..."
	leftChat = "has left our chat..."
)

func init() {
	l, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Fatalln(err)
	}

	logo = string(l)
	messChan = make(chan Message)
	statChan = make(chan Message)
}

func main() {
	flag.StringVar(&port, "port", "8989", "Port for Net-Cat.")
	flag.Parse()

	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	sess := &Users{
		All: map[net.Conn]string{},
		mu:  sync.Mutex{},
	}

	hist := &History{
		mu: sync.Mutex{},
	}

	go BroadCast(sess, hist)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go Handle(conn, sess, hist)

	}
}

func ScanName(conn net.Conn, u *Users) string {
	var username string
	scanner := bufio.NewScanner(conn)
	fmt.Fprint(conn, yourName)
	for scanner.Scan() {
		username = scanner.Text()
		if IsEmptyMess(username) || !IsPrintable(username) || u.ContainName(username) {
			fmt.Fprintln(conn, "incorrect name or user is already existed")
			fmt.Fprint(conn, yourName)
			continue
		}
		break
	}

	return username
}

func Handle(conn net.Conn, u *Users, h *History) {
	defer conn.Close()
	var username string

	if u.IsOverflow() {
		fmt.Fprintln(conn, "room is overflow")
		return
	}

	fmt.Fprint(conn, logo)
	username = ScanName(conn, u)
	u.Add(username, conn)
	defer u.Del(conn)
	statChan <- Message{
		Author: username,
		Text:   joinChat,
	}

	h.PrintTo(conn)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		msg := Message{
			Author: username,
			Time:   time.Now(),
		}

		msg.Text = scanner.Text()
		messChan <- msg
	}

	statChan <- Message{
		Author: username,
		Text:   leftChat,
	}
}

func BroadCast(u *Users, h *History) {
	for {
		select {
		case msg := <-messChan:
			h.Add(msg.String())

			u.mu.Lock()
			for conn, name := range u.All {
				if name != msg.Author {
					fmt.Fprint(conn, "\n", msg, "\n")
				}
				fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format("2006-01-02 15:04:05 07:00"), name)
			}
			u.mu.Unlock()

		case stt := <-statChan:
			h.Add(stt.StatusString())

			u.mu.Lock()
			for conn, name := range u.All {
				if name != stt.Author {
					fmt.Fprint(conn, "\n", stt.StatusString(), "\n")
				}
				fmt.Fprintf(conn, "[%s][%s]:", time.Now().Format("2006-01-02 15:04:05 07:00"), name)
			}
			u.mu.Unlock()

		}
	}
}
