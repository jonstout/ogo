package ogo

import (
	"errors"
	"encoding/binary"
	"github.com/jonstout/ogo/openflow/ofp10"
	"log"
	"net"
	"time"
)

type MessageStream struct {
	connection net.TCPConn
	errorMessage chan error
	newMessages chan ofp10.Packet
	parseRoutines chan int
	parsedMessage chan ofp10.Packet
}

func NewMessageStream(conn net.TCPConn) *MessageStream {
	m := &MessageStream{conn,
		make(chan error),
		make(chan ofp10.Packet),
		make(chan int, 10),
		make(chan ofp10.Packet)}
	go m.loop()
	return m
}

func (m *MessageStream) Close() {
	m.errorMessage <- errors.New("Stream closed by external process")
}

func (m *MessageStream) Updates() <- chan ofp10.Packet {
	return m.newMessages
}

func (m *MessageStream) loop() {
	go func() {
		for {
			select {
			case <- m.errorMessage: // Close the m.newMessages chan to end Updates()
				close(m.newMessages)
				m.connection.Close()
				return
			case msg := <- m.parsedMessage: // Forward parsed messages to Updates()
				m.newMessages <- msg
			}
		}
	}()

	cursor := 0
	unreadBytes := make([]byte, 0)
	for {
		buf := make([]byte, 2048)
		if n, err := m.connection.Read(buf); err != nil {
			log.Println("Read timeout...")
			m.errorMessage <- err
			return
		} else {
			cursor = 0
			unreadBytes = append( append([]byte(nil), unreadBytes[cursor:]...), buf[:n]...)
			unreadByteLength := len(unreadBytes)
			// A minimum of 4 bytes should be in the buffer
			for unreadByteLength >= 4 {
				messageLength := int( binary.BigEndian.Uint16(unreadBytes[cursor+2:cursor+4]) )
				if unreadByteLength >= messageLength {
					end := cursor + messageLength
					// Run m.parse upto the size of m.parseRoutines at a time
					m.parseRoutines <- 1
					go m.parse(unreadBytes[cursor:end])
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
		d = new(ofp10.Header)
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
	select {
	case m.parsedMessage <- d:
		<- m.parseRoutines
	case <- time.After(time.Millisecond * 100):
	}
}
