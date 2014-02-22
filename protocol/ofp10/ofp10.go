// OpenFlow Wire Protocol 0x01
// Package ofp10 provides OpenFlow 1.0 structs along with Read
// and Write methods for each.
//
// Struct documentation is taken from the OpenFlow Switch
// Specification Version 1.0.0.
// https://www.opennetworking.org/images/stories/downloads/sdn-resources/onf-specifications/openflow/openflow-spec-v1.0.0.pdf
package ofp10

import (
	//"fmt"
	"bytes"
	"encoding/binary"
	"github.com/jonstout/ogo/protocol/eth"
	"io"
	"net"
)

type Packetish interface {
	io.ReadWriter
	Len() (n uint16)
}

// Packet is any OpenFlow packet that includes a header.
type Packet interface {
	io.ReadWriter
	GetHeader() *Header
}

// Msg is any Packet with its originating DPID.
type Msg struct {
	Data Packet
	DPID net.HardwareAddr
}

const (
	VERSION = 1
)

// The version specifies the OpenFlow protocol version being
// used. During the current draft phase of the OpenFlow
// Protocol, the most significant bit will be set to indicate an
// experimental version and the lower bits will indicate a
// revision number. The current version is 0x01. The final
// version for a Type 0 switch will be 0x00. The length field
// indicates the total length of the message, so no additional
// framing is used to distinguish one frame from the next.
type Header struct {
	Version uint8
	Type    uint8
	Length  uint16
	XID     uint32
}

var NewHeader func() *Header = newHeaderGenerator()

func newHeaderGenerator() func() *Header {
	var xid uint32 = 1
	return func() *Header {
		p := new(Header)
		p.Version = 1
		p.Type = 0
		p.Length = 8
		p.XID = xid
		xid += 1

		return p
	}
}

func (h *Header) GetHeader() *Header {
	return h
}

func (h *Header) Len() (n uint16) {
	return 8
}

func (h *Header) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (h *Header) ReadFrom(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, binary.BigEndian, &h.Version); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &h.Type); err != nil {
		return
	}
	n += 1
	if err = binary.Read(r, binary.BigEndian, &h.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(r, binary.BigEndian, &h.XID); err != nil {
		return
	}
	n += 4
	return
}

func (h *Header) Write(b []byte) (n int, err error) {
	r := bytes.NewBuffer(b)
	binary.Read(r, binary.BigEndian, &h.Version)
	n += 1
	binary.Read(r, binary.BigEndian, &h.Type)
	n += 1
	binary.Read(r, binary.BigEndian, &h.Length)
	n += 2
	binary.Read(r, binary.BigEndian, &h.XID)
	n += 4
	return n, err
}

func NewHello() *Header {
	h := NewHeader()
	h.Type = T_HELLO
	return h
}

// Echo request/reply messages can be sent from either the
// switch or the controller, and must return an echo reply. They
// can be used to indicate the latency, bandwidth, and/or
// liveness of a controller-switch connection.
func NewEchoRequest() *Header {
	h := NewHeader()
	h.Type = T_ECHO_REQUEST
	return h
}

// Echo request/reply messages can be sent from either the
// switch or the controller, and must return an echo reply. They
// can be used to indicate the latency, bandwidth, and/or
// liveness of a controller-switch connection.
func NewEchoReply() *Header {
	h := NewHeader()
	h.Type = T_ECHO_REPLY
	return h
}

// ofp_type 1.0
const (
	/* Immutable messages. */
	T_HELLO = iota
	T_ERROR
	T_ECHO_REQUEST
	T_ECHO_REPLY
	T_VENDOR

	/* Switch configuration messages. */
	T_FEATURES_REQUEST
	T_FEATURES_REPLY
	T_GET_CONFIG_REQUEST
	T_GET_CONFIG_REPLY
	T_SET_CONFIG

	/* Asynchronous messages. */
	T_PACKET_IN
	T_FLOW_REMOVED
	T_PORT_STATUS

	/* Controller command messages. */
	T_PACKET_OUT
	T_FLOW_MOD
	T_PORT_MOD

	/* Statistics messages. */
	T_STATS_REQUEST
	T_STATS_REPLY

	/* Barrier messages. */
	T_BARRIER_REQUEST
	T_BARRIER_REPLY

	/* Queue Configuration messages. */
	T_QUEUE_GET_CONFIG_REQUEST
	T_QUEUE_GET_CONFIG_REPLY
)

// When the controller wishes to send a packet out through the
// datapath, it uses the OFPT_PACKET_OUT message: The buffer_id
// is the same given in the ofp_packet_in message. If the
// buffer_id is -1, then the packet data is included in the data
// array. If OFPP_TABLE is specified as the output port of an
// action, the in_port in the packet_out message is used in the
// flow table lookup.
type PacketOut struct {
	Header     Header
	BufferID   uint32
	InPort     uint16
	ActionsLen uint16
	Actions    []Action
	Data       Packetish
}

func NewPacketOut() *PacketOut {
	p := new(PacketOut)
	p.Header = *NewHeader()
	p.Header.Type = T_PACKET_OUT
	p.BufferID = 0xffffffff
	p.InPort = P_NONE
	p.ActionsLen = 0
	p.Actions = make([]Action, 0)
	return p
}

func (p *PacketOut) AddAction(act Action) {
	p.Actions = append(p.Actions, act)
	p.ActionsLen += act.Len()
}

func (p *PacketOut) GetHeader() *Header {
	return &p.Header
}

func (p *PacketOut) Len() (n uint16) {
	n += p.Header.Len()
	n += p.ActionsLen
	n += 8
	n += p.Data.Len()
	//if n < 72 { return 72 }
	return
}

func (p *PacketOut) Read(b []byte) (n int, err error) {
	p.Header.Length = p.Len()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(&p.Header)
	binary.Write(buf, binary.BigEndian, p.BufferID)
	binary.Write(buf, binary.BigEndian, p.InPort)
	binary.Write(buf, binary.BigEndian, p.ActionsLen)
	for _, e := range p.Actions {
		_, err = buf.ReadFrom(e)
	}
	_, err = buf.ReadFrom(p.Data)

	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (p *PacketOut) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = p.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.BufferID); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &p.InPort); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &p.ActionsLen); err != nil {
		return
	}
	n += 2
	/*actionCount := buf.Len() / 8
	//p.Actions = make([]Action, actionCount)
	for i := 0; i < actionCount; i++ {
		a := new(ActionOutput)//Header)
		m := 0
		m, err = a.Write(buf.Next(8))
		if m == 0 {
			return
		}
		n += m
		p.Actions[i] = a
	}*/
	return
}

// ofp_packet_in 1.0
type PacketIn struct {
	Header   Header
	BufferID uint32
	TotalLen uint16
	InPort   uint16
	Reason   uint8
	Data     eth.Ethernet
}

func (p *PacketIn) GetHeader() *Header {
	return &p.Header
}

func (p *PacketIn) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (p *PacketIn) Write(b []byte) (n int, err error) {
	//buf := bytes.NewBuffer(b)
	n, err = p.Header.Write(b[:8])
	/*if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.BufferID); err != nil {
		return
	}*/
	p.BufferID = binary.BigEndian.Uint32(b[8:12])
	n += 4
	/*if err = binary.Read(buf, binary.BigEndian, &p.TotalLen); err != nil {
		return
	}*/
	p.TotalLen = binary.BigEndian.Uint16(b[12:14])
	n += 2
	/*if err = binary.Read(buf, binary.BigEndian, &p.InPort); err != nil {
		return
	}*/
	p.InPort = binary.BigEndian.Uint16(b[14:16])
	n += 2
	/*if err = binary.Read(buf, binary.BigEndian, &p.Reason); err != nil {
		return
	}*/
	p.Reason = b[16]
	n += 1
	//TODO::Parse Data
	p.Data = eth.Ethernet{}
	if m, err := p.Data.Write(b[n:]); m == 0 {
		return m, err
	} else {
		n += m
	}
	return
}

// ofp_packet_in_reason 1.0
const (
	R_NO_MATCH = iota
	R_ACTION
)

// ofp_vendor_header 1.0
type VendorHeader struct {
	Header Header /*Type OFPT_VENDOR*/
	Vendor uint32
}

func (v *VendorHeader) GetHeader() *Header {
	return &v.Header
}

func (v *VendorHeader) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (v *VendorHeader) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = v.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &v.Vendor); err != nil {
		return
	}
	n += 4
	return
}
