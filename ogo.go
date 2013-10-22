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
	hostsMutex sync.RWMutex
}

func NewHostMap() *HostMap {
	h := new(HostMap)
	h.hosts = make(map[string]Host)
	return h

}
func (hosts *HostMap) Host(mac net.HardwareAddr) (h Host, ok bool) {
	hosts.RLock()
	defer hosts.RUnlock()
	h, ok = dc.hosts[mac.String()]
	return
}

func (hosts *HostMap) SetHost(mac net.HardwareAddr, port uint16) {
	hosts.Lock()
	defer hosts.Unlock()
	dc.hosts[mac.String()] = Host{mac, port}
}

// Application to spawn per switch instances. May hold global
// variables such as connections to databases or channels to
// to web services.
type Demo struct {
}

func NewDemo() *DemoApplication {
	dc := new(DemoCore)
	return dc
}

// Returns a new instance that implements one of the many
// interfaces found in ofp/ofp10/interface.go
func (d *Demo) NewInstance() interface{} {
	// The instance is passed a pointer to the application
	// for global variables and its own unique HostMap. One
	// instance is spawned per OpenFlow Switch. Of course
	// you could return the same pointer every time as well.
	return &DemoInstance{d, NewHostMap()}
}

// The instance is passed a pointer to the application
// for global variables and its own unique HostMap. Each
// unique instance will act as its own learning switch.
type DemoInstance struct {
	*Demo
	*HostMap
}

func (b *DemoInstance) PacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	eth := pkt.Data
	b.SetHost(eth.HWSrc, pkt.InPort)

	if host, ok := b.Host(eth.HWDst); ok {
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
	ctrl.RegisterApplication(NewDemo())
	ctrl.Listen(":6633")
}
