package ogo

import (
	"log"
	"net"
	"errors"
	"github.com/jonstout/ogo/openflow/ofp10"
)

var Switches map[string]*Switch

type Switch struct {
	conn net.TCPConn
	outbound chan ofp10.OfpPacket
	DPID net.HardwareAddr
	Ports map[int]ofp10.OfpPhyPort
	requests map[uint32]chan ofp10.OfpMsg
}

/* Builds and populates Switch struct then starts listening
for OpenFlow messages on conn. */
func NewOpenFlowSwitch(conn *net.TCPConn) {

	if _, err := conn.ReadFrom(ofp10.NewHello()); err != nil {
		log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
		conn.Close()
	}
	buf := make([]byte, 1500)
	n, _ := conn.Read(buf)
	res := ofp10.NewHello()
	res.Write(buf[:n])

	if _, err := conn.ReadFrom(ofp10.NewFeaturesRequest()); err != nil {
		log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
		conn.Close()
	}
	buf2 := make([]byte, 1500)
	fres := ofp10.NewFeaturesReply()
	n, _ = conn.Read(buf2)
	fres.Write(buf2[:n])

	if sw, ok := Switches[fres.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID)
		sw.conn = *conn
		go sw.SendSync()
		go sw.Receive()
	} else {
		log.Printf("Openflow 1.%d Connection: %s", res.Version - 1, fres.DPID.String())
		s := new(Switch)
		s.conn = *conn
		s.DPID = fres.DPID
		s.Ports = make(map[int]ofp10.OfpPhyPort)
		s.requests = make(map[uint32]chan ofp10.OfpMsg)
		for _, p := range fres.Ports {
			s.Ports[int(p.PortNo)] = p
		}
		go s.SendSync()
		go s.Receive()
		Switches[s.DPID.String()] = s
	}
}

/* Returns a pointer to the Switch found at dpid. */
func GetSwitch(dpid string) (*Switch, bool) {
	sw, ok := Switches[dpid]
	if ok != true {
		return nil, false
	}
	return sw, ok
}

/* Disconnects Switch found at dpid. */
func DisconnectSwitch(dpid string) {
	log.Printf("Closing connection with: %s", dpid)
	Switches[dpid].conn.Close()
	delete(Switches, dpid)
}

/* Returns OfpPhyPort p from Switch s. */
func (s *Switch) GetPort(portNo int) (*ofp10.OfpPhyPort, error) {
	port, ok := s.Ports[portNo]
	if ok != true {
		return nil, errors.New("ERROR::Undefined port number")
	}
	return &port, nil
}

/* Return a map of ints to OfpPhyPorts associated with s. */
func (s *Switch) AllPorts() map[int]ofp10.OfpPhyPort {
	return s.Ports
}

/* Sends an OpenFlow message to s. Any error encountered during the
send except io.EOF is returned. */
func (s *Switch) Send(req ofp10.OfpPacket) (err error) {
	go func() {
		s.outbound<- req
	}()
	return nil
}

func (s *Switch) SendSync() {
	s.outbound = make(chan ofp10.OfpPacket)
	for {
		msg := <-s.outbound
		if _, err := s.conn.ReadFrom(msg); err != nil {
			log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
			s.conn.Close()
			break
		}
	}
}

/* Receive loop for each Switch. */
func (s *Switch) Receive() {
	for {
		buf := make([]byte, 750)
		if _, err := s.conn.Read(buf); err != nil {
			log.Println("ERROR::Switch.Receive::Read:", err)
			DisconnectSwitch(s.DPID.String())
			break
		}
		switch buf[1] {
		case ofp10.OFPT_HELLO:
			m := ofp10.OfpMsg{new(ofp10.OfpHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_ERROR:
			m := ofp10.OfpMsg{new(ofp10.OfpErrorMsg), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_ECHO_REPLY:
			m := ofp10.OfpMsg{new(ofp10.OfpHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_ECHO_REQUEST:
			m := ofp10.OfpMsg{new(ofp10.OfpHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_VENDOR:
			m := ofp10.OfpMsg{new(ofp10.OfpVendorHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_FEATURES_REPLY:
			m := ofp10.OfpMsg{new(ofp10.OfpHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_GET_CONFIG_REPLY:
			m := ofp10.OfpMsg{new(ofp10.OfpSwitchConfig), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_PACKET_IN:
			m := ofp10.OfpMsg{new(ofp10.OfpPacketIn), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_FLOW_REMOVED:
			m := ofp10.OfpMsg{new(ofp10.OfpFlowRemoved), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_PORT_STATUS:
			m := ofp10.OfpMsg{new(ofp10.OfpPortStatus), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_STATS_REPLY:
			m := ofp10.OfpMsg{new(ofp10.OfpStatsReply), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		case ofp10.OFPT_BARRIER_REPLY:
			m := ofp10.OfpMsg{new(ofp10.OfpHeader), s.DPID.String()}
			m.Data.Write(buf)
			s.distributeReceived(m)
		default:
		}
	}
}

func (s *Switch) distributeReceived(p ofp10.OfpMsg) {
	h := p.Data.GetHeader()
	if pktChan, ok := s.requests[h.XID]; ok {
		go func() {
			pktChan<- p
			delete(s.requests, h.XID)
			}()
	} else {
		for _, ch := range messageChans[h.Type] {
			go func() {
				ch<- p
			}()
		}
	}
}

/* Sends an OpenFlow message to s, and returns a channel to receive 
a response on. Any error encountered during the send except io.EOF
is returned. */
func (s *Switch) SendAndReceive(req ofp10.OfpPacket) (p chan ofp10.OfpMsg, err error) {
	p = make(chan ofp10.OfpMsg)
	s.requests[req.GetHeader().XID] = p
	err = s.Send(req)
	if err != nil {
		delete(s.requests, req.GetHeader().XID)
		return nil, err
	}
	return
}
