package core

import (
	"log"
	"time"
	"github.com/jonstout/pacit"
	"github.com/jonstout/ogo/openflow/ofp10"
)


type Core struct {
	echoRequest chan ofp10.Msg
	portStatus chan ofp10.Msg
	packetIn chan ofp10.Msg
}


func (b *Core) InitApplication(args map[string]string) {
	b.echoRequest = SubscribeTo(ofp10.T_ECHO_REQUEST)
	b.portStatus = SubscribeTo(ofp10.T_PORT_STATUS)
	b.packetIn = SubscribeTo(ofp10.T_PACKET_IN)
}


func (b *Core) Name() string {
	return "Core"
}


func (b *Core) Receive() {
	for {
		select {
		case m := <-b.echoRequest:
			go b.SendEchoReply(m.DPID)
		case m := <-b.portStatus:
			go b.UpdatePortStatus(m)
		case m := <-b.packetIn:
			if pkt, ok := m.Data.(*ofp10.PacketIn); ok {
				b.handlePacketIn(m.DPID, pkt)
			}
		case <- time.After(time.Second * 1):
			go b.discoverLinks()
		}
	}
}


func (b *Core) discoverLinks() {
	for _, sw := range Switches() {
		pkt := ofp10.NewPacketOut()
		act := ofp10.NewActionOutput(ofp10.P_ALL)
		pkt.AddAction(act)
		if data, err := NewListDiscovery(sw.DPID().String()); err != nil {
			log.Println(err)
		} else {
			pkt.Data = data
			sw.Send(pkt)
		}
	}
}


func (b *Core) handlePacketIn(dpid string, msg *ofp10.PacketIn) {
	eth := msg.Data
	log.Println("Handling packet ins!")
	if buf, ok := eth.Data.(*pacit.PacitBuffer); ok {
		log.Println(buf.String())
	}
}


func (b *Core) SendEchoReply(dpid string) {
	if s, ok := Switch(dpid); ok {
		<-time.After(time.Second * 3)
		res := ofp10.NewEchoReply()
		s.Send(res)
	}
}


func (b *Core) UpdatePortStatus(msg ofp10.Msg) {
	if _, ok := Switch(msg.DPID); ok {		
	}
}
