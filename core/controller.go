package core

import (
	"log"
	"net"
	"github.com/jonstout/ogo/openflow/ofp10"
)

type OgoController struct { }

func NewController() *OgoController {
	o := new(OgoController)
	Applications = make(map[string]Application)
	messageChans = make(map[uint8][]chan ofp10.Msg)
	Switches = make(map[string]*OFPSwitch)
	// Register ogo core
	b := new(Core)
	o.RegisterApplication(b)
	return o
}

func (o *OgoController) Start(port string) {
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
		go NewOpenFlowSwitch(conn)
	}
}

func (o *OgoController) RegisterApplication(app Application) {
	// Setup Openflow Message Channels
	app.InitApplication(make(map[string]string))
	go app.Receive()
	Applications[app.Name()] = app
}
