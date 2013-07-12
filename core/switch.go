package core

import (
	"log"
	"net"
	"time"
	"errors"
	//"bufio"
	//"io"
	//"bytes"
	//"encoding/binary"
	"github.com/jonstout/ogo/openflow/ofp10"
)

/* A map from DPIDs to all Switches that have connected since
Ogo started. */
var switches map[string]*OFPSwitch

type OFPSwitch struct {
	conn net.TCPConn
	messageStream *MessageStream
	outbound chan ofp10.Packet
	DPID net.HardwareAddr
	Ports map[int]ofp10.PhyPort
	requests map[uint32]chan ofp10.Msg
}

/* Builds and populates a Switch struct then starts listening
for OpenFlow messages on conn. */
func NewOFPSwitch(conn *net.TCPConn) {
	if _, err := conn.ReadFrom(ofp10.NewHello()); err != nil {
		log.Println("Could not send initial Hello message", err)
		conn.Close()
		return
	}
	if _, err := ofp10.NewHello().ReadFrom(conn); err != nil {
		log.Println("Did not receive Hello message from connection", err)
		conn.Close()
		return
	}

	if _, err := conn.ReadFrom(ofp10.NewFeaturesRequest()); err != nil {
		log.Println("Could not send initial Features Request", err)
		conn.Close()
		return
	}
	res := ofp10.NewFeaturesReply()
	if _, err := res.ReadFrom(conn); err != nil {
		log.Println("Did not receive Features Reply", err)
		conn.Close()
		return
	}

	if sw, ok := switches[res.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID)
		sw.conn = *conn
		sw.messageStream = NewMessageStream(*conn)
		go sw.sendSync()
		go sw.Receive()
	} else {
		log.Printf("Openflow 1.%d Connection: %s", res.Header.Version - 1, res.DPID.String())
		s := new(OFPSwitch)
		s.conn = *conn
		s.outbound = make(chan ofp10.Packet)
		s.DPID = res.DPID
		s.Ports = make(map[int]ofp10.PhyPort)
		s.requests = make(map[uint32]chan ofp10.Msg)
		for _, p := range res.Ports {
			s.Ports[int(p.PortNo)] = p
		}
		s.messageStream = NewMessageStream(*conn)
		switches[s.DPID.String()] = s
		go s.sendSync()
		go s.Receive()
	}
}

/* Returns a pointer to the Switch mapped to dpid. */
func Switch(dpid string) (*OFPSwitch, bool) {
	if sw, ok := switches[dpid]; ok {
		return sw, ok
	} else {
		return nil, false
	}
}

/* Disconnects Switch mapped to dpid. */
func DisconnectSwitch(dpid string) {
	log.Printf("Closing connection with: %s", dpid)
	switches[dpid].conn.Close()
	delete(switches, dpid)
}

/* Returns a pointer to portNo OfpPhyPort from this Switch. */
func (s *OFPSwitch) Port(portNo int) (*ofp10.PhyPort, error) {
	if port, ok := s.Ports[portNo]; ok {
		return &port, nil
	} else {
		return nil, errors.New("ERROR::Undefined port number")
	}
}

/* Returns a map of all the OfpPhyPorts from this Switch. */
func (s *OFPSwitch) AllPorts() map[int]ofp10.PhyPort {
	return s.Ports
}

/* Sends an OpenFlow message to this Switch. */
func (s *OFPSwitch) Send(req ofp10.Packet) (err error) {
	s.outbound <- req
	return nil
}

func (s *OFPSwitch) sendSync() {
	for {
		if _, err := s.conn.ReadFrom(<-s.outbound); err != nil {
			log.Println("Closing connection from", s.DPID)
			s.conn.Close()
			s.messageStream.Close()
			break
		}
	}
}

/* Receive loop for each Switch. */
func (s *OFPSwitch) Receive() {
	for p := range s.messageStream.Updates() {
		s.distributeReceived( ofp10.Msg{p, s.DPID.String()} )
	}
}

func (s *OFPSwitch) distributeReceived(p ofp10.Msg) {
	h := p.Data.GetHeader()
	if pktChan, ok := s.requests[h.XID]; ok {
		select {
		case pktChan <- p:
		case <- time.After(time.Millisecond * 100):
		}
		delete(s.requests, h.XID)
	} else {
		for _, ch := range messageChans[h.Type] {
			select {
			case ch <- p:
			case <- time.After(time.Millisecond * 100):
			}
		}
	}
}

/* Sends an OpenFlow message to s, and returns a channel to receive 
a response on. Any error encountered during the send except io.EOF
is returned. */
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
