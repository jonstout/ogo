package ofp10

import (
	"io"
	"bytes"
	"net"
	"encoding/binary"
)

// ofp_phy_port 1.0
type PhyPort struct {
     PortNo uint16
     HWAddr net.HardwareAddr
     Name []byte
     
     Config uint32
     State uint32
     
     Curr uint32
     Advertised uint32
     Supported uint32
     Peer uint32
}

func (p *PhyPort) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (p *PhyPort) ReadFrom(r io.Reader) (n int64, err error) {
	if err = binary.Read(r, binary.BigEndian, &p.PortNo); err != nil {
		return
	}
	n += 2
	p.HWAddr = make([]byte, ETH_ALEN)
	if err = binary.Read(r, binary.BigEndian, &p.HWAddr); err != nil {
		return
	}
	n += int64(ETH_ALEN)
	p.Name = make([]byte, MAX_PORT_NAME_LEN)
	if err = binary.Read(r, binary.BigEndian, &p.Name); err != nil {
		return
	}
	n += int64(MAX_PORT_NAME_LEN)
	if err = binary.Read(r, binary.BigEndian, &p.Config); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, binary.BigEndian, &p.State); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, binary.BigEndian, &p.Curr); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, binary.BigEndian, &p.Advertised); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, binary.BigEndian, &p.Supported); err != nil {
		return
	}
	n += 4
	if err = binary.Read(r, binary.BigEndian, &p.Peer); err != nil {
		return
	}
	n += 4
	return
}

func (p *PhyPort) Write(b []byte) (n int, err error) {
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
type PortMod struct {
	Header Header
	PortNo uint16
	HWAddr [ETH_ALEN]uint8

	Config uint32
	Mask uint32
	Advertise uint32
	Pad [4]uint8
}

func (p *PortMod) GetHeader() *Header {
	return &p.Header
}

func (p *PortMod) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (p *PortMod) Write(b []byte) (n int, err error) {
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
	ETH_ALEN = 6
	MAX_PORT_NAME_LEN = 16
)

// ofp_port_config 1.0
const (
	PC_PORT_DOWN = 1 << 0

	PC_NO_STP = 1 << 1
	PC_NO_RECV = 1 << 2
	
	PC_NO_STP_RECV = 1 << 3
	PC_NO_FLOOD = 1 << 4
	PC_NO_FWD = 1 << 5
	PC_NO_PACKET_IN = 1 << 6
)

// ofp_port_state 1.0
const (
	PS_LINK_DOWN = 1 << 0

	PS_STP_LISTEN = 0 << 8 /* Not learning or relaying frames. */
	PS_STP_LEARN = 1 << 8 /* Learning but not relaying frames. */
	PS_STP_FORWARD = 2 << 8 /* Learning and relaying frames. */
	PS_STP_BLOCK = 3 << 8 /* Not part of spanning tree. */
	PS_STP_MASK = 3 << 8 /* Bit mask for OFPPS_STP_* values. */
)

// ofp_port 1.0
const (
	P_MAX = 0Xff00

	P_IN_PORT = 0xfff8
	P_TABLE = 0xfff9
	
	P_NORMAL = 0xfffa
	P_FLOOD = 0xfffb
	
	P_ALL = 0xfffc
	P_CONTROLLER = 0xfffd
	P_LOCAL = 0xfffe
	P_NONE = 0xffff
)

// ofp_port_features 1.0
const (
	PF_10MB_HD = 1 << 0
	PF_10MB_FD = 1 << 1
	PF_100MB_HD = 1 << 2
	PF_100MB_FD = 1 << 3
	PF_1GB_HD = 1 << 4
	PF_1GB_FD = 1 << 5
	PF_10GB_FD = 1 << 6

	PF_COPPER = 1 << 7
	PF_FIBER = 1 << 8
	PF_AUTONEG = 1 << 9
	PF_PAUSE = 1 << 10
	PF_PAUSE_ASYM = 1 << 11
)
// END: 10 - 5.2.1
