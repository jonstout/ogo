package ofp10

import (
	"encoding/binary"
	//"errors"
	"net"

	"github.com/jonstout/ogo/protocol/ofpxx"
)

type SwitchFeatures struct {
	ofpxx.Header
	DPID net.HardwareAddr // Size 8
	Buffers uint32
	Tables  uint8
	pad     []uint8 // Size 3
	Capabilities uint32
	Actions      uint32

	Ports []PhyPort
}

// FeaturesRequest constructor
func NewFeaturesRequest() *ofpxx.Header {
	req := ofpxx.NewOfp10Header()
	req.Type = Type_Features_Request
	return &req
}

// FeaturesReply constructor
func NewFeaturesReply() *SwitchFeatures {
	res := new(SwitchFeatures)
	res.Header.Type = Type_Features_Reply
	res.DPID = make([]byte, 8)
	res.pad = make([]byte, 3)
	res.Ports = make([]PhyPort, 0)
	return res
}

func (s *SwitchFeatures) Len() (n uint16) {
	n = s.Header.Len()
	n += uint16(len(s.DPID))
	n += 16
	for _, p := range s.Ports {
		n += p.Len()
	}
	return
}

func (s *SwitchFeatures) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	s.Header.Length = s.Len()
	next := 0

	bytes, err := s.Header.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[next:], bytes)
	next = len(bytes)
	binary.BigEndian.PutUint32(data[next:], s.Buffers)
	next += 4
	data[next] = s.Tables
	next += 1
	copy(data, s.pad)
	next += len(s.pad)
	binary.BigEndian.PutUint32(data[next:], s.Capabilities)
	next += 4
	binary.BigEndian.PutUint32(data[next:], s.Actions)
	next += 4
	
	for _, p := range s.Ports {
		bytes, err = p.MarshalBinary()
		if err != nil {
			return
		}
		copy(data[next:], bytes)
		next += len(bytes)
	}
	return
}

func (s *SwitchFeatures) UnmarshalBinary(data []byte) error {
	err := s.Header.UnmarshalBinary(data)
	if err != nil {
		return err
	}
	next := int(s.Header.Len())
	copy(s.DPID, data[next:])
	next += len(s.DPID)
	s.Buffers = binary.BigEndian.Uint32(data[next:])
	next += 4
	s.Tables = data[next]
	next += 1
	copy(s.pad, data[next:])
	next += 3
	s.Capabilities = binary.BigEndian.Uint32(data[next:])
	next += 4
	s.Actions = binary.BigEndian.Uint32(data[next:])
	next += 4

	for _, p := range s.Ports {
		err = p.UnmarshalBinary(data[next:])
		if err != nil {
			return err
		}
		next += int(p.Len())
	}	
	return nil
}

// ofp_capabilities 1.0
const (
	C_FLOW_STATS   = 1 << 0
	C_TABLE_STATS  = 1 << 1
	C_PORT_STATS   = 1 << 2
	C_STP          = 1 << 3
	C_RESERVED     = 1 << 4
	C_IP_REASM     = 1 << 5
	C_QUEUE_STATS  = 1 << 6
	C_ARP_MATCH_IP = 1 << 7
)
