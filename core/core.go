package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"github.com/jonstout/pacit"
	"log"
	"net"
	"time"
)

type Core struct {
	echoRequest chan ofp10.Msg
	portStatus  chan ofp10.Msg
	packetIn    chan ofp10.Msg
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
		case <-time.After(time.Second * 2):
			b.discoverLinks()
		}
	}
}

func (b *Core) discoverLinks() {
	for _, sw := range Switches() {
		pkt := ofp10.NewPacketOut()

		act := ofp10.NewActionOutput(ofp10.P_FLOOD)
		pkt.AddAction(act)

		if data, err := NewListDiscovery(sw.DPID()); err != nil {
			log.Println(err)
		} else {
			eth := pacit.NewEthernet()
			eth.Ethertype = 0xa0f1
			eth.HWSrc = sw.DPID()[2:]
			eth.Data = data
			pkt.Data = eth
			sw.Send(pkt)
		}
	}
}

func (b *Core) handlePacketIn(dpid net.HardwareAddr, msg *ofp10.PacketIn) {
	eth := msg.Data
	if buf, ok := eth.Data.(*pacit.PacitBuffer); ok {
		lmsg := new(LinkDiscovery)
		lmsg.Write(buf.Bytes())

		latency := time.Since(time.Unix(0, lmsg.Nsec))
		l := &Link{lmsg.SrcDPID, msg.InPort, latency, -1}
		if sw, ok := Switch(dpid); ok {
			sw.setLink(dpid, l)
		}
	}
}

func (b *Core) SendEchoReply(dpid net.HardwareAddr) {
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
