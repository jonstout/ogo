package main

import (
	"fmt"
	"github.com/jonstout/ogo/core"
	"github.com/jonstout/ogo/protocol/ofp10"
	"net"
	"sync"
	"runtime"
)

// Structure to track hosts that we discover.
type Host struct {
	mac  net.HardwareAddr
	port uint16
}

// A thread safe map to store our hosts.
type HostMap struct {
	hosts map[string]Host
	sync.RWMutex
}

func NewHostMap() *HostMap {
	h := new(HostMap)
	h.hosts = make(map[string]Host)
	return h
}

// Returns the host associated with mac.
func (m *HostMap) Host(mac net.HardwareAddr) (h Host, ok bool) {
	m.RLock()
	defer m.RUnlock()
	h, ok = m.hosts[mac.String()]
	return
}

// Records the host mac address and the port where mac was discovered.
func (m *HostMap) SetHost(mac net.HardwareAddr, port uint16) {
	m.Lock()
	defer m.Unlock()
	m.hosts[mac.String()] = Host{mac, port}
}

var hostMap HostMap

// Returns a new instance that implements one of the many
// interfaces found in ofp/ofp10/interface.go. One
// DemoInstance will be created for each switch that connects
// to the network.
func NewDemoInstance() interface{} {
	return &DemoInstance{&hostMap}
}

// Acts as a simple learning switch.
type DemoInstance struct {
	*HostMap
}

func (b *DemoInstance) PacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
	eth := pkt.Data
	// Ignore link discovery packet types.
	if eth.Ethertype == 0xa0f1 || eth.Ethertype == 0x88cc {
		return
	}

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
		p.InPort = pkt.InPort
		p.AddAction(ofp10.NewActionOutput(ofp10.P_ALL))
		p.Data = &eth
		if sw, ok := core.Switch(dpid); ok {
			sw.Send(p)
		}
	}
}

func main() {
	fmt.Println("Ogo 2013")
	runtime.GOMAXPROCS(runtime.NumCPU())
	ctrl := core.NewController()
	hostMap = *NewHostMap()
	ctrl.RegisterApplication(NewDemoInstance)
	ctrl.Listen(":6633")
}
