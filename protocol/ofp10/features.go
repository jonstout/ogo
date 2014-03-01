package ofp10

import (
	"bytes"
	"encoding/binary"
	//"errors"
	"net"

	"github.com/jonstout/ogo/protocol/ofpxx"
)

type SwitchFeatures struct {
	Header ofpxx.Header
	//DPID uint64
	DPID net.HardwareAddr
	//DPID [8]uint8
	Buffers uint32
	Tables  uint8
	Pad     [3]uint8

	Capabilities uint32
	Actions      uint32

	Ports []PhyPort
}

// OFP_ASSERT(len(SwitchFeatures) == 32)

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

func (f *SwitchFeatures) GetHeader() *ofpxx.Header {
	return &f.Header
}

func (f *SwitchFeatures) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, f)
	n, err = buf.Read(b)
	return
}

func (f *SwitchFeatures) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err := f.Header.UnmarshalBinary(buf.Next(8)); err != nil {
		return 0, err
	}
	n += 8
	dpid := make([]uint8, 8)
	if err = binary.Read(buf, binary.BigEndian, &dpid); err != nil {
		return
	}
	n += 8
	f.DPID = net.HardwareAddr(dpid)
	if err = binary.Read(buf, binary.BigEndian, &f.Buffers); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &f.Tables); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &f.Pad); err != nil {
		return
	}
	n += 3
	if err = binary.Read(buf, binary.BigEndian, &f.Capabilities); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &f.Actions); err != nil {
		return
	}
	n += 4

	f.Ports = make([]PhyPort, 0)
	for n < len(b) {
		p := NewPhyPort()
		p.UnmarshalBinary(b[n:])
		f.Ports = append(f.Ports, *p)
	}
	return
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
