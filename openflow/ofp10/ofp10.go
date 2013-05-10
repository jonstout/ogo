/************************************
 * Jonathan M. Stout 2012
 * ofp10.go
 * OpenFlow 1.0
 ***********************************/
package ofp10

import (
	"io"
	"bytes"
	"encoding/binary"
	"github.com/jonstout/pacit"
)

type Packetish interface {
	io.ReadWriter
	Len() (n uint16)
}

type Packet interface {
	io.ReadWriter
	GetHeader() *Header
}

type Msg struct {
	Data Packet
	DPID string
}

const (
	VERSION = 1
)

type Header struct {
	Version uint8
	Type uint8
	Length uint16
	XID uint32
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

func (h *Header) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, h)
	return 8, err
}

func NewHello() *Header {
	h := NewHeader()
	h.Type = T_HELLO
	return h
}

func NewEchoRequest() *Header {
	h := NewHeader()
	h.Type = T_ECHO_REQUEST
	return h
}

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

// BEGIN: ofp10 - 5.3.6
// ofp_packet_out 1.0
type PacketOut struct {
	Header Header
	BufferID uint32
	InPort uint16
	ActionsLen uint16
	Actions []Packetish//Header
	Data Packetish
}

func NewPacketOut() *PacketOut {
	p := new(PacketOut)
	p.Header = *NewHeader()
	//p.Header.Length = 71
	p.Header.Type = T_PACKET_OUT
	p.BufferID = 0xffffffff
	p.InPort = 0
	//p.ActionsLen = 8
	p.Actions = make([]Packetish,0)
	return p
}

func (p *PacketOut) GetHeader() *Header {
	return &p.Header
}

func (p *PacketOut) Len() (n uint16) {
	n += p.Header.Len()
	for _, e := range p.Actions {
		n += e.Len()
	}
	n += 8
	n += p.Data.Len()
	//if n < 72 { return 72 }
	return
}

func (p *PacketOut) Read(b []byte) (n int, err error) {
	p.Header.Length = p.Len()
	for _, e := range p.Actions {
		p.ActionsLen += e.Len()
	}

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
	actionCount := buf.Len() / 8
	p.Actions = make([]Packetish, actionCount)
	for i := 0; i < actionCount; i++ {
		a := new(ActionOutput)//Header)
		m := 0
		m, err = a.Write(buf.Next(8))
		if m == 0 {
			return
		}
		n += m
		p.Actions[i] = a
	}
	return
}

// ofp_packet_in 1.0
type PacketIn struct {
	Header Header
	BufferID uint32
	TotalLen uint16
	InPort uint16
	Reason uint8
	Data pacit.Ethernet
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
	buf := bytes.NewBuffer(b)
	n, err = p.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.BufferID); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &p.TotalLen); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &p.InPort); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &p.Reason); err != nil {
		return
	}
	n += 1
	//TODO::Parse Data
	m := 0
	p.Data = pacit.Ethernet{}
	if m, err := p.Data.Write(b[n:]); m == 0 {
		return m, err
	}
	n += m
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
