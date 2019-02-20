// Package comm contains Client API and communication layer
// for theia server.
//
// Communication to theia server is established via websocket
// channel. The default implementation of the Client library
// that communicates to theia server is based on websockets
// implementation.
//
// The Client API implements the basic functionalities that theia
// offers:
// 	- send (publish) an Event to theia server
//  - find past events, and
//  - receive events from the server is real time.
//
// Here is an example of establishing connection to theia server
// and publishing an event:
//	import (
//		"github.com/theia-log/selene/comm"
//		"github.com/theia-log/selene/model"
//		uuid "github.com/satori/go.uuid"
//	)
//
//	func main() {
//		client := comm.NewWebsocketClient("ws://localhost:6433")
//		if err := client.Send(&model.Event{
//			ID: 		uuid.Must(uuid.NewV4()).String(),
//			Timestamp: 	1550695140.89999200,
//			Source: 	"/dev/sensors/temp-sensor",
//			Tags: 		[]string{ "sensor", "home", "temp" },
//			Content: 	"10C",
//		}); err != nil {
//			panic(err)
//		}
//	}
package comm
