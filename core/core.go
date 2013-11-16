package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"github.com/jonstout/ogo/pacit"
	"log"
	"net"
	"time"
)

// OgoInstance generator.
func NewInstance() interface{} {
	return new(OgoInstance)
}

type OgoInstance struct {
	shutdown chan bool
}

func (o *OgoInstance) ConnectionUp(dpid net.HardwareAddr) {
	arpFmod := ofp10.NewFlowMod()
	arpFmod.HardTimeout = 0
	arpFmod.Priority = 2
	arpFmod.Match.DLType = 0x0806 // ARP Messages
	arpFmod.AddAction(ofp10.NewActionOutput(ofp10.P_CONTROLLER))

	dscFmod := ofp10.NewFlowMod()
	dscFmod.HardTimeout = 0
	dscFmod.Priority = 3
	dscFmod.Match.DLType = 0xa0f1 // Link Discovery Messages
	dscFmod.AddAction(ofp10.NewActionOutput(ofp10.P_CONTROLLER))

	if sw, ok := Switch(dpid); ok {
		sw.Send(ofp10.NewFeaturesRequest())
		sw.Send(arpFmod)
		sw.Send(dscFmod)
		sw.Send(ofp10.NewEchoRequest())
	}
	go o.linkDiscoveryLoop(dpid)
}

func (o *OgoInstance) ConnectionDown(dpid net.HardwareAddr) {
	o.shutdown <- true
	log.Println("Switch Disconnected:", dpid)
}

func (o *OgoInstance) EchoRequest(dpid net.HardwareAddr) {
	// Wait three seconds then send an echo_reply message.
	go func() {
		<- time.After(time.Second * 3)
		if sw, ok := Switch(dpid); ok {
			res := ofp10.NewEchoReply()
			sw.Send(res)
		}
	}()
}

func (o *OgoInstance) EchoReply(dpid net.HardwareAddr) {
	// Wait three seconds then send an echo_reply message.
	go func() {
		<- time.After(time.Second * 3)
		if sw, ok := Switch(dpid); ok {
			res := ofp10.NewEchoRequest()
			sw.Send(res)
		}
	}()
}

func (o *OgoInstance) FeaturesReply(dpid net.HardwareAddr, features *ofp10.SwitchFeatures) {
	if sw, ok := Switch(dpid); ok {
		for _, p := range features.Ports {
			sw.SetPort(p.PortNo, p)
		}
	}
}

func (o *OgoInstance) PacketIn(dpid net.HardwareAddr, msg *ofp10.PacketIn) {
	eth := msg.Data
	if buf, ok := eth.Data.(*pacit.PacitBuffer); ok {
		linkMsg := new(LinkDiscovery)
		linkMsg.Write(buf.Bytes())

		latency := time.Since(time.Unix(0, linkMsg.Nsec))
		l := &Link{linkMsg.SrcDPID, msg.InPort, latency, -1}

		if sw, ok := Switch(dpid); ok {
			sw.setLink(dpid, l)
			//log.Println(sw.Links())
		}
	}
}

func (o *OgoInstance) linkDiscoveryLoop(dpid net.HardwareAddr) {
	for {
		select {
		case <- o.shutdown:
			return
		// Every two seconds send a link discovery packet.
		case <-time.After(time.Second * 2):
			eth := pacit.NewEthernet()
			eth.Ethertype = 0xa0f1
			eth.HWSrc = dpid[2:]
			eth.Data = NewLinkDiscovery(dpid)

			pkt := ofp10.NewPacketOut()
			pkt.Data = eth
			pkt.AddAction(ofp10.NewActionOutput(ofp10.P_FLOOD))
			
			/*if sw, ok := Switch(dpid); ok {
				sw.Send(pkt)
			}*/
		}
	}
}
