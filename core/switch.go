package core

import (
	//"errors"
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
	"time"
	"sync"
)

// A map from DPIDs to all Switches that have connected since
// Ogo started.
type Network struct {
	sync.RWMutex
	Switches map[string]*Switch
}

func NewNetwork() *Network {
	n := new(Network)
	n.Switches = make(map[string]*Switch)
	return n
}

var network *Network

type Switch struct {
	stream *MessageStream
	dpid          net.HardwareAddr
	ports         map[int]*ofp10.PhyPort
	portsMu sync.RWMutex
	links         map[string]*Link
	linksMu sync.RWMutex
	requests      map[uint32]chan ofp10.Msg
}

// Builds and populates a Switch struct then starts listening
// for OpenFlow messages on conn.
func NewSwitch(stream *MessageStream, msg *ofp10.FeaturesReply) {

	network.Lock()
	if sw, ok := network.Switches[msg.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID())
		sw.stream = stream
		go sw.sendSync()
		go sw.receive()
	} else {
		log.Println("Openflow Connection:", msg.DPID)
		s := new(Switch)
		s.stream = stream
		s.dpid = msg.DPID
		s.ports = make(map[int]*ofp10.PhyPort)
		s.links = make(map[string]*Link)
		s.requests = make(map[uint32]chan ofp10.Msg)
		for _, p := range msg.Ports {
			s.ports[int(p.PortNo)] = &p
		}
		network.Switches[msg.DPID.String()] = s
		go s.sendSync()
		go s.receive()
	}
	network.Unlock()
}

// Returns a pointer to the Switch mapped to dpid.
func Switch(dpid net.HardwareAddr) (*OFPSwitch, bool) {
	network.RLock()
	defer network.RUnlock()
	if sw, ok := network.Switches[dpid.String()]; ok {
		return sw, ok
	} else {
		return nil, false
	}
}

// Returns a slice of *OFPSwitches for operations across all
// switches.
func Switches() []*OFPSwitch {
	network.RLock()
	defer network.RUnlock()
	a := make([]*OFPSwitch, len(network.Switches))
	i := 0
	for _, v := range network.Switches {
		a[i] = v
		i++
	}
	return a
}

// Disconnects Switch dpid.
func disconnect(dpid net.HardwareAddr) {
	network.Lock()
	defer network.Unlock()
	log.Printf("Closing connection with: %s", dpid)
	network.Switches[dpid.String()].conn.Close()
	delete(network.Switches, dpid.String())
}

// Returns a slice of all links connected to Switch s.
func (s *OFPSwitch) Links() []Link {
	s.linksMu.RLock()
	a := make([]Link, len(s.links))
	i := 0
	for _, v := range s.links {
		a[i] = *v
		i++
	}
	s.linksMu.RUnlock()
	return a
}

// Returns the link between Switch s and the Switch dpid.
func (s *OFPSwitch) Link(dpid net.HardwareAddr) (l Link, ok bool) {
	s.linksMu.RLock()
	if n, k := s.links[dpid.String()]; k {
		l = *n
		ok = true
	}
	s.linksMu.RUnlock()
	return
}

// Updates the link between s.DPID and l.DPID.
func (s *OFPSwitch) setLink(dpid net.HardwareAddr, l *Link) {
	s.linksMu.Lock()
	s.links[l.DPID.String()] = l
	s.linksMu.Unlock()
}

// Returns the dpid of Switch s.
func (s *OFPSwitch) DPID() net.HardwareAddr {
	return s.dpid
}

// Returns a slice of all the ports from Switch s.
func (s *OFPSwitch) Ports() []ofp10.PhyPort {
	s.portsMu.RLock()
	a := make([]ofp10.PhyPort, len(s.ports))
	i := 0
	for _, v := range s.ports {
		a[i] = *v
		i++
	}
	s.portsMu.RUnlock()
	return a
}

// Returns a pointer to the OfpPhyPort at port number from Switch s.
func (s *OFPSwitch) Port(number int) (q ofp10.PhyPort, ok bool) {
	s.portsMu.RLock()
	if p, k := s.ports[number]; k {
		q = *p
		ok = true
	}
	s.portsMu.RUnlock()
	return
}

// Sends an OpenFlow message to this Switch.
func (s *OFPSwitch) Send(req ofp10.Packet) {
	s.stream.Outbound <- req
}

// Receive loop for each Switch.
func (s *OFPSwitch) receive() {
	for {
		select {
		case msg <- s.stream.Inbound:
			s.distributeMessages(s.dpid, msg)
		case err <- s.stream.Error:
			return
		}
	}
}

func (s *OFPSwitch) distributeReceived(dpid net.HardwareAddr, msg ofp10.Msg) {
	h := msg.Data.GetHeader()
	if pktChan, ok := s.requests[h.XID]; ok {
		select {
		case pktChan <- p:
		case <-time.After(time.Millisecond * 100):
		}
		delete(s.requests, h.XID)
	} else {
		for _, ch := range messageChans[h.Type] {
			select {
			case ch <- p:
			case <-time.After(time.Millisecond * 100):
			}
		}
	}
}

// Sends an OpenFlow message to s, and returns a channel to receive
// a response on. Any error encountered during the send except io.EOF
// is returned.
func (s *OFPSwitch) SendAndReceive(req ofp10.Packet) (p chan ofp10.Msg, err error) {
	p = make(chan ofp10.Msg)
	s.requests[req.GetHeader().XID] = p
	err = s.Send(req)
	if err != nil {
		delete(s.requests, req.GetHeader().XID)
		return nil, err
	}
	return
}
