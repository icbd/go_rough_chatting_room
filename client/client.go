package client

import (
	"io"
	"log"
	"net"
	"os"
)

type client struct {
	conn    net.Conn
	allDown chan struct{}
}

func New(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[ClientMode]\t" + addr)
	c := client{
		conn:    conn,
		allDown: make(chan struct{}),
	}

	go c.printMsgFromConn()
	c.typeMsgToConn()
	<-c.allDown
}

func (c *client) typeMsgToConn() {
	if _, err := io.Copy(c.conn, os.Stdin); err != nil {
		log.Fatal(err)
	}
	// stdin closed
	if err := c.conn.Close(); err != nil {
		log.Fatal(err)
	}
}

func (c *client) printMsgFromConn() {
	_, err := io.Copy(os.Stdout, c.conn)
	// conn disconnected
	if opErr, ok := err.(*net.OpError); ok {
		if opErr.Err.Error() == "use of closed network connection" {
			log.Println("BYE.")
			c.allDown <- struct{}{}
			return
		}
	}
	// other fatal errors
	log.Fatal(err)
}
