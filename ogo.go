package main

import (
	"net"
	"fmt"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/openflow/ofp10"
)

// This is a basic learning switch implementation
type DemoApplication struct {
	packetIn chan ofp10.Msg
	hostMap map[string]uint16
}

func (b *DemoApplication) InitApplication(args map[string]string) {
	// SubscribeTo returns a chan to receive a specific message type.
	b.packetIn = core.SubscribeTo(ofp10.T_PACKET_IN)
	b.hostMap = make(map[string]uint16)
}

func (b *DemoApplication) Name() string {
	// Every application needs a name
	return "Demo"
}

func (b *DemoApplication) Receive() {
	for {
		select {
		case m := <-b.packetIn:
			if pkt, ok := m.Data.(*ofp10.PacketIn); ok {
				if pkt.Data.Ethertype == 0x806 {
					b.parsePacketIn(m.DPID, pkt)
				}
			}
		}
	}
}

func (b *DemoApplication) parsePacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	eth := pkt.Data
	hwSrc := eth.HWSrc.String()
	hwDst := eth.HWDst.String()
	if _, ok := b.hostMap[hwSrc]; !ok {
		fmt.Println("Learning host", hwSrc)
		b.hostMap[hwSrc] = pkt.InPort
	}
	if _, ok := b.hostMap[hwDst]; ok {
		f1 := ofp10.NewFlowMod()
		f1.AddAction( ofp10.NewActionOutput(b.hostMap[hwDst]) )
		f1.Match.DLSrc = eth.HWSrc
		f1.Match.DLDst = eth.HWDst
		f1.IdleTimeout = 3
		
		f2 := ofp10.NewFlowMod()
		f2.AddAction( ofp10.NewActionOutput(b.hostMap[hwSrc]) )
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
	ctrl.Start(":6633")
}
