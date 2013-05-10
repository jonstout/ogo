package main

import (
	"fmt"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/openflow/ofp10"
)

// This is a basic learning switch implementation
type DemoApplication struct {
	packetIn chan ofp10.OfpMsg
	hostMap map[string]map[string]uint16
}

func (b *DemoApplication) InitApplication(args map[string]string) {
	// SubscribeTo returns a chan to receive a specific message type.
	b.packetIn = ogo.SubscribeTo(ofp10.OFPT_PACKET_IN)
	// A place to store the source ports of MAC Addresses
	b.hostMap = make(map[string]map[string]uint16)
}

func (b *DemoApplication) Name() string {
	// Every application needs a name
	return "Demo"
}

func (b *DemoApplication) Receive() {
	for {
		select {
		case m := <-b.packetIn:
			if pkt, ok := m.Data.(*ofp10.OfpPacketIn); ok {
				// This could be launched in a separate goroutine,
				// but maps in Go aren't thread safe.
				if _, ok := b.hostMap[m.DPID]; !ok {
					b.hostMap[m.DPID] = make(map[string]uint16)
				}
				b.parsePacketIn(m.DPID, pkt)
			}
		}
	}
}

func (b *DemoApplication) parsePacketIn(dpid string, pkt *ofp10.OfpPacketIn) {
	eth := pkt.Data
	hwSrc := eth.HWSrc.String()
	hwDst := eth.HWDst.String()
	if _, ok := b.hostMap[dpid][hwSrc]; !ok {
		b.hostMap[dpid][hwSrc] = pkt.InPort
	}
	if _, ok := b.hostMap[dpid][hwDst]; ok {
		f1 := ofp10.NewFlowMod()
		act1 := ofp10.NewActionOutput()
		act1.Port = b.hostMap[dpid][hwDst]
		f1.Actions = append(f1.Actions, act1)
		m1 := ofp10.NewMatch()
		m1.DLSrc = eth.HWSrc
		m1.DLDst = eth.HWDst
		f1.Match = *m1
		f1.IdleTimeout = 3
		
		f2 := ofp10.NewFlowMod()
		act2 := ofp10.NewActionOutput()
		act2.Port = b.hostMap[dpid][hwSrc]
		f2.Actions = append(f1.Actions, act2)
		m2 := ofp10.NewMatch()
		m2.DLSrc = eth.HWDst
		m2.DLDst = eth.HWSrc
		f2.Match = *m2
		f2.IdleTimeout = 3
		if s, ok := ogo.GetSwitch(dpid); ok {
			s.Send(f1)
			s.Send(f2)
		}
	} else {
		// Flood
		pktOut := ofp10.NewPacketOut()
		pktOut.Actions = append(pktOut.Actions, ofp10.NewActionOutput())
		pktOut.Data = &pkt.Data
		if s, ok := ogo.GetSwitch(dpid); ok {
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
