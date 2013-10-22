package main

import (
	"fmt"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/openflow/ofp10"
	"net"
)

// This is a basic learning switch implementation
type DemoApplication struct {
	hosts  map[string]uint16
	newHost chan Host
}

type Host struct {
	mac net.HardwareAddr
	port uint16
}

func (b *DemoApplication) Initialize(args map[string]string, shutdown chan bool) {
	// SubscribeTo returns a chan to receive a specific message type.
	b.hosts = make(map[string]uint16)
	go b.loop()
}

func (b *DemoApplication) Name() string {
	// Every application needs a name.
	return "Demo"
}

func (b *DemoApplication) loop() {
	for {
		select {
		case h := <- b.newHost:
			if _, ok := b.hosts[h.mac.String()]; !ok {
				fmt.Println("Learning host", h.mac)
				b.hosts[h.mac.String()] = h.port
			}
		}
	}
}

func (b *DemoApplication) PacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	eth := pkt.Data
	hwSrc := eth.HWSrc.String()
	hwDst := eth.HWDst.String()

	b.newHost <- Host{eth.HWSrc, pkt.InPort}

	if _, ok := b.hosts[hwDst]; ok {
		f1 := ofp10.NewFlowMod()
		f1.AddAction(ofp10.NewActionOutput(b.hosts[hwDst]))
		f1.Match.DLSrc = eth.HWSrc
		f1.Match.DLDst = eth.HWDst
		f1.IdleTimeout = 3

		f2 := ofp10.NewFlowMod()
		f2.AddAction(ofp10.NewActionOutput(b.hosts[hwSrc]))
		f2.Match.DLSrc = eth.HWDst
		f2.Match.DLDst = eth.HWSrc
		f2.IdleTimeout = 3

		if s, ok := core.Switch(dpid); ok {
			s.Send(f1)
			s.Send(f2)
		}
	} else {
		p := ofp10.NewPacketOut()
		a := ofp10.NewActionOutput(ofp10.P_FLOOD)
		p.AddAction(a)
		p.Data = &eth
		if sw, ok := core.Switch(dpid); ok {
			sw.Send(p)
		}
	}
}

func main() {
	fmt.Println("Ogo 2013")
	ctrl := core.NewController()
	ctrl.RegisterApplication(new(DemoApplication))
	ctrl.Listen(":6633")
}
