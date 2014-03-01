// OpenFlow Wire Protocol 0x01
// Package ofp10 provides OpenFlow 1.0 structs along with Read
// and Write methods for each.
//
// Struct documentation is taken from the OpenFlow Switch
// Specification Version 1.0.0.
// https://www.opennetworking.org/images/stories/downloads/sdn-resources/onf-specifications/openflow/openflow-spec-v1.0.0.pdf
package ofp10

import (
	"encoding/binary"
	"errors"

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
	Type_Packet_In
	T_FLOW_REMOVED
	T_PORT_STATUS

	/* Controller command messages. */
	Type_Packet_Out
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
	ofpxx.Header
	BufferId   uint32
	InPort     uint16
	ActionsLen uint16
	Actions    []Action
	Data       util.Message
}

func NewPacketOut() *PacketOut {
	p := new(PacketOut)
	p.Header = ofpxx.NewOfp10Header()
	p.Header.Type = Type_Packet_Out
	p.BufferId = 0xffffffff
	p.InPort = P_NONE
	p.ActionsLen = 0
	p.Actions = make([]Action, 0)
	return p
}

func (p *PacketOut) AddAction(act Action) {
	p.Actions = append(p.Actions, act)
	p.ActionsLen += act.Len()
}

func (p *PacketOut) Len() (n uint16) {
	n += p.Header.Len()
	n += 8
	n += p.ActionsLen
	n += p.Data.Len()
	//if n < 72 { return 72 }
	return
}

func (p *PacketOut) MarshalBinary() (data []byte, err error) {
	p.Header.Length = p.Len()

	data, err = p.Header.MarshalBinary()

	b := make([]byte, 4)
	n := 0
	binary.BigEndian.PutUint32(b, p.BufferId)
	n += 4
	binary.BigEndian.PutUint16(b[n:], p.InPort)
	n += 2
	binary.BigEndian.PutUint16(b[n:], p.ActionsLen)
	data = append(data, b...)

	for _, a := range p.Actions {
		b, err = a.MarshalBinary()
		data = append(data, b...)
	}

	b, err = p.Data.MarshalBinary()
	data = append(data, b...)
	return
}

func (p *PacketOut) UnmarshalBinary(data []byte) error {
	err := p.Header.UnmarshalBinary(data)
	n := p.Header.Len()

	p.BufferId = binary.BigEndian.Uint32(data[n:])
	n += 4
	p.InPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	p.ActionsLen = binary.BigEndian.Uint16(data[n:])
	n += 2

	for n < (n + p.ActionsLen) {
		a := DecodeAction(data[n:])
		p.Actions = append(p.Actions, a)
		n += a.Len()
	}

	err = p.Data.UnmarshalBinary(data[n:])
	return err
}

// ofp_packet_in 1.0
type PacketIn struct {
	ofpxx.Header
	BufferId uint32
	TotalLen uint16
	InPort   uint16
	Reason   uint8
	Data     eth.Ethernet
}

func NewPacketIn() *PacketIn {
	p := new(PacketIn)
	p.Header = ofpxx.NewOfp10Header()
	p.Header.Type = Type_Packet_In
	p.BufferId = 0xffffffff
	p.InPort = P_NONE
	p.Reason = 0
	return p
}

func (p *PacketIn) Len() (n uint16) {
	n += p.Header.Len()
	n += 9
	n += p.Data.Len()
	return
}

func (p *PacketIn) MarshalBinary() (data []byte, err error) {
	data, err = p.Header.MarshalBinary()

	b := make([]byte, 9)
	n := 0
	binary.BigEndian.PutUint32(b, p.BufferId)
	n += 4
	binary.BigEndian.PutUint16(b[n:], p.TotalLen)
	n += 2
	binary.BigEndian.PutUint16(b[n:], p.InPort)
	n += 2
	b[n] = p.Reason
	n += 1
	data = append(data, b...)

	b, err = p.Data.MarshalBinary()
	data = append(data, b...)
	return
}

func (p *PacketIn) UnmarshalBinary(data []byte) error {
	err := p.Header.UnmarshalBinary(data)
	n := p.Header.Len()

	p.BufferId = binary.BigEndian.Uint32(data[n:])
	n += 4
	p.TotalLen = binary.BigEndian.Uint16(data[n:])
	n += 2
	p.InPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	p.Reason = data[n]
	n += 1

	err = p.Data.UnmarshalBinary(data[n:])
	return err
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

func (v *VendorHeader) Len() (n uint16) {
	return v.Header.Len() + 4
}

func (v *VendorHeader) MarshalBinary() (data []byte, err error) {
	data, err = v.Header.MarshalBinary()

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(data[:4], v.Vendor)

	data = append(data, b...)
	return
}

func (v *VendorHeader) UnmarshalBinary(data []byte) error {
	if len(data) < int(v.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"VendorHeader message.")
	}
	v.Header.UnmarshalBinary(data)
	n := int(v.Header.Len())
	v.Vendor = binary.BigEndian.Uint32(data[n:])
	return nil
}
