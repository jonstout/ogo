package main

import (
	"fmt"
	"github.com/jonstout/pacit"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/openflow/ofp10"
)


type DemoApplication struct {
	ogo.OgoApplication
	packetIn chan ofp10.OfpMsg
}

func (b *DemoApplication) InitApplication(args map[string]string) {
	b.Name = "Demo"
	b.packetIn = ogo.SubscribeTo(ofp10.OFPT_PACKET_IN)
}

func (b *DemoApplication) Receive() {
	for {
		select {
		case m := <-b.packetIn:
			if pkt, ok := m.Data.(*ofp10.OfpPacketIn); ok {
				b.parsePacketIn(m.DPID, pkt)
			}
		}
	}
}

func (b *DemoApplication) parsePacketIn(dpid string, pkt *ofp10.OfpPacketIn) {
	if s, ok := ogo.GetSwitch(dpid); ok {
		pktOut := ofp10.NewPacketOut()
		pktOut.Actions = append(pktOut.Actions, *ofp10.NewActionOutput())
		if arp, ok := pkt.Data.(*pacit.ARP); ok {
			pktOut.Data = arp
			s.Send(pktOut)
		}
	}
}

func main() {
	fmt.Println("Ogo 2013")
	ctrl := ogo.NewController()
	ctrl.RegisterApplication(new(DemoApplication))
	ctrl.Start(":6633")
}
