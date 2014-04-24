package ogo

import (
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/ofp10"
	"log"
	"net"
	"time"
)

type Controller struct{}
type ApplicationInstanceGenerator func() interface{}

var Applications []ApplicationInstanceGenerator

func NewController() *Controller {
	c := new(Controller)
	Applications = *new([]ApplicationInstanceGenerator)
	network = NewNetwork()

	c.RegisterApplication(NewInstance)
	return c
}

func (c *Controller) Listen(port string) {
	addr, _ := net.ResolveTCPAddr("tcp", port)

	sock, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()

	log.Println("Listening for connections on", addr)
	for {
		conn, err := sock.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go c.handleConnection(conn)
	}
}

func (c *Controller) handleConnection(conn *net.TCPConn) {
	stream := NewMessageStream(conn)
	h, err := ofpxx.NewHello(1)
	if err != nil {
		return
	}
	stream.Outbound <- h

	for {
		select {
		//case stream.Outbound <- ofp10.NewHello():
		// Send hello message with latest protocol version.
		case msg := <-stream.Inbound:
			switch m := msg.(type) {
			// A Hello message of the appropriate type
			// completes version negotiation. If version
			// types are incompatable, it is possible the
			// connection may be servered without error.
			case *ofpxx.Header:
				if m.Version == ofp10.VERSION {
					// Version negotiation is
					// considered complete. Create
					// new Switch and notifiy listening
					// applications.
					stream.Version = m.Version
					stream.Outbound <- ofp10.NewFeaturesRequest()
				} else {
					// Connection should be severed if controller
					// doesn't support switch version.
					log.Println("Received unsupported ofp version", m.Version)
					stream.Shutdown <- true
				}
			// After a vaild FeaturesReply has been received we
			// have all the information we need. Create a new
			// switch object and notify applications.
			case *ofp10.SwitchFeatures:
				NewSwitch(stream, *m)
				for _, newInstance := range Applications {
					if sw, ok := Switch(m.DPID); ok {
						i := newInstance()
						sw.AddInstance(i)
					}
				}
				return
			// An error message may indicate a version mismatch. We
			// disconnect if an error occurs this early.
			case *ofp10.ErrorMsg:
				log.Println(m)
				stream.Version = m.Header.Version
				stream.Shutdown <- true
			}
		case err := <-stream.Error:
			// The connection has been shutdown.
			log.Println(err)
			return
		case <-time.After(time.Second * 3):
			// This shouldn't happen. If it does, both the controller
			// and switch are no longer communicating. The TCPConn is
			// still established though.
			log.Println("Connection timed out.")
			return
		}
	}
}

// Setup OpenFlow Message chans for each message type.
func (c *Controller) RegisterApplication(fn ApplicationInstanceGenerator) {
	Applications = append(Applications, fn)
}
