package core

import (
	"encoding/binary"
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
)

type MessageStream struct {
	conn *net.TCPConn
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
		0,
		make(chan error, 1),        // Error
		make(chan ofp10.Packet, 1), // Inbound
		make(chan ofp10.Packet, 1), // Outbound
		make(chan bool, 1),         // Shutdown
	}
	go m.outbound()
	go m.inbound()
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

func (m *MessageStream) inbound() {

	cursor := 0
	unreadBytes := make([]byte, 1500)
	unreadByteLength := 0
	for {
		buf := make([]byte, 512)
		if n, err := m.conn.Read(buf); err != nil {
			// Likely a read timeout. Send error to any listening
			// threads. Trigger shutdown to close outbound loop.
			log.Println("InboundError:", err)
			m.Error <- err
			m.Shutdown <- true
			return
		} else {

			copy(unreadBytes, unreadBytes[cursor:])
			copy(unreadBytes[unreadByteLength:], buf)

			cursor = 0
			unreadByteLength = unreadByteLength + n

			// A minimum of 4 bytes should be in the buffer
			for unreadByteLength >= 4 {
				messageLength := int(binary.BigEndian.Uint16(unreadBytes[cursor+2 : cursor+4]))

				if unreadByteLength >= messageLength {
					end := cursor + messageLength
					m.parse(unreadBytes[cursor:end])

					cursor = end
					unreadByteLength = unreadByteLength - messageLength
				} else {
					break
				}
			}
		}
	}
}

func (m *MessageStream) parse(buf []byte) {
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
	m.Inbound <- d
}
