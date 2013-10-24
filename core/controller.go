package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
	"time"
)

type Controller struct { }

type InstanceGen func() interface{}

func NewController() *Controller {
	c := new(Controller)
	Applications = *new([]InstanceGen)
	network = NewNetwork()

	if f, ok := NewInstance().(InstanceGen); ok {
		c.RegisterApplication(f)
	}
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
		case stream.Outbound <- ofp10.NewHello():
			// Send hello message with latest protocol version.
		case msg := <- stream.Inbound:
			switch m := msg.(type) {
			// A Hello message of the appropriate type
			// completes version negotiation. If version
			// types are incompatable, it is possible the
			// connection may be servered without error.
			case *ofp10.Header:
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
				stream.Version = m.Header.Version
				stream.Shutdown <- true
			}
		case err := <- stream.Error:
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
func (c *Controller) RegisterApplication(fn InstanceGen) {
	Applications = append(Applications, fn)
}
