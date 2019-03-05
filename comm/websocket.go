package comm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/theia-log/selene/model"
)

// theiaData represents a data packet or error returned from the raw websocket
// channel to theia server.
type theiaData struct {
	data []byte
	err  error
}

// IsACK checks if this data packet is server ACK.
func (d *theiaData) IsACK() bool {
	if d.data == nil {
		return false
	}
	return string(d.data) == "ok"
}

// GetServerError checks if the data packet from the server is actually an error
// instead of event data.
// If so, it returns an error. Otherwise returns nil.
func (d *theiaData) GetServerError() error {
	if d.data == nil {
		return nil
	}
	str := string(d.data)
	if strings.HasPrefix(str, "{\"error\"") {
		errMap := map[string]string{}
		if err := json.Unmarshal(d.data, &errMap); err != nil {
			// it is not a server error
			return nil
		}
		return fmt.Errorf(errMap["error"])
	}
	return nil
}

// theiaConn represents an open websocket connection to Theia sever.
// This connection can be reused.
type theiaConn struct {
	url  string
	conn *websocket.Conn
}

// Open connects and opens the actual connection to Theia.
func (t *theiaConn) Open() error {
	c, _, err := websocket.DefaultDialer.Dial(t.url, nil)
	if err != nil {
		return err
	}
	t.conn = c
	return nil
}

// Send sends raw data to theia server.
func (t *theiaConn) Send(data []byte) error {
	return t.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Read consumes data (websocket messages) from the open websocket channel to
//  theia server.
// Once the data is received, it is wrapped in theiaData packet.
// These packets are then published on a theiaData channel for further
// processing.
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
			case websocket.TextMessage:
				dataChan <- &theiaData{
					data: []byte(data),
				}
			default:
				log.Println("Ignoring message of type: ", messageType)
			}
		}
	}()
	return dataChan
}

// Close closes the underlying websocket connection with the given reason.
// A formal close message is issued to the server before breaking up
// the connection.
func (t *theiaConn) Close(reason string) error {
	return t.conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason))
}

// newConn creates new raw theia connection to a given server and for a
// particular action.
func newConn(baseURL, action string) *theiaConn {
	return &theiaConn{
		url: fmt.Sprintf("%s/%s", baseURL, action),
	}
}

// WebsocketClient implements the Client interface.
// Implements a client to a particular Theia server.
// Connections to the theia actions, like /event, /find and /live are reused if
// possible - new connections will not be opened if a channel is already
// established on the endpoint.
type WebsocketClient struct {
	baseURL     string
	connections map[string]*theiaConn
}

// getConn returns an exiting connection for a particular endpoint.
// The connections are reused, new connection is created only if a connection to
// the requested endpoint does not exist.
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

// Send send an event to the server.
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

// doReceive sends EventFilter data to the endpoint on the server, then listens
// for events from the server.
// The returned theiaData packets are decoded to EventResponse structure and
// published on the EventResponse channel.
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
			if data.IsACK() {
				// server ACK. We san safely ignore this message.
				continue
			}
			if err = data.GetServerError(); err != nil {
				// Server responded with an error.
				eventChan <- &EventResponse{
					Error: err,
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

// Receive opens a channel for real-time events that match the EventFilter.
func (w *WebsocketClient) Receive(filter *EventFilter) (chan *EventResponse, error) {
	return w.doReceive("live", filter)
}

// Find looks up past events that match the given EventFilter.
func (w WebsocketClient) Find(filter *EventFilter) (chan *EventResponse, error) {
	return w.doReceive("find", filter)
}

// NewWebsocketClient creates new websocket Client to theia server on the given
// server URL.
func NewWebsocketClient(serverURL string) *WebsocketClient {
	return &WebsocketClient{
		baseURL: serverURL,
		connections: map[string]*theiaConn{
			"event": newConn(serverURL, "event"),
			"live":  newConn(serverURL, "live"),
			"find":  newConn(serverURL, "find"),
		},
	}
}
