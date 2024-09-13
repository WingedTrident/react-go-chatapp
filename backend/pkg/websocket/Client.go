package websocket

import (
	"fmt"
	"log"
	//"sync"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type Client struct {
	ID string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int `json:"type"`
	Body string `json:"body"`
	Chatters []string `json:"chatters"`
	Sender string `json:"sender"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()
	c.ID = uuid.New().String()	//generates a unique ID
	fmt.Println("NEW ID", c.ID)
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p), Sender: c.ID}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}