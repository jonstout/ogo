package ofp10

import (
	"net"
	"bytes"
	"errors"
	"encoding/binary"
)

type OfpSwitchFeatures struct {
	Header OfpHeader
	//DPID uint64
	DPID net.HardwareAddr
	//DPID [8]uint8
	Buffers uint32
	Tables uint8
	Pad [3]uint8

	Capabilities uint32
	Actions uint32

	Ports []OfpPhyPort
}
// OFP_ASSERT(len(OfpSwitchFeatures) == 32)

// FeaturesRequest constructor
func NewFeaturesRequest() *OfpHeader {
	req := NewHeader()
	req.Type = OFPT_FEATURES_REQUEST
	return req
}

// FeaturesReply constructor
func NewFeaturesReply() *OfpSwitchFeatures {
	res := new(OfpSwitchFeatures)
	res.Header.Type = OFPT_FEATURES_REPLY
	return res
}

func (f *OfpSwitchFeatures) GetHeader() *OfpHeader {
	return &f.Header
}

func (f *OfpSwitchFeatures) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, f)
	n, err = buf.Read(b)
	return
}

func (f *OfpSwitchFeatures) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = f.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
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
	
	// Verify port data structures are the correct size.
	if buf.Len() % 48 != 0 {
		return n, errors.New("Ports recieved are malformed.")
	}
	portCount := buf.Len() / 48
	f.Ports = make([]OfpPhyPort, portCount)
	for i := 0; i < portCount; i++ {
		p := new(OfpPhyPort)
		m, portErr := p.Write(buf.Next(48))
		if portErr != nil {
			return n, portErr
		}
		n += m
		f.Ports[i] = *p
	}
	return
}

// ofp_capabilities 1.0
const (
	OFPC_FLOW_STATS = 1 << 0
	OFPC_TABLE_STATS = 1 << 1
	OFPC_PORT_STATS = 1 << 2
	OFPC_STP = 1 << 3
	OFPC_RESERVED = 1 << 4
	OFPC_IP_REASM = 1 << 5
	OFPC_QUEUE_STATS = 1 << 6
	OFPC_ARP_MATCH_IP = 1 << 7
)
