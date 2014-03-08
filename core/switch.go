package core

import (
	"github.com/jonstout/ogo/protocol/ofp10"
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
	"log"
	"net"
	"sync"
)

// A map from DPIDs to all Switches that have connected since
// Ogo started.
type Network struct {
	sync.RWMutex
	Switches map[string]*OFSwitch
}

func NewNetwork() *Network {
	n := new(Network)
	n.Switches = make(map[string]*OFSwitch)
	return n
}

var network *Network

type OFSwitch struct {
	stream      *MessageStream
	appInstance []interface{}
	dpid        net.HardwareAddr
	ports       map[uint16]ofp10.PhyPort
	portsMu     sync.RWMutex
	links       map[string]*Link
	linksMu     sync.RWMutex
	reqs        map[uint32]chan util.Message
	reqsMu      sync.RWMutex
}

// Builds and populates a Switch struct then starts listening
// for OpenFlow messages on conn.
func NewSwitch(stream *MessageStream, msg ofp10.SwitchFeatures) {
	network.Lock()
	if sw, ok := network.Switches[msg.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID())
		sw.stream = stream
		go sw.receive()
	} else {
		log.Println("Openflow Connection:", msg.DPID)
		s := new(OFSwitch)
		s.stream = stream
		s.appInstance = *new([]interface{})
		s.dpid = msg.DPID
		s.ports = make(map[uint16]ofp10.PhyPort)
		s.links = make(map[string]*Link)
		s.reqs = make(map[uint32]chan util.Message)
		for _, p := range msg.Ports {
			s.ports[p.PortNo] = p
		}
		network.Switches[msg.DPID.String()] = s
		go s.receive()
	}
	network.Unlock()
}

func (sw *OFSwitch) AddInstance(inst interface{}) {
	if actor, ok := inst.(ofp10.ConnectionUpReactor); ok {
		actor.ConnectionUp(sw.DPID())
	}
	sw.appInstance = append(sw.appInstance, inst)
}

func (sw *OFSwitch) SetPort(portNo uint16, port ofp10.PhyPort) {
	sw.portsMu.Lock()
	defer sw.portsMu.Unlock()
	sw.ports[portNo] = port
}

// Returns a pointer to the Switch mapped to dpid.
func Switch(dpid net.HardwareAddr) (*OFSwitch, bool) {
	network.RLock()
	defer network.RUnlock()
	if sw, ok := network.Switches[dpid.String()]; ok {
		return sw, ok
	}
	return nil, false
}

// Returns a slice of *OFPSwitches for operations across all
// switches.
func Switches() []*OFSwitch {
	network.RLock()
	defer network.RUnlock()
	a := make([]*OFSwitch, len(network.Switches))
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
	network.Switches[dpid.String()].stream.Shutdown <- true
	delete(network.Switches, dpid.String())
}

// Returns a slice of all links connected to Switch s.
func (s *OFSwitch) Links() []Link {
	s.linksMu.RLock()
	a := make([]Link, 0)
	for _, v := range s.links {
		a = append(a, *v)
	}
	s.linksMu.RUnlock()
	return a
}

// Returns the link between Switch s and the Switch dpid.
func (s *OFSwitch) Link(dpid net.HardwareAddr) (l Link, ok bool) {
	s.linksMu.RLock()
	if n, k := s.links[dpid.String()]; k {
		l = *n
		ok = true
	}
	s.linksMu.RUnlock()
	return
}

// Updates the link between s.DPID and l.DPID.
func (s *OFSwitch) setLink(dpid net.HardwareAddr, l *Link) {
	s.linksMu.Lock()
	if _, ok := s.links[l.DPID.String()]; !ok {
		log.Println("Link discovered:", dpid, l.Port, l.DPID)
	}
	s.links[l.DPID.String()] = l
	s.linksMu.Unlock()
}

// Returns the dpid of Switch s.
func (s *OFSwitch) DPID() net.HardwareAddr {
	return s.dpid
}

// Returns a slice of all the ports from Switch s.
func (s *OFSwitch) Ports() []ofp10.PhyPort {
	s.portsMu.RLock()
	a := make([]ofp10.PhyPort, len(s.ports))
	i := 0
	for _, v := range s.ports {
		a[i] = v
		i++
	}
	s.portsMu.RUnlock()
	return a
}

// Returns a pointer to the OfpPhyPort at port number from Switch s.
func (sw *OFSwitch) Port(portNo uint16) (port ofp10.PhyPort, ok bool) {
	sw.portsMu.RLock()
	defer sw.portsMu.RUnlock()

	port, ok = sw.ports[portNo]
	return
}

// Sends an OpenFlow message to this Switch.
func (s *OFSwitch) Send(req util.Message) {
	s.stream.Outbound <- req
}

// Receive loop for each Switch.
func (s *OFSwitch) receive() {
	for {
		select {
		case msg := <-s.stream.Inbound:
			// New message has been received from message
			// stream.
			go s.distributeMessages(s.dpid, msg)
		case err := <-s.stream.Error:
			// Message stream has been disconnected.
			for _, app := range s.appInstance {
				if actor, ok := app.(ofp10.ConnectionDownReactor); ok {
					actor.ConnectionDown(s.DPID(), err)
				}
			}
			return
		}
	}
}

func (s *OFSwitch) distributeMessages(dpid net.HardwareAddr, msg util.Message) {
	s.reqsMu.RLock()
	for _, app := range s.appInstance {
		switch t := msg.(type) {
		case *ofp10.SwitchFeatures:
			if actor, ok := app.(ofp10.SwitchFeaturesReplyReactor); ok {
				actor.FeaturesReply(s.DPID(), t)
			}
		case *ofp10.PacketIn:
			if actor, ok := app.(ofp10.PacketInReactor); ok {
				actor.PacketIn(s.DPID(), t)
			}
		case *ofpxx.Header:
			switch t.Header().Type {
			case ofp10.Type_EchoReply:
				if actor, ok := app.(ofp10.EchoReplyReactor); ok {
					actor.EchoReply(s.DPID())
				}
			case ofp10.Type_EchoRequest:
				if actor, ok := app.(ofp10.EchoRequestReactor); ok {
					actor.EchoRequest(s.DPID())
				}
				
			}
		}
	}
	s.reqsMu.RUnlock()
}

// Sends an OpenFlow message to s, and returns a channel to receive
// a response on. Any error encountered during the send except io.EOF
// is returned.
//
// TODO: Actually make work.
func (s *OFSwitch) SendAndReceive(msg util.Message) chan util.Message {
	ch := make(chan util.Message)
	s.Send(msg)
	return ch
}
