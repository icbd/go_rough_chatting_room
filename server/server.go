package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type clientChan chan string
type server struct {
	clients    map[clientChan]bool
	entering   chan clientChan
	leaving    chan clientChan
	microphone clientChan
}

func New(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[ServerMode]\t" + addr)

	s := server{
		clients:    make(map[clientChan]bool),
		entering:   make(chan clientChan),
		leaving:    make(chan clientChan),
		microphone: make(clientChan),
	}
	go s.broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("[Err]" + err.Error())
			continue
		}
		go s.handle(conn)
	}
}

func (s *server) broadcast() {
	for {
		select {
		case c := <-s.entering:
			s.clients[c] = true
		case c := <-s.leaving:
			delete(s.clients, c)
			close(c)
		case msg := <-s.microphone:
			log.Println("[msg]" + msg)
			for c := range s.clients {
				c <- msg
			}
		}
	}
}

func (s *server) handle(conn net.Conn) {
	c := make(clientChan)
	go c.sendMsg(conn)

	who := conn.RemoteAddr().String()
	c <- "你的ID: " + who
	s.microphone <- who + "上线了!"
	s.entering <- c

	input := bufio.NewScanner(conn)
	for input.Scan() {
		s.microphone <- "\t\t" + who + ":" + input.Text()
	}
	// conn disconnected
	s.leaving <- c
	s.microphone <- who + "下线了."
	if err := conn.Close(); err != nil {
		log.Println("[Err]" + err.Error())
	}
}

func (cC clientChan) sendMsg(conn net.Conn) {
	for msg := range cC {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Println("[Err]" + err.Error())
		}
	}
}
