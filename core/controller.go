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
	switches = make(map[string]*OFPSwitch)
	// Register ogo core
	b := new(Core)
	o.RegisterApplication(b)
	return o
}

func (o *Controller) Start(port string) {
	// Listen for Switch Connections
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

func (o *Controller) RegisterApplication(app Application) {
	// Setup Openflow Message Channels
	app.InitApplication(make(map[string]string))
	go app.Receive()
	Applications[app.Name()] = app
}
