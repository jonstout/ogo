package ofp10

import (
	"net"
	"bytes"
	"errors"
	"encoding/binary"
)

type SwitchFeatures struct {
	Header Header
	//DPID uint64
	DPID net.HardwareAddr
	//DPID [8]uint8
	Buffers uint32
	Tables uint8
	Pad [3]uint8

	Capabilities uint32
	Actions uint32

	Ports []PhyPort
}
// OFP_ASSERT(len(SwitchFeatures) == 32)

// FeaturesRequest constructor
func NewFeaturesRequest() *Header {
	req := NewHeader()
	req.Type = T_FEATURES_REQUEST
	return req
}

// FeaturesReply constructor
func NewFeaturesReply() *SwitchFeatures {
	res := new(SwitchFeatures)
	res.Header.Type = T_FEATURES_REPLY
	return res
}

func (f *SwitchFeatures) GetHeader() *Header {
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
	f.Ports = make([]PhyPort, portCount)
	for i := 0; i < portCount; i++ {
		p := new(PhyPort)
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
	C_FLOW_STATS = 1 << 0
	C_TABLE_STATS = 1 << 1
	C_PORT_STATS = 1 << 2
	C_STP = 1 << 3
	C_RESERVED = 1 << 4
	C_IP_REASM = 1 << 5
	C_QUEUE_STATS = 1 << 6
	C_ARP_MATCH_IP = 1 << 7
)
