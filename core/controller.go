package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
)

type Controller struct { }

func NewController() *Controller {
	c := &new(Controller)
	Applications = make(map[string]*Application)
	network = NewNetwork()

	coreApplication := &new(Core)
	c.RegisterApplication(coreApplication)
	return c
}


func (c *Controller) Listen(port string) {
	addr, _ := net.ResolveTCPAddr("tcp", port)
	if sock, err := net.ListenTCP("tcp", addr); err != nil {
		if conn, e := sock.AcceptTCP(); e != nil {
			log.Println(e)
		} else {
			go c.handleConnection(conn)
		}
	} else {
		log.Fatal(err)
	}
}


func (c *Controller) handleConnection(conn *net.TCPConn) {
	stream := NewMessageStream(conn)

	for {
		select {
		case stream.Outbound <- ofp.NewHello():
			// Send hello message with latest protocol version.
		case msg <- stream.Inbound:
			switch m := msg.(type) {
			// A Hello message of the appropriate type
			// completes version negotiation. If version
			// types are incompatable, it is possible the
			// connection may be servered without error.
			case *ofp10.Hello:
				if msg.Version == ofp10.Version {
					// Version negotiation is
					// considered complete. Create
					// new Switch and notifiy listening
					// applications.
					stream.Version = msg.Version
					NewSwitch(stream, msg)

					for a := range Applications {
						if app, ok := a.(ofp10.ConnectionUpReactor); ok {
							app.ConnectionUp(msg.DPID)
						}
					}
					return
				} else {
					// Connection should be severed if controller
					// doesn't support switch version.
					stream.Shutdown <- true
				}
			// An error message may indicate a version mismatch. We
			// can attempt to continue with a vaild
			// FeaturesRequest.
			case *ofp10.ErrorMsg:
				stream.Version = msg.Version
				stream.Outbound <- ofp10.NewFeaturesRequest()
			}
		case err <- stream.Error:
			// The connection has been shutdown.
			log.Println(err)
			return
		case <- time.After(time.Second * 3):
			// This shouldn't happen. If it does, both the controller
			// and switch are no longer communicating. The TCPConn is
			// still established though.
			log.Println("Connection timed out.")
			return
		}
	}
}


// Setup OpenFlow Message chans for each message type.
func (c *Controller) RegisterApplication(app *Application) {
	app.Initialize(make(map[string]string))
	go app.Receive()
	Applications[app.Name()] = app
}
