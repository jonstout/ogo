package core

type Controller struct {

}

func NewController() *Controller {
	c := new(Controller)
}

func (c *Controller) Listen(port string) {
	addr, _ := net.ResolveTCPAddr("tcp", port)
	if sock, err := net.ListenTCP("tcp", addr); err != nil {
		if conn, e := net.AcceptTCP(); e != nil {
			go c.handleConnection(conn)
		} else {
			log.Println(e)
		}
	} else {
		log.Fatal(err)
	}
}


func (c *Controller) handleConnection(conn *net.TCPConn) {
	stream := NewMessageStream(conn)
	stream.Version = 1
	record := new(Record)

	// Send a hello message
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
				for v := range ofp.SupportedVersions {
					if v == msg.Version {
						if record.Sent >= msg.Version {
							// Version negotiation is
							// considered complete. Ask for
							// switch features with received
							// version. Should record version
							// in case connection is severed.
							// stream.Version is lost in case
							// of a disconnect.
							stream.Version = msg.Version
							stream.Outbound <- ofp.NewFeaturesRequest()
						} else {
							// This should never happen, as
							// the expected behavior is to
							// start with the highest
							// supported protocol and adapt to
							// the switches version.
						}
						return
					}
				}
				// Connection should be severed if controller
				// doesn't support switch version.
				stream.Shutdown <- true
				return
			case *ofp10.ErrorMsg:
				// An error message may indicate a version mismatch. We
				// can attempt to continue with a vaild
				// FeaturesRequest.
				for v := range ofp.SupportedVersions {
					if v == msg.Version {
						stream.Version = msg.Version
						stream.Outbound <- ofp.NewFeaturesRequest()
					}
				}
				stream.Shutdown <- true
				return
			}
			return
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

type Record struct {

}

var unresolvedConnections []Record

// Maps openflow versions to ip addresses. Until a successful
// version negotiation has occoured.
func (c *Controller) nextHello(ip net.IP) pacit.Message {
}

// Let's the controller know about the last hello message's
// version number.l
func (c *Controller) recvHello(ip net.IP, ver uint16) {

}
