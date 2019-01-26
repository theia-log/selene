package comm

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type theiaData struct {
	data []byte
	err  error
}

type theiaConn struct {
	conn *websocket.Conn
}

func (t *theiaConn) Open(url string) error {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	t.conn = c
	return nil
}

func (t *theiaConn) Send(data []byte) error {
	return t.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (t *theiaConn) Read() chan *theiaData {
	dataChan := make(chan *theiaData)
	go func() {
		for {
			messageType, data, err := t.conn.ReadMessage()
			if err != nil {
				dataChan <- &theiaData{
					err: err,
				}
				close(dataChan)
				return
			}
			switch messageType {
			case websocket.CloseMessage:
				dataChan <- &theiaData{
					err: fmt.Errorf(string(data)),
				}
				close(dataChan)
			case websocket.BinaryMessage:
				dataChan <- &theiaData{
					data: data,
				}
			default:
				log.Println("Ignoring message of type: ", messageType)
			}
		}
	}()
	return dataChan
}

func (t *theiaConn) Close(reason string) error {
	return t.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason))
}

type WebsocketClient struct {
	baseURL     string
	connections map[string]*theiaConn
}
