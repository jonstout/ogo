package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
)

type Controller struct{}

func NewController() *Controller {
	o := new(Controller)
	Applications = make(map[string]Application)
	messageChans = make(map[uint8][]chan ofp10.Msg)
	network = NewNetwork()

	// Register ogo core
	b := new(Core)
	o.RegisterApplication(b)
	return o
}

// Start listening for Switch connections.
func (o *Controller) Start(port string) {
	addr, _ := net.ResolveTCPAddr("tcp", port)
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Println(err)
		}
		go NewOFPSwitch(conn)
	}
}

// Setup OpenFlow Message chans for each message type.
func (o *Controller) RegisterApplication(app Application) {
	app.InitApplication(make(map[string]string))
	go app.Receive()
	Applications[app.Name()] = app
}
