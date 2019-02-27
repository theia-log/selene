// Package comm contains Client API and communication layer for theia server.
//
// Communication to theia server is established via websocket channel. The
// default implementation of the Client library that communicates to theia
// server is based on websockets implementation.
//
// The Client API implements the basic functionalities that theia offers:
// send (publish) an Event to theia server;
// find past events;
// and receive events from the server in real time.
//
// Here is an example of establishing connection to theia server and publishing
// an event:
//	import (
//		"github.com/theia-log/selene/comm"
//		"github.com/theia-log/selene/model"
//	)
//
//	func main() {
//		// Create new client to the server
//		client := comm.NewWebsocketClient("ws://localhost:6433")
//
//		// Send the event
//		if err := client.Send(&model.Event{
//			ID: 		model.NewEventID(),
//			Timestamp: 	1550695140.89999200,
//			Source: 	"/dev/sensors/temp-sensor",
//			Tags: 		[]string{"sensor", "home", "temp"},
//			Content: 	"10C",
//		}); err != nil {
//			panic(err)
//		}
//	}
//
// Reading events from the server is done in an asynchronous way. The functions
// that read from the server return a read channel that will push new events as
// they come from the server.
//
// An example of looking up past event using EventFilter:
//
//	import (
//		"log"
//
//		"github.com/theia-log/selene/comm"
//		"github.com/theia-log/selene/model"
//	)
//
//	func main() {
//		// Create new client to the server
//		client := comm.NewWebsocketClient("ws://localhost:6433")
//
//		respChan, err := client.Find(&comm.EventFilter{
//			Start: 	1550710745.10, 	// return only events that happened after
//									// this time
//			End:	1550710746.90,	// but before this time
//			Tags: 	[]string{"sensor", "temp.+"},	// events that contain
//													// these tags
//			Content: "\\d+C",		// content that matches this regex
//			Order: 	"asc",			// ascending, by timestamp
//		})
//
//		if err != nil {
//			panic(err)
//		}
//
//		for {
//			resp, ok := <- respChan
//			if !ok {
//				break	// we're done, no more events
//			}
//			if resp.Error != nil {
//				// an error occurred, log it
//				log.Println("[ERROR]: ", resp.Error.Error())
//				continue
//			}
//			log.Println(resp.Event.Dump())	// print the event
//		}
//	}
//
// Tracking events in real-time looks exactly the same as looking up past events.
// To receive the events in real time, we use comm.Client.Receive() function:
//	import (
//		"log"
//
//		"github.com/theia-log/selene/comm"
//		"github.com/theia-log/selene/model"
//	)
//
//	func main() {
//		// Create new client to the server
//		client := comm.NewWebsocketClient("ws://localhost:6433")
//
//		// Receive reads events in real-time
//		respChan, err := client.Receive(&comm.EventFilter{
//			Start: 	1550710745.10, 	// return only events that happened after
//									// this time
//			Tags: 	[]string{"sensor", "temp.+"},	// events that contain
//													// these tags
//			Content: "\\d+C",		// content that matches this regex
//		})
//
//		if err != nil {
//			panic(err)
//		}
//
//		for {
//			resp, ok := <- respChan
//			if !ok {
//				break	// Server closed the connection, we're done.
//			}
//			if resp.Error != nil {
//				// an error occurred, log it
//				log.Println("[ERROR]: ", resp.Error.Error())
//				continue
//			}
//			log.Println(resp.Event.Dump())	// print the event
//		}
//	}
package comm
