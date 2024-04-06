package mywbskt

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		ty, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mmsg := Message{Type: ty, Data: string(msg)}
		c.Pool.Broadcast <- mmsg
		fmt.Printf("Received message: %+v\n", mmsg)
	}
}
