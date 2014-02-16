package core

import (
	"encoding/binary"
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
	"bytes"
)

type MessageBuffer struct {
	Empty chan *bytes.Buffer
	Full chan *bytes.Buffer
}

func NewMessageBuffer() *MessageBuffer {
	m := new(MessageBuffer)
	m.Empty = make(chan *bytes.Buffer, 50)
	m.Full = make(chan *bytes.Buffer, 50)

	for i := 0; i < 50; i++ {
		m.Empty <- bytes.NewBuffer(make([]byte, 0, 2048))
	}
	return m
}

func (b *MessageBuffer) ReadFrom(conn *net.TCPConn) error {
	msg := 0
	hdr := 0
	hdrBuf := make([]byte, 4)

	tmp := make([]byte, 2048)
	buf := <- b.Empty
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			log.Println("InboundError", err)
			return err
		}		
		
		for i := 0; i < n; i++ {
			if hdr < 4 {
				hdrBuf[hdr] = tmp[i]
				buf.WriteByte(tmp[i])
				hdr += 1
				if hdr >= 4 {
					msg = int(binary.BigEndian.Uint16(hdrBuf[2:])) - 4
				}
				continue
			}
			if msg > 0 {
				buf.WriteByte(tmp[i])
				msg = msg - 1
				if msg == 0 {
					hdr = 0
					b.Full <- buf
					buf = <- b.Empty
				}
				continue
			}
		}
	}
}

type MessageStream struct {
	conn *net.TCPConn
	Buffer *MessageBuffer
	// OpenFlow Version
	Version uint8
	// Channel on which to publish connection errors
	Error chan error
	// Channel on which to publish inbound messages
	Inbound chan ofp10.Packet
	// Channel on which to receive outbound messages
	Outbound chan ofp10.Packet
	// Channel on which to receive a shutdown command
	Shutdown chan bool
}

// Returns a pointer to a new MessageStream. Used to parse
// OpenFlow messages from conn.
func NewMessageStream(conn *net.TCPConn) *MessageStream {
	m := &MessageStream{
		conn,
		NewMessageBuffer(),
		0,
		make(chan error, 1),        // Error
		make(chan ofp10.Packet, 1), // Inbound
		make(chan ofp10.Packet, 1), // Outbound
		make(chan bool, 1),         // Shutdown
	}

	go m.outbound()
	go m.Buffer.ReadFrom(conn)

	for i := 0; i < 25; i++ {
		go m.parse()
	}
	return m
}

func (m *MessageStream) GetAddr() net.Addr {
	return m.conn.RemoteAddr()
}

// Listen for a Shutdown signal or Outbound messages.
func (m *MessageStream) outbound() {
	for {
		select {
		case <-m.Shutdown:
			log.Println("Closing OpenFlow message stream.")
			m.conn.Close()
			return
		case msg := <-m.Outbound:
			// Forward outbound messages to conn
			if _, err := m.conn.ReadFrom(msg); err != nil {
				log.Println("OutboundError:", err)
				m.Error <- err
				m.Shutdown <- true
			}
		}
	}
}

func (m *MessageStream) parse() {
	for {
		b := <- m.Buffer.Full
		buf := b.Bytes()
		var d ofp10.Packet

		switch buf[1] {
		case ofp10.T_PACKET_IN:
			d = new(ofp10.PacketIn)
			d.Write(buf)
		case ofp10.T_HELLO:
			d = new(ofp10.Header)
			d.Write(buf)
		case ofp10.T_ECHO_REPLY:
			d = new(ofp10.Header)
			d.Write(buf)
		case ofp10.T_ECHO_REQUEST:
			d = new(ofp10.Header)
			d.Write(buf)
		case ofp10.T_ERROR:
			d = new(ofp10.ErrorMsg)
			d.Write(buf)
		case ofp10.T_VENDOR:
			d = new(ofp10.VendorHeader)
			d.Write(buf)
		case ofp10.T_FEATURES_REPLY:
			d = new(ofp10.SwitchFeatures)
			d.Write(buf)
		case ofp10.T_GET_CONFIG_REPLY:
			d = new(ofp10.SwitchConfig)
			d.Write(buf)
		case ofp10.T_FLOW_REMOVED:
			d = new(ofp10.FlowRemoved)
			d.Write(buf)
		case ofp10.T_PORT_STATUS:
			d = new(ofp10.PortStatus)
			d.Write(buf)
		case ofp10.T_STATS_REPLY:
			d = new(ofp10.StatsReply)
			d.Write(buf)
		case ofp10.T_BARRIER_REPLY:
			d = new(ofp10.Header)
			d.Write(buf)
		default:
			// Unrecognized packet do nothing
		}
		b.Reset()
		m.Buffer.Empty <- b
		m.Inbound <- d
	}
}
