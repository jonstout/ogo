package main

import (
	"fmt"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/openflow/ofp10"
	"net"
	"sync"
)

// Structure to track hosts that we discover.
type Host struct {
	mac net.HardwareAddr
	port uint16
}

// A thread safe map to store our hosts. We are unlikely to
// actually need a thread safe data structure in this demo.
type HostMap struct {
	hosts  map[string]Host
	sync.RWMutex
}

func NewHostMap() *HostMap {
	h := new(HostMap)
	h.hosts = make(map[string]Host)
	return h

}

func (m *HostMap) Host(mac net.HardwareAddr) (h Host, ok bool) {
	m.RLock()
	defer m.RUnlock()
	h, ok = m.hosts[mac.String()]
	return
}

func (m *HostMap) SetHost(mac net.HardwareAddr, port uint16) {
	m.Lock()
	defer m.Unlock()
	m.hosts[mac.String()] = Host{mac, port}
}

// Returns a new instance that implements one of the many
// interfaces found in ofp/ofp10/interface.go. One
// DemoInstance will be created for each switch that connects
// to the network. Of course you could return the same
// pointer every time as well.
func NewDemoInstance() interface{} {
	return &DemoInstance{NewHostMap()}
}

// The instance is passed a pointer to the application
// for global variables and its own unique HostMap. Each
// unique instance will act as its own learning switch.
type DemoInstance struct {
	*HostMap
}

func (b *DemoInstance) PacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	eth := pkt.Data
	if eth.Ethertype != 0x0806 {
		return
	}

	b.SetHost(eth.HWSrc, pkt.InPort)
	if host, ok := b.Host(eth.HWDst); ok {
		fmt.Println(eth.HWSrc, ":", pkt.InPort, "to", eth.HWDst, ":", host.port)
		f1 := ofp10.NewFlowMod()
		f1.Match.DLSrc = eth.HWSrc
		f1.Match.DLDst = eth.HWDst
		f1.AddAction(ofp10.NewActionOutput(host.port))
		f1.IdleTimeout = 3

		f2 := ofp10.NewFlowMod()
		f2.Match.DLSrc = eth.HWDst
		f2.Match.DLDst = eth.HWSrc
		f2.AddAction(ofp10.NewActionOutput(pkt.InPort))
		f2.IdleTimeout = 3

		if s, ok := core.Switch(dpid); ok {
			s.Send(f1)
			s.Send(f2)
		}
	} else {
		p := ofp10.NewPacketOut()
		a := ofp10.NewActionOutput(ofp10.P_ALL)
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

	ctrl.RegisterApplication(NewDemoInstance)

	ctrl.Listen(":6633")
}
