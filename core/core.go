package core

import (
	"github.com/jonstout/ogo/openflow/ofp10"
	"github.com/jonstout/ogo/pacit"
	"log"
	"net"
	"time"
	"sync"
)

type Ogo struct {
	switchesMutex sync.RWMutex
}

func NewOgo() *Ogo {
	ogo := new(Ogo)
	return ogo
}

func (c *Ogo) NewInstance() interface{} {
	return &OgoInstance{c}
}

type OgoInstance struct {
	*Ogo
	shutdown chan bool
}

func (o *OgoInstance) ConnectionUp(dpid net.HardwareAddr) {
	log.Println("Switch Connected:", dpid)

	if sw, ok := Switch(dpid); ok {
		sw.Send(ofp10.NewFeaturesRequest())
	}
	go linkDiscoveryLoop(dpid)
}

func (c *OgoInstance) ConnectionDown(dpid net.HardwareAddr) {
	o.shutdown <- true
	log.Println("Switch Disconnected:", dpid)
}

func (c *OgoInstance) EchoRequest(dpid net.HardwareAddr) {
	// Wait three seconds then send an echo_reply message.
	<- time.After(time.Second * 3)
	if sw, ok := Switch(dpid); ok {
		res := ofp10.NewEchoReply()
		sw.Send(res)
	}
}

func (c *OgoInstance) FeaturesReply(dpid net.HardwareAddr, features *ofp10.SwitchFeatures) {
	if sw, ok := Switch(dpid); ok {
		for p := range features.Ports {
			sw.SetPort(features.PortNo, features)
		}
	}
}

func (c *OgoInstance) PacketIn(dpid net.HardwareAddr, msg *ofp10.PacketIn) {
	eth := msg.Data
	if buf, ok := eth.Data.(*pacit.PacitBuffer); ok {
		linkMsg := new(LinkDiscovery)
		linkMsg.Write(buf.Bytes())

		latency := time.Since(time.Unix(0, linkMsg.Nsec))
		l := &Link{linkMsg.SrcDPID, msg.InPort, latency, -1}

		if sw, ok := Switch(dpid); ok {
			sw.setLink(dpid, l)
		}
	}
}

func (o *OgoInstance) linkDiscoveryLoop(dpid net.HardwareAddr) {
	for {
		case <- o.shutdown:
			return
		// Every two seconds send a link discovery packet.
		case <-time.After(time.Second * 2):
			eth := pacit.NewEthernet()
			eth.Ethertype = 0xa0f1
			eth.HWSrc = sw.DPID()[2:]
			eth.Data = NewLinkDiscovery(dpid)

			pkt := ofp10.NewPacketOut()
			pkt.Data = eth
			pkt.AddAction(ofp10.NewActionOutput(ofp10.P_FLOOD))
			
			if sw, ok := Switch(dpid); ok {
				sw.Send(pkt)
			}
		}
	}
}
