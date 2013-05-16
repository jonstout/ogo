package ogo

import (
	"log"
	"net"
	"time"
	"errors"
	//"bufio"
	"io"
	"bytes"
	"github.com/jonstout/ogo/openflow/ofp10"
)

var Switches map[string]*Switch

type Switch struct {
	conn net.TCPConn
	outbound chan ofp10.Packet
	DPID net.HardwareAddr
	Ports map[int]ofp10.PhyPort
	requests map[uint32]chan ofp10.Msg
}

/* Builds and populates Switch struct then starts listening
for OpenFlow messages on conn. */
func NewOpenFlowSwitch(conn *net.TCPConn) {

	if _, err := conn.ReadFrom(ofp10.NewHello()); err != nil {
		log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
		conn.Close()
		return
	}
	if _, err := ofp10.NewHello().ReadFrom(conn); err != nil {
		log.Println("ERROR::Failed with Hello:", err)
		conn.Close()
		return
	}

	if _, err := conn.ReadFrom(ofp10.NewFeaturesRequest()); err != nil {
		log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
		conn.Close()
	}
	res := ofp10.NewFeaturesReply()
	if _, err := res.ReadFrom(conn); err != nil {
		log.Println("ERROR::Failed with FeaturesReply", err)
		conn.Close()
	}

	if sw, ok := Switches[res.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID)
		sw.conn = *conn
		go sw.SendSync()
		go sw.Receive()
	} else {
		log.Printf("Openflow 1.%d Connection: %s", res.Header.Version - 1, res.DPID.String())
		s := new(Switch)
		s.conn = *conn
		s.outbound = make(chan ofp10.Packet)
		s.DPID = res.DPID
		s.Ports = make(map[int]ofp10.PhyPort)
		s.requests = make(map[uint32]chan ofp10.Msg)
		for _, p := range res.Ports {
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
func (s *Switch) GetPort(portNo int) (*ofp10.PhyPort, error) {
	port, ok := s.Ports[portNo]
	if ok != true {
		return nil, errors.New("ERROR::Undefined port number")
	}
	return &port, nil
}

/* Return a map of ints to OfpPhyPorts associated with s. */
func (s *Switch) AllPorts() map[int]ofp10.PhyPort {
	return s.Ports
}

/* Sends an OpenFlow message to s. Any error encountered during the
send except io.EOF is returned. */
func (s *Switch) Send(req ofp10.Packet) (err error) {
	s.outbound <- req
	return nil
}

func (s *Switch) SendSync() {
	for {
		if _, err := s.conn.ReadFrom(<-s.outbound); err != nil {
			log.Println("ERROR::Switch.SendSync::ReadFrom:", err)
			s.conn.Close()
			break
		}
	}
}

/* Receive loop for each Switch. */
func (s *Switch) Receive() {
	//length := uint16(b[2]) << 8 + uint16(b[3])
	parse := make(chan io.Reader)

	go func(parseBuffer chan io.Reader) {
		buf := <- parseBuffer
		if a, ok := buf.(*bytes.Buffer); ok {
			for {
				c := a.Bytes()[:4]
				packetLen := (int(c[2]) << 8) + int(c[3])
				for c[1] >= 4 && a.Len() >= packetLen {
					switch a.Bytes()[1] {
					case ofp10.T_PACKET_IN:
						d := new(ofp10.PacketIn)
						d.ReadFrom(buf)
						m := ofp10.Msg{d, s.DPID.String()}
						s.distributeReceived(m)
					case ofp10.T_HELLO:
						d := new(ofp10.Header)
						d.ReadFrom(buf)
						m := ofp10.Msg{d, s.DPID.String()}
						s.distributeReceived(m)
					case ofp10.T_ECHO_REPLY:
						d := new(ofp10.Header)
						d.ReadFrom(buf)
						m := ofp10.Msg{d, s.DPID.String()}
						s.distributeReceived(m)
					case ofp10.T_ECHO_REQUEST:
						d := new(ofp10.Header)
						d.ReadFrom(buf)
						m := ofp10.Msg{d, s.DPID.String()}
						s.distributeReceived(m)
					}
				}
				a.ReadFrom(<- parseBuffer)
			}
		}
	}(parse)

	for {
		byteSlice := make([]byte, 2500)		
		if _, err := s.conn.Read(byteSlice); err != nil {
			DisconnectSwitch(s.DPID.String())
			break
		}
		parse <- bytes.NewBuffer(byteSlice)
	}
		
}

func (s *Switch) distributeReceived(p ofp10.Msg) {
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
func (s *Switch) SendAndReceive(req ofp10.Packet) (p chan ofp10.Msg, err error) {
	p = make(chan ofp10.Msg)
	s.requests[req.GetHeader().XID] = p
	err = s.Send(req)
	if err != nil {
		delete(s.requests, req.GetHeader().XID)
		return nil, err
	}
	return
}
