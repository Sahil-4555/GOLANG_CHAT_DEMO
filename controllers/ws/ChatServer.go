package ws

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/Sahil-4555/mvc/shared/log"
)

type WsServer struct {
	clients           map[*Client]bool          // clients registered in server wsclients.go
	register          chan *Client              // request to register a client
	unregister        chan *Client              // request to unregister a client
	broadcast         chan []byte               // broadcast channel listens for messages, sent by the client readPump. It in turn pushes this messages in the send channel of all the clients registered.
	broadcastToClient chan map[*Client]*Message // broadcast to only particular client
	rooms             map[*Room]bool            // rooms created
}

// NewWebSocketServer initialize new websocket server
func NewWebSocketServer() *WsServer {
	return &WsServer{
		clients:           make(map[*Client]bool),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan []byte), //unbuffered channel unlike of send of client cause it will recieve only when readpump sends in it else it will block
		broadcastToClient: make(chan map[*Client]*Message),
		rooms:             make(map[*Room]bool),
	}
}

// Run websocket server accepting various requests
func (server *WsServer) Run() {
	log.GetLog().Info("INFO : ", "Chat server controller called(Run).")

	for {
		select {
		case client := <-server.register:
			server.registerClient(client) //add the client
			server.broadcastActiveMessage()

		case client := <-server.unregister:
			server.unregisterClient(client) //remove the client
			server.broadcastActiveMessage()

		case message := <-server.broadcast: //this broadcaster will broadcast to all clients
			server.broadcastToClients(message) //broadcast the message from readpump

		case message := <-server.broadcastToClient: //this broadcaster will broadcast to particular clients
			server.broadcastToAClient(message)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	log.GetLog().Info("INFO : ", "Chat server controller called(registerClient).")

	server.clients[client] = true
}

func (server *WsServer) broadcastActiveMessage() {
	log.GetLog().Info("INFO : ", "Chat server controller called(broadcastActiveMessage).")

	activeusers := map[string]string{
		"ActiveNow": strconv.Itoa(len(server.clients)),
	}

	jsonStr, err := json.Marshal(activeusers)
	if err != nil {
		log.GetLog().Info("ERROR(websocket) : ", err.Error())
	}

	server.broadcastToClients(jsonStr)
}

func (server *WsServer) broadcastToClients(message []byte) {
	log.GetLog().Info("INFO : ", "Chat server controller called(broadcastToClients).")

	for client := range server.clients {
		client.send <- message //Client
	}
}

func (server *WsServer) findRoom(name string) *Room {
	log.GetLog().Info("INFO : ", "Chat server controller called(findRoom).")

	var foundroom *Room
	for room := range server.rooms {
		if room.name == name {
			foundroom = room
			log.GetLog().Info("INFO : ", fmt.Sprintf("Found Room %s.", room.name))
			break
		}
	}

	return foundroom
}

func (server *WsServer) unregisterClient(client *Client) {
	log.GetLog().Info("INFO : ", "Chat server controller called(unregisterClient).")

	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToAClient(message map[*Client]*Message) {
	log.GetLog().Info("INFO : ", "Chat server controller called(broadcastToAClient).")

	for client, message := range message {
		client.send <- message.encode() //Client
	}
}

func (server *WsServer) createRoom(name string, client *Client, maxRoomCap int) *Room {
	log.GetLog().Info("INFO : ", "Chat server controller called(createRoom).")

	room := NewRoom(name, server)
	room.Grouplimit = maxRoomCap
	go room.RunRoom()
	server.rooms[room] = true
	log.GetLog().Info("INFO :", fmt.Sprintf("Started Goroutine For Room %s By %s.", room.name, client.Name))

	return room
}

func (server *WsServer) deleteRoom(room *Room) {
	log.GetLog().Info("INFO : ", "Chat server controller called(deleteRoom).")

	if _, ok := server.rooms[room]; ok {
		delete(server.rooms, room)
	}
}
