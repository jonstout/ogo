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
	hostMap map[string]int
}

func (b *DemoApplication) InitApplication(args map[string]string) {
	b.Name = "Demo"
	b.packetIn = ogo.SubscribeTo(ofp10.OFPT_PACKET_IN)
	b.hostMap = make(map[string]int)
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
		if eth, ok := pkt.Data.(*pacit.Ethernet); ok {
			if src, ok := b.hostMap[eth.HWSrc.String()]; ok {
				if dst, ok := b.hostMap[eth.HWDst.String()]; ok {
					// Match := src,dst Fwd := dst.port
					// Match := dst,src Fwd := src.port
					fmt.Println("Src:", src, "and Dst:", dst, "known.")
				} else {
					// Flood
					pktOut := ofp10.NewPacketOut()
					pktOut.Actions = append(pktOut.Actions, ofp10.NewActionOutput())
					pktOut.Data = pkt.Data
					s.Send(pktOut)
				}
			} else {
				// Add to map
				b.hostMap[eth.HWSrc.String()] = int(pkt.InPort)
			}
		}
	}
}

func main() {
	fmt.Println("Ogo 2013")
	ctrl := ogo.NewController()
	ctrl.RegisterApplication(new(DemoApplication))
	ctrl.Start(":6633")
}
