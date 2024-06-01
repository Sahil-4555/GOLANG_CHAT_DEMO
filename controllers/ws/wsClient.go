package ws

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"chat-demo-golang/configs/database"
	"chat-demo-golang/configs/middleware"
	"chat-demo-golang/models"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client represents the websocket client at the server [server] <--> [client]
type Client struct {
	Name     string          `json:"name"`      // name of the client
	Id       string          `json:"_id"`       // id of the client
	UserName string          `json:"user_name"` // username of the client
	conn     *websocket.Conn // websocket connection of the client
	wsServer *WsServer       // keep reference of webserver to every client
	send     chan []byte     // send channel then moves this to broadcast channel in server
	room     *Room           // room in which client is connected
}

// NewClient initialize new websocket client like App server in routes.go
func NewClient(wscon *websocket.Conn, wsServer *WsServer, name string, id string, username string) *Client {
	c := &Client{
		Name:     name,
		Id:       id,
		UserName: username,
		conn:     wscon,
		wsServer: wsServer,
		send:     make(chan []byte, 256), //needs to be buffered cause it should not block when channel is not receiving from broadcast
		room:     &Room{},
	}

	wsServer.register <- c
	return c
}

// handleNewMessage will handle Client messages
func (client *Client) handleNewMessage(message Message) {
	log.GetLog().Info("INFO : ", "Started handle new message server.")
	message.Sender = client

	switch message.Action {
	case SendMessageAction:
		// sending message to the users who joined the room 'roomname'
		roomname := message.Target

		// after sending the message on brodcast channel of room `roomname`, now storing that message in the DB
		id := client.StoreInDB(message)
		message.Id = id
		message.CreatedAt = time.Now()
		message.UpdatedAt = time.Now()

		var room *Room
		if room = client.wsServer.findRoom(roomname); room != nil {
			room.broadcast <- &message
		}

		// add the id of the users which are connected with the room in the message's read_by field
		client.UpdateMessageReadby(id, room)

	case JoinRoomAction:
		// to join the room
		client.HandleJoinRoomMessage(message)
		// to update the last opened channel by an client
		client.UpdateChannelStatus(message)

	case LeaveRoomAction:
		// to leave the room
		client.HandleLeaveRoomMessage(message)

	case ClientNameAction:
		clientname := message.Message
		client.Name = clientname

	case PostGlobalNotification:
		client.HandleNotification(message)

	case ChangeUserStatus:
		client.HandleUserStatus(message)

	case UpdateChannelDataAcrossChannel:
		client.UpdateChannelDataAcrossChannel(message)
	}
}

func (client *Client) HandleUserStatus(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(HandleUserStatus).")

	switch message.Status {
	case common.USER_STATUS_ONLINE:
		client.ChangeUserStatusToOnline(message)

	case common.USER_STATUS_OFFLINE:
		client.ChangeUserStatusToOffline(message)

	case common.USER_STATUS_AWAY:
		client.ChangeUserStatusToAway(message)

	case common.USER_STATUS_DO_NOT_DISTURB:
		client.ChangeUserStatusToDoNotDisturb(message)

	}
}

func (client *Client) UpdateChannelDataAcrossChannel(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(UpdateChannelDataAcrossChannel).")
	var room *Room

	if room = client.wsServer.findRoom(message.Target); room != nil {
		room.broadcast <- &message
	} else {
		log.GetLog().Info("ERROR : ", "Room not found.")
	}
}

func (client *Client) ChangeUserStatusToOnline(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(ChangeUserStatusToOnline).")

	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(client.Id)
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: common.USER_STATUS_ONLINE}}}}

	_, err := conn.UserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	m := &Message{
		Action:  ChangeUserStatus,
		Payload: message.Payload,
	}

	msg := m.encode()
	client.wsServer.broadcast <- msg
}

func (client *Client) ChangeUserStatusToOffline(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(ChangeUserStatusToOffline).")

	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(client.Id)
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: common.USER_STATUS_OFFLINE}}}}

	_, err := conn.UserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	m := &Message{
		Action:  ChangeUserStatus,
		Payload: message.Payload,
	}

	msg := m.encode()
	client.wsServer.broadcast <- msg
}

func (client *Client) ChangeUserStatusToDoNotDisturb(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(ChangeUserStatusToDoNotDisturb).")

	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(client.Id)
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: common.USER_STATUS_DO_NOT_DISTURB}}}}

	_, err := conn.UserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	m := &Message{
		Action:  ChangeUserStatus,
		Payload: message.Payload,
	}

	msg := m.encode()
	client.wsServer.broadcast <- msg
}

func (client *Client) ChangeUserStatusToAway(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(ChangeUserStatusToAway).")

	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, _ := primitive.ObjectIDFromHex(client.Id)
	filter := bson.M{"_id": id, "deleted_at": nil}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: common.USER_STATUS_AWAY}}}}

	_, err := conn.UserCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	m := &Message{
		Action:  ChangeUserStatus,
		Payload: message.Payload,
	}

	msg := m.encode()
	client.wsServer.broadcast <- msg
}

func (client *Client) UpdateChannelStatus(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(UpdateChannelStatus).")

	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	channId, _ := primitive.ObjectIDFromHex(message.Message)
	senderId, _ := primitive.ObjectIDFromHex(client.Id)

	// to add and update the last opened channel by an client
	var ch models.Channel
	filter := bson.M{"_id": channId, "last_opened.user_id": bson.M{"$in": bson.A{senderId}}}
	err := conn.ChannelCollection().FindOne(ctx, filter).Decode(&ch)

	if err == mongo.ErrNoDocuments {
		var d models.LastOpenedBy
		d.UserId = senderId
		d.TimeStamp()
		filter := bson.M{"_id": channId}
		update := bson.M{"$push": bson.M{"last_opened": d}}
		_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			m := &Message{
				Action:  Error,
				Message: err.Error(),
			}
			client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
		}
	} else if err == nil {
		filter := bson.M{"_id": channId, "last_opened.user_id": senderId}
		update := bson.M{"$set": bson.M{"last_opened.$.last_opened_at": time.Now()}}
		_, err := conn.ChannelCollection().UpdateOne(ctx, filter, update)
		if err != nil {
			m := &Message{
				Action:  Error,
				Message: err.Error(),
			}
			client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
		}
	} else {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	// to update the messages to seen which are still unseen by client.
	filter = bson.M{"channel_id": channId, "read_by": bson.M{"$ne": senderId}}
	update := bson.M{"$addToSet": bson.M{"read_by": senderId}}

	_, err = conn.MessageCollection().UpdateMany(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

}

func (client *Client) UpdateMessageReadby(id primitive.ObjectID, room *Room) {
	log.GetLog().Info("INFO : ", "Client controller called(UpdateMessageReadby).")

	var users []primitive.ObjectID
	for key := range room.clients {
		_, ok := room.clients[key]
		if ok {
			userId, _ := primitive.ObjectIDFromHex(key.Id)
			users = append(users, userId)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$addToSet": bson.M{"read_by": bson.M{"$each": users}}}
	_, err := database.NewConnection().MessageCollection().UpdateOne(ctx, filter, update)

	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}
}

func (client *Client) StoreInDB(message Message) primitive.ObjectID {
	log.GetLog().Info("INFO : ", "Client controller called(StoreInDB).")
	conn := database.NewConnection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	channId, _ := primitive.ObjectIDFromHex(message.Target)
	senderId, _ := primitive.ObjectIDFromHex(client.Id)

	var msg models.Message
	msg.Sender = senderId
	msg.ChannelId = channId
	msg.ContentType = message.ContentType
	msg.ReadBy = []primitive.ObjectID{senderId}

	if message.ContentType == common.CONTENT_TYPE_TEXT {
		msg.Content = message.Message
	} else {
		msg.FileName = message.FileName
	}

	msg.TimeStamp()
	msg.NewID()

	result, err := conn.MessageCollection().InsertOne(ctx, msg)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: "Failed To Store The Message In DB",
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	filter := bson.M{"_id": channId, "deleted_at": bson.M{"$eq": nil}}
	update := bson.M{"$set": bson.M{"last_activity.$[].last_activity_at": time.Now()}}
	_, err = conn.ChannelCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		m := &Message{
			Action:  Error,
			Message: err.Error(),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: m}
	}

	var channel models.Channel
	conn.ChannelCollection().FindOne(ctx, filter).Decode(&channel)
	var reciver primitive.ObjectID
	if channel.ChannelType == common.ONE_TO_ONE_COMMUNICATION {
		for index := range channel.Users {
			if channel.Users[index] != senderId {
				reciver = channel.Users[index]
				break
			}
		}
		update = bson.M{"$pull": bson.M{"close_conversation": bson.M{"$in": bson.A{reciver}}}}
		conn.ChannelCollection().UpdateOne(ctx, filter, update)
	}

	return result.InsertedID.(primitive.ObjectID)
}

func (client *Client) HandleNotification(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(HandleNotification).")
	user := make(map[string]bool)

	for _, id := range message.NotifyUsers {
		user[id.Hex()] = true
	}

	// to send the message to the users who are in the room with event type
	var msg Message
	switch message.NotificationType {
	case MetionInMessageNotification:
		msg.Action = MetionInMessageNotification
		msg.ChannelId = message.Target
		msg.Message = message.Message
	case UpdateChannelDataAcrossUser:
		msg.Action = UpdateChannelDataAcrossUser
		msg.Target = message.Target
		msg.Payload = message.Payload
	case UpdateChannelDataAcrossChannel:
		msg.Action = UpdateChannelDataAcrossChannel
		msg.Target = message.Target
	case UpdateChannelOnMessage:
		msg.Action = UpdateChannelOnMessage
		msg.Payload = message.Payload
	case UpdateMessage:
		msg.Action = UpdateMessage
		msg.Payload = message.Payload
	case DeleteMessage:
		msg.Action = DeleteMessage
		msg.Payload = message.Payload
	case AddChannelOnAddingMember:
		msg.Action = AddChannelOnAddingMember
		msg.Target = message.Target
		msg.Payload = message.Payload
	}

	for clients := range client.wsServer.clients {
		_, ok := user[clients.Id]

		if ok {
			client.wsServer.broadcastToClient <- map[*Client]*Message{clients: &msg}
		}
	}
}

// handleLeaveRoomMessage will handle leaving a room.unregister <- Client
func (client *Client) HandleLeaveRoomMessage(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(HandleLeaveRoomMessage).")
	client.room.unregister <- client
}

// handleJoinRoomMessage will handle joining a room.register <- Client
func (client *Client) HandleJoinRoomMessage(message Message) {
	log.GetLog().Info("INFO : ", "Client controller called(HandleJoinRoomMessage).")

	if client.room.name != "" {
		log.GetLog().Info("INFO : ", fmt.Sprintf("%s client leave the room %s.\n", client.Name, client.room.name))
		client.room.unregister <- client
	}

	roomname := message.Message
	room := client.wsServer.findRoom(roomname)
	if room == nil {
		if len(roomname) != 24 {
			message := &Message{
				Action:  FailJoinRoomNotification,
				Target:  roomname,
				Message: fmt.Sprintf("%s room name is not valid.", roomname),
			}
			client.wsServer.broadcastToClient <- map[*Client]*Message{client: message}
			return
		}
		room = client.wsServer.createRoom(roomname, client, MAX_ROOM_CAPACITY)
	}

	if len(room.clients) < room.Grouplimit {
		client.room = room
		room.register <- client
	} else {
		message := &Message{
			Action:  FailJoinRoomNotification,
			Target:  room.name,
			Message: fmt.Sprintf("oops the room %v is occupied :-P", room.name),
		}
		client.wsServer.broadcastToClient <- map[*Client]*Message{client: message}
	}
}

// disconnect will handle the client disconnection from the server
func (client *Client) disconnect() {
	log.GetLog().Info("INFO : ", "Client controller called(disconnect).")

	msg := Message{
		Status: common.USER_STATUS_OFFLINE,
		Payload: map[string]interface{}{
			"user_id": client.Id,
			"status":  common.USER_STATUS_OFFLINE,
		},
	}

	client.ChangeUserStatusToOffline(msg)

	client.wsServer.unregister <- client //remove client from webserver map list
	if _, ok := client.wsServer.rooms[client.room]; ok {
		client.room.unregister <- client //unregister to room
	}

	close(client.send)  //close the sending channel
	client.conn.Close() //close the client connection
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	log.GetLog().Info("INFO : ", "Client controller called(ServeWs).")
	conn, err := websocket.Upgrade(w, r, nil, 4096, 4096)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return
	}

	// ws://url?_id=id&token=token
	id := r.URL.Query().Get("_id")
	user := middleware.GetUserById(id)
	client := NewClient(conn, wsServer, user.Name, id, user.UserName)

	log.GetLog().Info("INFO : ", fmt.Sprintf("%s Client Joined The Hub", client.Name))

	go client.readPump()
	go client.writePump()
}

// writePump goroutine handles sending the messages to the connected client. It runs in an endless loop waiting for new messages in the client.send channel. When receiving new messages it writes them to the client, if there are multiple messages available they will be combined in one write.
// writePump is also responsible for keeping the connection alive by sending ping messages to the client with the interval given in pingPeriod. If the client does not respond with a pong, the connection is closed.
func (client *Client) writePump() {
	log.GetLog().Info("INFO : ", "Client controller called(writePump).")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.GetLog().Info("ERROR : ", err.Error())
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				log.GetLog().Info("ERROR : ", err.Error())
				return
			}

		case <-ticker.C: //make a ping request
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.GetLog().Info("ERROR : ", err.Error())
				return
			}
		}
	}
}

// readPump Goroutine, the client will read new messages send over the WebSocket connection. It will do so in an endless loop until the client is disconnected. When the connection is closed, the client will call its own disconnect method to clean up.
// upon receiving new messages the client will push them in the WsServer broadcast channel.
func (client *Client) readPump() {
	log.GetLog().Info("INFO : ", "Client controller called(readPump).")
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)

	// Frontend client will give a pong message to the routine we have to handle it
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(appData string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for message from client
	for {
		var message Message
		err := client.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.GetLog().Info("ERROR : ", fmt.Sprintf("Unexepected Close Error: %v.\n", err))
				break
			}

			log.GetLog().Info("ERROR : ", err.Error())
			room := client.room
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				if len(room.clients) == 0 {
					room.stoproom <- true
				}
			}
			break
		}

		client.handleNewMessage(message) //broadcast to room

	}
}
