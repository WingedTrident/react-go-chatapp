package websocket

import (
	"fmt"
	"strconv"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	clientNames := []string{}
	idToName := map[string]string{}
	visitorNumber := 0
	for {
		select {
		case client := <-pool.Register:
			visitorNumber++
			clientNames = append(clientNames, "Chatter"+strconv.Itoa(visitorNumber))
			idToName[client.ID] = "Chatter" + strconv.Itoa(visitorNumber)
			joinedChatter := idToName[client.ID]

			fmt.Println(idToName)
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 2, Body: joinedChatter + " Joined", Chatters: clientNames})
			}
			fmt.Println(pool.Clients)
		case client := <-pool.Unregister:
			var indexToDelete int
			leftChatter := idToName[client.ID]

			for i, item := range clientNames {
				if item == idToName[client.ID] {
					indexToDelete = i
					break
				}
			}

			if len(clientNames) == 1 {
				clientNames = []string{}
			} else {
				fmt.Println(indexToDelete, "index")
				clientNames = append(clientNames[:indexToDelete], clientNames[indexToDelete+1:]...)
			}

			delete(idToName, client.ID)
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 2, Body: leftChatter + " Disconnected", Chatters: clientNames})
			}
		case message := <-pool.Broadcast:
			sender := idToName[message.Sender]
			message.Body = sender + ": " + message.Body
			message.Chatters = clientNames
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
