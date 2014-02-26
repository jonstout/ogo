// OpenFlow Wire Protocol 0x01
// Package ofp10 provides OpenFlow 1.0 structs along with Read
// and Write methods for each.
//
// Struct documentation is taken from the OpenFlow Switch
// Specification Version 1.0.0.
// https://www.opennetworking.org/images/stories/downloads/sdn-resources/onf-specifications/openflow/openflow-spec-v1.0.0.pdf
package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/jonstout/ogo/protocol/eth"
	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
)

const (
	VERSION = 1
)

// Echo request/reply messages can be sent from either the
// switch or the controller, and must return an echo reply. They
// can be used to indicate the latency, bandwidth, and/or
// liveness of a controller-switch connection.
func NewEchoRequest() *ofpxx.Header {
	h := ofpxx.NewOfp10Header()
	h.Type = Type_Echo_Request
	return &h
}

// Echo request/reply messages can be sent from either the
// switch or the controller, and must return an echo reply. They
// can be used to indicate the latency, bandwidth, and/or
// liveness of a controller-switch connection.
func NewEchoReply() *ofpxx.Header {
	h := ofpxx.NewOfp10Header()
	h.Type = Type_Echo_Reply
	return &h
}

// ofp_type 1.0
const (
	/* Immutable messages. */
	Type_Hello = iota
	Type_Error
	Type_Echo_Request
	Type_Echo_Reply
	Type_Vendor

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
	Header     ofpxx.Header
	BufferID   uint32
	InPort     uint16
	ActionsLen uint16
	Actions    []Action
	Data       util.Message
}

func NewPacketOut() *PacketOut {
	p := new(PacketOut)
	p.Header = ofpxx.NewOfp10Header()
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

func (p *PacketOut) GetHeader() *ofpxx.Header {
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
	data, _ := p.Header.MarshelBinary()
	buf.Write(data)
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
	err = p.Header.UnmarshelBinary(buf.Next(8))
	n += 8
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
	Header   ofpxx.Header
	BufferID uint32
	TotalLen uint16
	InPort   uint16
	Reason   uint8
	Data     eth.Ethernet
}

func (p *PacketIn) GetHeader() *ofpxx.Header {
	return &p.Header
}

func (p *PacketIn) Len() (n uint16) {
	n += p.Header.Len()
	n += 9
	n += p.Data.Len()
	return
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
	err = p.Header.UnmarshelBinary(buf.Next(8))
	n += 8
	if err = binary.Read(buf, binary.BigEndian, &p.TotalLen); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.InPort); err != nil {
		return
	}
	if err = binary.Read(buf, binary.BigEndian, &p.Reason); err != nil {
		return
	}
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
	Header ofpxx.Header /*Type OFPT_VENDOR*/
	Vendor uint32
}

func (v *VendorHeader) GetHeader() *ofpxx.Header {
	return &v.Header
}

func (v *VendorHeader) Len() (n uint16) {
	n = v.Header.Len()
	n += 4
	return
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
	err = v.Header.UnmarshelBinary(buf.Next(8))
	n += 8
	if err = binary.Read(buf, binary.BigEndian, &v.Vendor); err != nil {
		return
	}
	n += 4
	return
}
