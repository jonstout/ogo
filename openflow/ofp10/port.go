package ofp10

import (
	"io"
	"bytes"
	"encoding/binary"
)

// ofp_phy_port 1.0
type OfpPhyPort struct {
     PortNo uint16
     HWAddr [OFP_ETH_ALEN]uint8
     Name [OFP_MAX_PORT_NAME_LEN]byte
     
     Config uint32
     State uint32
     
     Curr uint32
     Advertised uint32
     Supported uint32
     Peer uint32
}

func (p *OfpPhyPort) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (p *OfpPhyPort) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &p.PortNo)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &p.HWAddr)
	if err != nil {
		return
	}
	n += 6
	err = binary.Read(buf, binary.BigEndian, &p.Name)
	if err != nil {
		return
	}
	n += 16
	err = binary.Read(buf, binary.BigEndian, &p.Config)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.State)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Curr)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Advertised)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Supported)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Peer)
	if err != nil {
		return
	}
	n += 4
	return n, err
}

// ofp_port_mod 1.0
type OfpPortMod struct {
	Header OfpHeader
	PortNo uint16
	HWAddr [OFP_ETH_ALEN]uint8

	Config uint32
	Mask uint32
	Advertise uint32
	Pad [4]uint8
}

func (p *OfpPortMod) GetHeader() *OfpHeader {
	return &p.Header
}

func (p *OfpPortMod) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (p *OfpPortMod) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = p.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &p.PortNo)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &p.HWAddr)
	if err != nil {
		return
	}
	n += 6
	err = binary.Read(buf, binary.BigEndian, &p.Config)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Mask)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Advertise)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &p.Pad)
	if err != nil {
		return
	}
	n += 4
	return n, nil
}

const (
	OFP_ETH_ALEN = 6
	OFP_MAX_PORT_NAME_LEN = 16
)

// ofp_port_config 1.0
const (
	OFPPC_PORT_DOWN = 1 << 0

	OFPPC_NO_STP = 1 << 1
	OFPPC_NO_RECV = 1 << 2
	
	OFPPC_NO_STP_RECV = 1 << 3
	OFPPC_NO_FLOOD = 1 << 4
	OFPPC_NO_FWD = 1 << 5
	OFPPC_NO_PACKET_IN = 1 << 6
)

// ofp_port_state 1.0
const (
	OFPPS_LINK_DOWN = 1 << 0

	OFPPS_STP_LISTEN = 0 << 8 /* Not learning or relaying frames. */
	OFPPS_STP_LEARN = 1 << 8 /* Learning but not relaying frames. */
	OFPPS_STP_FORWARD = 2 << 8 /* Learning and relaying frames. */
	OFPPS_STP_BLOCK = 3 << 8 /* Not part of spanning tree. */
	OFPPS_STP_MASK = 3 << 8 /* Bit mask for OFPPS_STP_* values. */
)

// ofp_port 1.0
const (
	OFPP_MAX = 0Xff00

	OFPP_IN_PORT = 0xfff8
	OFPP_TABLE = 0xfff9
	
	OFPP_NORMAL = 0xfffa
	OFPP_FLOOD = 0xfffb
	
	OFPP_ALL = 0xfffc
	OFPP_CONTROLLER = 0xfffd
	OFPP_LOCAL = 0xfffe
	OFPP_NONE = 0xffff
)

// ofp_port_features 1.0
const (
	OFPPF_10MB_HD = 1 << 0
	OFPPF_10MB_FD = 1 << 1
	OFPPF_100MB_HD = 1 << 2
	OFPPF_100MB_FD = 1 << 3
	OFPPF_1GB_HD = 1 << 4
	OFPPF_1GB_FD = 1 << 5
	OFPPF_10GB_FD = 1 << 6

	OFPPF_COPPER = 1 << 7
	OFPPF_FIBER = 1 << 8
	OFPPF_AUTONEG = 1 << 9
	OFPPF_PAUSE = 1 << 10
	OFPPF_PAUSE_ASYM = 1 << 11
)
// END: ofp10 - 5.2.1
