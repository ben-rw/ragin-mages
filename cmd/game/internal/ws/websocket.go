package ws

import (
	"context"
	"log"
	"syscall/js"

	"github.com/ben-rw/ragin-mages/internal/protocol"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Connection struct {
	*websocket.Conn
	incomingMsgs chan protocol.Message
	outgoingMsgs chan protocol.Message
}

// marshal msgdata to Message, add to outgoingMsgs
func (c *Connection) WriteMsg(mt protocol.MessageType, msgData any) {
	msg, err := protocol.MarshalToMessage(mt, msgData)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case c.outgoingMsgs <- *msg:
	default:
		log.Println("write buffer full: dropped a message")
	}
}

// write outgoingMsgs to server
func (c *Connection) writeLoop() {
	for {
		msg := <-c.outgoingMsgs
		err := wsjson.Write(context.Background(), c.Conn, msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				c.Close(websocket.StatusNormalClosure, "connection closed normally")
				return
			} else {
				log.Printf("wsjson write error: %v\n", err)
				log.Println("terminating write loop")
				return
			}
		}
	}
}

// check for messages to pass to Update in main
func (c *Connection) Check() []protocol.Message {
	msgs := []protocol.Message{}
	for {
		select {
		case msg := <-c.incomingMsgs:
			msgs = append(msgs, msg)
		default:
			return msgs
		}
	}
}

// read incoming messages, close dead connections, drop messages if buffer fills up
func (c *Connection) readLoop() {
	defer c.Close(websocket.StatusInternalError, "connection closed unexpectedly")

	for {
		msg := protocol.Message{}
		err := wsjson.Read(context.Background(), c.Conn, &msg)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				c.Close(websocket.StatusNormalClosure, "connection closed normally")
				return
			} else {
				log.Printf("wsjson read error: %v\n", err)
				log.Println("terminating read loop")
				return
			}
		}
		select {
		case c.incomingMsgs <- msg:
		default:
			log.Println("read buffer full: dropped a message")
		}
	}

}

// set up websocket connection, send roomID to server, start read loop
func ConnectToWebsocket() (*Connection, error) {
	ctx := context.Background()

	prot, host, roomID := getClientInfo()
	wsURL := prot + host + "/room/" + roomID + "/ws"

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		return nil, err
	}

	joinData := protocol.JoinRequestData{
		RoomID: roomID,
	}

	msg, err := protocol.MarshalToMessage(protocol.JoinRequest, joinData)
	if err != nil {
		return nil, err
	}

	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		return nil, err
	}

	incMsgs := make(chan protocol.Message, 20)
	outMsgs := make(chan protocol.Message, 20)
	c := &Connection{
		conn,
		incMsgs,
		outMsgs,
	}

	go c.readLoop()
	go c.writeLoop()

	return c, nil
}

// use syscall/js to get url info and roomID
func getClientInfo() (string, string, string) {
	loc := js.Global().Get("location")

	prot := "ws://"
	if loc.Get("protocol").String() == "https:" {
		prot = "wss://"
	}
	host := loc.Get("host").String()

	params := js.Global().Get("URLSearchParams").New(loc.Get("search"))
	roomID := params.Call("get", "room").String()

	return prot, host, roomID
}
