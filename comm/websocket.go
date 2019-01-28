package comm

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/theia-log/selene/model"
)

type theiaData struct {
	data []byte
	err  error
}

type theiaConn struct {
	url  string
	conn *websocket.Conn
}

func (t *theiaConn) Open() error {
	c, _, err := websocket.DefaultDialer.Dial(t.url, nil)
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

func NewConn(baseUrl, action string) *theiaConn {
	return &theiaConn{
		url: fmt.Sprintf("%s/%s", baseUrl, action),
	}
}

type WebsocketClient struct {
	baseURL     string
	connections map[string]*theiaConn
}

func (w *WebsocketClient) getConn(endpoint string) (*theiaConn, error) {
	conn, ok := w.connections[endpoint]
	if !ok {
		return nil, fmt.Errorf("no such endpoint")
	}
	if conn.conn == nil {
		if err := conn.Open(); err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (w *WebsocketClient) Send(event *model.Event) error {
	conn, err := w.getConn("event")
	if err != nil {
		return err
	}
	data, err := event.DumpBytes()
	if err != nil {
		return err
	}
	return conn.Send(data)
}

func (w *WebsocketClient) doReceive(endpoint string, filter *EventFilter) (chan *EventResponse, error) {
	conn, err := w.getConn(endpoint)
	if err != nil {
		return nil, err
	}

	filterData, err := filter.DumpBytes()
	if err != nil {
		return nil, err
	}

	if err = conn.Send(filterData); err != nil {
		return nil, err
	}

	dataChan := conn.Read()
	eventChan := make(chan *EventResponse)

	go func() {
		for {
			data, ok := <-dataChan
			if !ok {
				// channel closed
				close(eventChan)
				break
			}
			if data.err != nil {
				eventChan <- &EventResponse{
					Error: data.err,
				}
				continue
			}
			ev := &model.Event{}
			if err = ev.LoadBytes(data.data); err != nil {
				eventChan <- &EventResponse{
					Error: err,
				}
				continue
			}

			eventChan <- &EventResponse{
				Event: ev,
			}
		}
	}()

	return eventChan, nil
}

func (w *WebsocketClient) Receive(filter *EventFilter) (chan *EventResponse, error) {
	return w.doReceive("live", filter)
}

func (w WebsocketClient) Find(filter *EventFilter) (chan *EventResponse, error) {
	return w.doReceive("find", filter)
}

func NewWebsocketClient(serverURL string) *WebsocketClient {
	return &WebsocketClient{
		baseURL:     serverURL,
		connections: map[string]*theiaConn{},
	}
}
