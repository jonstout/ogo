package main

import (
	//"runtime"
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
	// A place to store the source ports of MAC Addresses
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
				// This could be launched in a separate goroutine,
				// but maps in Go aren't thread safe.
				b.parsePacketIn(m.DPID, pkt)
			}
		}
	}
}

func (b *DemoApplication) parsePacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	/*eth := pkt.Data
	hwSrc := eth.HWSrc.String()
	hwDst := eth.HWDst.String()
	if _, ok := b.hostMap[hwSrc]; !ok {
		b.hostMap[hwSrc] = pkt.InPort
	}
	if _, ok := b.hostMap[hwDst]; ok {
		f1 := ofp10.NewFlowMod()
		act1 := ofp10.NewActionOutput(b.hostMap[hwDst])
		f1.Actions = append(f1.Actions, act1)
		m1 := ofp10.NewMatch()
		m1.DLSrc = eth.HWSrc
		m1.DLDst = eth.HWDst
		f1.Match = *m1
		f1.IdleTimeout = 3
		
		f2 := ofp10.NewFlowMod()
		act2 := ofp10.NewActionOutput(b.hostMap[hwSrc])
		f2.Actions = append(f1.Actions, act2)
		m2 := ofp10.NewMatch()
		m2.DLSrc = eth.HWDst
		m2.DLDst = eth.HWSrc
		f2.Match = *m2
		f2.IdleTimeout = 3
		if s, ok := core.Switch(dpid); ok {
			s.Send(f1)
			s.Send(f2)
		}
	}*/
}

func main() {
	//runtime.GOMAXPROCS(16)
	fmt.Println("Ogo 2013")
	ctrl := core.NewController()
	//new(DemoApplication)
	//ctrl.RegisterApplication(new(DemoApplication))
	ctrl.Start(":6633")
}
