package ws

import (
	"fmt"

	"github.com/Sahil-4555/mvc/shared/log"
)

// var UnregisterDone = make(chan struct{})

type Room struct {
	name                    string
	clients                 map[*Client]bool
	wsServer                *WsServer //keep reference of webserver to every client
	register                chan *Client
	unregister              chan *Client
	broadcast               chan *Message //message send in a room
	broadcastToClientinRoom chan *Message //message to a client
	Grouplimit              int
	stoproom                chan bool //end the server routine
}

func NewRoom(name string, server *WsServer) *Room {
	return &Room{
		name:                    name,
		wsServer:                server,
		clients:                 make(map[*Client]bool),
		register:                make(chan *Client),
		unregister:              make(chan *Client),
		broadcast:               make(chan *Message), //unbuffered channel unlike of send of client cause it will recieve only when readpump sends in it else it will block
		broadcastToClientinRoom: make(chan *Message),
		stoproom:                make(chan bool),
	}
}

// Run websocket server accepting various requests
func (room *Room) RunRoom() {
	log.GetLog().Info("INFO : ", "Room controller called(RunRoom).")
	for {
		select {
		case client := <-room.register:
			log.GetLog().Info("INFO : ", fmt.Sprintf("Client name %s registered a room %s.\n", client.Name, room.name))

			room.registerClientinRoom(client) //add the client
			room.notifyClientJoined(client)   // notify to clients in room

		case client := <-room.unregister:
			// Remove the client
			room.unregisterClientinRoom(client)

			//clients in room notification
			var clientlist []*Client
			for key := range room.clients {
				clientlist = append(clientlist, key)
			}
			message := &ClientsinRoomMessage{
				Action:     ClientListNotification,
				Target:     room.name,
				ClientList: clientlist,
			}
			room.broadcastToClientsInRoom(message.encode())
			if len(room.clients) == 0 {
				log.GetLog().Info("INFO : ", fmt.Sprintf("Room Shutdown %s", room.name))
				room.wsServer.deleteRoom(room)
				return
			}

		case message := <-room.broadcast:
			//broadcast the message from readpump to a room clients only
			room.broadcastToClientsInRoom(message.encode())

		case <-room.stoproom:
			log.GetLog().Info("INFO : ", fmt.Sprintf("Room Shutdown %s", room.name))
			room.wsServer.deleteRoom(room)
			return
		}
	}
}

func (room *Room) registerClientinRoom(client *Client) {
	log.GetLog().Info("INFO : ", "Room controller called(registerClientinRoom).")
	room.clients[client] = true
}

func (room *Room) notifyClientJoined(client *Client) {
	log.GetLog().Info("INFO : ", "Room controller called(notifyClientJoined).")
	message := &Message{
		Action:  JoinRoomNotification,
		Target:  room.name,
		Message: fmt.Sprintf("%s joined the room.", client.Name),
		Sender:  client,
	}

	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	log.GetLog().Info("INFO : ", "Room controller called(broadcastToClientsInRoom).")

	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) unregisterClientinRoom(client *Client) {
	log.GetLog().Info("INFO : ", "Room controller called(unregisterClientinRoom).")

	delete(room.clients, client)
}
