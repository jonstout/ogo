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
	c.Shutdown = shutdown
	go c.loop()
}

func (c *Ogo) NewInstance() interface{} {
	return &OgoInstance{c}
}

type OgoInstance struct {
	*Ogo
}

func (c *Core) ConnectionUp(dpid net.HardwareAddr) {
	log.Println("Switch Connected:", dpid)

	if sw, ok := Switch(dpid); ok {
		sw.Send(ofp10.NewFeaturesRequest())
	}
}

func (c *Core) ConnectionDown(dpid net.HardwareAddr) {
	log.Println("Switch Disconnected:", dpid)
}

func (c *Core) FeaturesReply(dpid net.HardwareAddr, features ofp10.SwitchFeatures) {

}

func (c *Core) EchoRequest(dpid net.HardwareAddr, req ofp10.Header) {
	<- time.After(time.Second * 3)
	if sw, ok := Switch(dpid); ok {
		res := ofp10.NewEchoReply()
		sw.Send(res)
	}
}

func (c *Core) PacketIn(dpid net.HardwareAddr, msg ofp10.PacketIn) {
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

func (c *Core) loop() {
	for {
		select {
		case <- c.Shutdown:
			return
		case <-time.After(time.Second * 1):
			c.discoverLinks()
		}
	}
}

func (c *Core) discoverLinks() {
	for _, sw := range Switches() {
		pkt := ofp10.NewPacketOut()
		pkt.AddAction(ofp10.NewActionOutput(ofp10.P_FLOOD))

		eth := pacit.NewEthernet()
		eth.Ethertype = 0xa0f1
		eth.HWSrc = sw.DPID()[2:]
		eth.Data = NewLinkDiscovery(sw.DPID())
		pkt.Data = eth
		sw.Send(pkt)
	}
}
