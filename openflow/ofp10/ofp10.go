/************************************
 * Jonathan M. Stout 2012
 * ofp10.go
 * OpenFlow 1.0
 ***********************************/
package ofp10

import (
	"io"
	//"log"
	"bytes"
	"encoding/binary"
	"github.com/jonstout/pacit"
)

type OfpPacket interface {
	io.ReadWriter
	GetHeader() *OfpHeader
}

type OfpMsg struct {
	Data OfpPacket
	DPID string
}

const (
	VERSION = 1
)

type OfpHeader struct {
	Version uint8
	Type uint8
	Length uint16
	XID uint32
}

var NewHeader func() *OfpHeader = newOfpHeaderGenerator()

func newOfpHeaderGenerator() func() *OfpHeader {
	var xid uint32 = 1
	return func() *OfpHeader {
		p := new(OfpHeader)
		p.Version = 1
		p.Type = 0
		p.Length = 8
		p.XID = xid
		xid += 1
		
		return p
	}
}

func (h *OfpHeader) GetHeader() *OfpHeader {
	return h
}

func (h *OfpHeader) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, h)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (h *OfpHeader) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, h)
	return 8, err
}

func NewHello() *OfpHeader {
	h := NewHeader()
	h.Type = OFPT_HELLO
	return h
}

func NewEchoRequest() *OfpHeader {
	h := NewHeader()
	h.Type = OFPT_ECHO_REQUEST
	return h
}

func NewEchoReply() *OfpHeader {
	h := NewHeader()
	h.Type = OFPT_ECHO_REPLY
	return h
}

// ofp_type 1.0
const (
	/* Immutable messages. */
	OFPT_HELLO = iota
	OFPT_ERROR
	OFPT_ECHO_REQUEST
	OFPT_ECHO_REPLY
	OFPT_VENDOR

	/* Switch configuration messages. */
	OFPT_FEATURES_REQUEST
	OFPT_FEATURES_REPLY
	OFPT_GET_CONFIG_REQUEST
	OFPT_GET_CONFIG_REPLY
	OFPT_SET_CONFIG

	/* Asynchronous messages. */
	OFPT_PACKET_IN
	OFPT_FLOW_REMOVED
	OFPT_PORT_STATUS

	/* Controller command messages. */
	OFPT_PACKET_OUT
	OFPT_FLOW_MOD
	OFPT_PORT_MOD

	/* Statistics messages. */
	OFPT_STATS_REQUEST
	OFPT_STATS_REPLY

	/* Barrier messages. */
	OFPT_BARRIER_REQUEST
	OFPT_BARRIER_REPLY

	/* Queue Configuration messages. */
	OFPT_QUEUE_GET_CONFIG_REQUEST
	OFPT_QUEUE_GET_CONFIG_REPLY
)

// BEGIN: ofp10 - 5.3.6
// ofp_packet_out 1.0
type OfpPacketOut struct {
	Header OfpHeader
	BufferID uint32
	InPort uint16
	ActionsLen uint16
	Actions []OfpActionOutput//Header
	Data io.ReadWriter
}

func NewPacketOut() *OfpPacketOut {
	p := new(OfpPacketOut)
	p.Header = *NewHeader()
	p.Header.Length = 71
	p.Header.Type = OFPT_PACKET_OUT
	p.BufferID = 0xffffffff
	p.InPort = 0
	p.ActionsLen = 8
	p.Actions = make([]OfpActionOutput,0)
	return p
}

func (p *OfpPacketOut) GetHeader() *OfpHeader {
	return &p.Header
}

func (p *OfpPacketOut) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(&p.Header)
	binary.Write(buf, binary.BigEndian, p.BufferID)
	binary.Write(buf, binary.BigEndian, p.InPort)
	binary.Write(buf, binary.BigEndian, p.ActionsLen)
	binary.Write(buf, binary.BigEndian, p.Actions)
	_, err = buf.ReadFrom(p.Data)
	
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (p *OfpPacketOut) Write(b []byte) (n int, err error) {
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
	p.Actions = make([]OfpActionOutput, actionCount)
	for i := 0; i < actionCount; i++ {
		a := new(OfpActionOutput)//Header)
		m := 0
		m, err = a.Write(buf.Next(8))
		if m == 0 {
			return
		}
		n += m
		p.Actions[i] = *a
	}
	return
}

// ofp_packet_in 1.0
type OfpPacketIn struct {
	Header OfpHeader
	BufferID uint32
	TotalLen uint16
	InPort uint16
	Reason uint8
	Data interface{}
}

func (p *OfpPacketIn) GetHeader() *OfpHeader {
	return &p.Header
}

func (p *OfpPacketIn) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (p *OfpPacketIn) Write(b []byte) (n int, err error) {
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
	e := new(pacit.Ethernet)
	if m, err = e.Write(b[n:]); m == 0 {
		return
	}
	switch e.Ethertype {
	case pacit.ARP_MSG:
		d := new(pacit.ARP)
		if m, err = d.Write(buf.Bytes()); m == 0 {
			return
		}
		p.Data = d
	case pacit.LLDP_MSG:
		d := new(pacit.LLDP)
		if m, err = d.Write(buf.Bytes()); m == 0 {
			return
		}
		p.Data = d
	}
	return
}

// ofp_packet_in_reason 1.0
const (
	OFPR_NO_MATCH = iota
	OFPR_ACTION
)

// ofp_vendor_header 1.0
type OfpVendorHeader struct {
	Header OfpHeader /*Type OFPT_VENDOR*/
	Vendor uint32
}

func (v *OfpVendorHeader) GetHeader() *OfpHeader {
	return &v.Header
}

func (v *OfpVendorHeader) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (v *OfpVendorHeader) Write(b []byte) (n int, err error) {
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
