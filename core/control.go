package core

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
					// considered complete. Ask for
					// switch features with received
					// version.
					stream.Version = msg.Version
					stream.Outbound <- ofp10.NewFeaturesRequest()
				} else {
					// Connection should be severed if controller
					// doesn't support switch version.
					stream.Shutdown <- true
				}
			// After a vaild FeaturesReply has been received we
			// have all the information we need. Create a new
			// switch object and notify applications.
			case *ofp10.FeaturesReply:
				registerSwitch(stream, msg)
				for a := range Applications {
					if app, ok := a.(ofp10.ConnectionUpReactor); ok {
						app.ConnectionUp(msg.DPID, *msg)
					}
				}
				return
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
	app.InitApplication(make(map[string]string))
	go app.Receive()
	Applications[app.Name()] = app
}
