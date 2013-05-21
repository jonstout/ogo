package ogo

import (
	"log"
	"net"
	"time"
	"errors"
	//"bufio"
	//"io"
	//"bytes"
	"encoding/binary"
	"github.com/jonstout/ogo/openflow/ofp10"
)

// A map from DPIDs to all Switches that have connected since
// Ogo started.
var Switches map[string]*Switch

type Switch struct {
	conn net.TCPConn
	outbound chan ofp10.Packet
	DPID net.HardwareAddr
	Ports map[int]ofp10.PhyPort
	requests map[uint32]chan ofp10.Msg
}

// Builds and populates a Switch struct then starts listening
// for OpenFlow messages on conn
func NewOpenFlowSwitch(conn *net.TCPConn) {
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

	if sw, ok := Switches[res.DPID.String()]; ok {
		log.Println("Recovered connection from:", sw.DPID)
		sw.conn = *conn
		go sw.sendSync()
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
		Switches[s.DPID.String()] = s
		go s.sendSync()
		go s.Receive()
	}
}

// Returns a pointer to the Switch mapped to dpid.
func GetSwitch(dpid string) (*Switch, bool) {
	if sw, ok := Switches[dpid]; ok {
		return sw, ok
	} else {
		return nil, false
	}
}

// Disconnects Switch mapped to dpid.
func DisconnectSwitch(dpid string) {
	log.Printf("Closing connection with: %s", dpid)
	Switches[dpid].conn.Close()
	delete(Switches, dpid)
}

// Returns an OfpPhyPort from this Switch
func (s *Switch) GetPort(portNo int) (*ofp10.PhyPort, error) {
	if port, ok := s.Ports[portNo]; ok {
		return &port, nil
	} else {
		return nil, errors.New("ERROR::Undefined port number")
	}
}

// Returns a map of all the OfpPhyPorts from this Switch
func (s *Switch) AllPorts() map[int]ofp10.PhyPort {
	return s.Ports
}

// Sends an OpenFlow message to this Switch
func (s *Switch) Send(req ofp10.Packet) (err error) {
	s.outbound <- req
	return nil
}

func (s *Switch) sendSync() {
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
	parse := make(chan []byte)
	end := make(chan bool)

	go func(parseBuffer chan []byte, end chan bool) {
		buf := <- parseBuffer
		bufLen := len(buf)
		packetLen := int(binary.BigEndian.Uint16(buf[2:4]))
		offset := 0
		for {
			for bufLen >= packetLen {
				switch buf[offset+1] {
				case ofp10.T_PACKET_IN:
					d := new(ofp10.PacketIn)
					n, _ := d.Write(buf[offset:offset+packetLen])
					offset += n
					bufLen = bufLen - n
					m := ofp10.Msg{d, s.DPID.String()}
					s.distributeReceived(m)
				case ofp10.T_HELLO:
					d := new(ofp10.Header)
					n, _ := d.Write(buf[offset:offset+packetLen])
					offset += n
					bufLen = bufLen - n
					m := ofp10.Msg{d, s.DPID.String()}
					s.distributeReceived(m)
				case ofp10.T_ECHO_REPLY:
					d := new(ofp10.Header)
					n, _ := d.Write(buf[offset:offset+packetLen])
					offset += n
					bufLen = bufLen - n
					m := ofp10.Msg{d, s.DPID.String()}
					s.distributeReceived(m)
				case ofp10.T_ECHO_REQUEST:
					d := new(ofp10.Header)
					n, _ := d.Write(buf[offset:offset+packetLen])
					offset += n
					bufLen = bufLen - n
					m := ofp10.Msg{d, s.DPID.String()}
					s.distributeReceived(m)
				default:
					offset += packetLen
					bufLen = bufLen - packetLen
				}
				if bufLen < 4 {
					break
				}
				packetLen = int(binary.BigEndian.Uint16(buf[offset+2:offset+4]))
			}
			select {
			case nextBytes := <- parseBuffer:
				buf = append( append([]byte(nil), buf[offset:]...), nextBytes...)
				bufLen = len(buf)
				offset = 0
				packetLen = int(binary.BigEndian.Uint16(buf[offset+2:offset+4]))
			case <- end:
				break
			}
		}
	}(parse, end)

	for {
		byteSlice := make([]byte, 2048)
		if n, err := s.conn.Read(byteSlice); err != nil {
			DisconnectSwitch(s.DPID.String())
			end <- true
			break
		} else {
			parse <- byteSlice[:n]
		}
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
