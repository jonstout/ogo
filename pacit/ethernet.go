package pacit

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

// see http://en.wikipedia.org/wiki/EtherType
const (
	IPv4_MSG = 0x0800
	ARP_MSG  = 0x0806
	LLDP_MSG = 0x88cc
	WOL_MSG  = 0x0842
	RARP_MSG = 0x8035
	VLAN_MSG = 0x8100
)

type ReadWriterMeasurer interface {
	io.Reader
	io.ReaderFrom
	Len() uint16
}

type ReadWriteMeasurer interface {
	io.ReadWriter
	Len() uint16
}

type Ethernet struct {
	Delimiter uint8
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLANID    VLAN
	Ethertype uint16
	Data      ReadWriteMeasurer
}

func NewEthernet() *Ethernet {
	eth := new(Ethernet)
	eth.HWDst = net.HardwareAddr(make([]byte, 6))
	eth.HWSrc = net.HardwareAddr(make([]byte, 6))
	eth.VLANID = *NewVLAN()
	eth.Ethertype = 0x800
	return eth
}

func (e *Ethernet) Len() (n uint16) {
	if e.VLANID.VID != 0 {
		n += 5
	}
	n += 12
	n += 2
	if e.Data != nil {
		n += e.Data.Len()
	}
	return
}

func (e *Ethernet) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	//If you send a packet with the delimiter to the wire
	//packets are incorrectly interpreted.
	binary.Write(buf, binary.BigEndian, e.HWDst)
	binary.Write(buf, binary.BigEndian, e.HWSrc)

	if e.VLANID.VID != 0 {
		c := []byte{0, 0}
		e.VLANID.Read(c)
		binary.Write(buf, binary.BigEndian, c)
	}
	binary.Write(buf, binary.BigEndian, e.Ethertype)
	// In case the data type isn't known
	if e.Data != nil {
		if n, err := buf.ReadFrom(e.Data); n == 0 {
			return int(n), err
		}
	}
	n, err = buf.Read(b)
	return n, io.EOF
}

type PacitBuffer struct{ *bytes.Buffer }

func NewBuffer(buf []byte) *PacitBuffer {
	b := &PacitBuffer{Buffer: bytes.NewBuffer(buf)}
	return b
}

func (b *PacitBuffer) Len() uint16 {
	return uint16(b.Buffer.Len())
}

func (e *Ethernet) Write(b []byte) (n int, err error) {
	// Delimiter comes in from the wire. Not sure why this is the case, ignore
	n += 1
	e.HWDst = make([]byte, 6)
	e.HWDst = b[1:7]
	n += 6
	e.HWSrc = make([]byte, 6)
	e.HWSrc = b[7:13]
	n += 6
	e.Ethertype = binary.BigEndian.Uint16(b[13:15])
	n += 2
	// If tagged
	if e.Ethertype == VLAN_MSG {
		e.VLANID = *new(VLAN)
		e.VLANID.Write(b[13:17])
		n += 2
		e.Ethertype = binary.BigEndian.Uint16(b[13:15])
		n += 2
	} else {
		e.VLANID = *new(VLAN)
		e.VLANID.VID = 0
	}

	switch e.Ethertype {
	case IPv4_MSG:
		e.Data = new(IPv4)
		m, _ := e.Data.Write(b[n:])
		n += m
	case ARP_MSG:
		e.Data = new(ARP)
		m, _ := e.Data.Write(b[n:])
		n += m
	case RARP_MSG:
	/*case LLDP_MSG:
	e.Data = new(LLDP)
	m, _ := e.Data.Write(b[n:])
	n += m*/
	default:
		e.Data = NewBuffer(b[n:])
		n += int(e.Data.Len())
		//n += len(b[n:])
	}
	return
}

const (
	PCP_MASK = 0xe000
	DEI_MASK = 0x1000
	VID_MASK = 0x0fff
)

type VLAN struct {
	TPID uint16
	PCP  uint8
	DEI  uint8
	VID  uint8
}

func NewVLAN() *VLAN {
	v := new(VLAN)
	v.TPID = 0x8100
	v.VID = 0
	return v
}

func (v *VLAN) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v.TPID)
	var tci uint16 = 0
	tci = (tci | uint16(v.PCP)<<13) + (tci | uint16(v.DEI)<<12) + (tci | uint16(v.VID))
	binary.Write(buf, binary.BigEndian, tci)
	n, err = buf.Read(b)
	return
}

func (v *VLAN) ReadFrom(r io.Reader) (n int64, err error) {
	var tci uint16 = 0
	if err = binary.Read(r, binary.BigEndian, &tci); err != nil {
		return
	}
	n += 2
	v.PCP = uint8(PCP_MASK & tci >> 13)
	v.DEI = uint8(DEI_MASK & tci >> 12)
	v.VID = uint8(VID_MASK & tci)
	return
}

func (v *VLAN) Write(b []byte) (n int, err error) {
	var tci uint16 = 0
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &v.TPID); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &tci); err != nil {
		return
	}
	n += 2
	v.PCP = uint8(PCP_MASK & tci >> 13)
	v.DEI = uint8(DEI_MASK & tci >> 12)
	v.VID = uint8(VID_MASK & tci)
	return
}
