package eth

import (
	"encoding/binary"
	"errors"
	"net"
	
	"github.com/jonstout/ogo/protocol/arp"
	"github.com/jonstout/ogo/protocol/ipv4"
	"github.com/jonstout/ogo/protocol/util"
)

// see http://en.wikipedia.org/wiki/EtherType
const (
	IPv4_MSG = 0x0800
	ARP_MSG  = 0x0806
	LLDP_MSG = 0x88cc
	WOL_MSG  = 0x0842
	RARP_MSG = 0x8035
	VLAN_MSG = 0x8100

	IPv6_MSG     = 0x86DD
	STP_MSG      = 0x4242
	STP_BPDU_MSG = 0xAAAA
)

type Ethernet struct {
	Delimiter uint8
	HWDst     net.HardwareAddr
	HWSrc     net.HardwareAddr
	VLANID    VLAN
	Ethertype uint16
	Data      util.Message
}

func New() *Ethernet {
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

func (e *Ethernet) MarshalBinary() (data []byte, err error) {
	data = make([]byte, e.Len() - e.Data.Len())
	n := 0
	copy(data[:n+len(e.HWDst)], e.HWDst)
	n += len(e.HWDst)
	copy(data[n:n+len(e.HWSrc)], e.HWSrc)
	n += len(e.HWSrc)

	bytes, err := e.VLANID.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[n:n+len(bytes)], bytes)
	n += len(bytes)

	binary.BigEndian.PutUint16(data[n:n+2], e.Ethertype)
	n += 2

	bytes, err = e.Data.MarshalBinary()
	if err != nil {
		return
	}
	copy(data[n:n+len(bytes)], bytes)
	return
}

func (e *Ethernet) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return errors.New("The []byte is too short to unmarshal a full Ethernet message.")
	}
	n := 0
	e.HWDst = data[:n+len(e.HWDst)]
	n += len(e.HWDst)

	e.HWSrc = data[n:n+len(e.HWSrc)]
	n += len(e.HWSrc)

	e.Ethertype = binary.BigEndian.Uint16(data[n:n+2])
	n += 2
	if e.Ethertype == VLAN_MSG {
		e.VLANID = *new(VLAN)
		err := e.VLANID.UnmarshalBinary(data[n:n+5])
		if err != nil {
			return err
		}
		n += 5

		e.Ethertype = binary.BigEndian.Uint16(data[n:n+2])
		n += 2
	} else {
		e.VLANID = *new(VLAN)
		e.VLANID.VID = 0
	}

	switch e.Ethertype {
	case IPv4_MSG:
		e.Data = new(ipv4.IPv4)
	case ARP_MSG:
		e.Data = new(arp.ARP)
	default:
		e.Data = new(util.Buffer)
	}
	err := e.Data.UnmarshalBinary(data[n:])
	return err
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

func (v *VLAN) Len() (n uint16) {
	return 4
}

func (v *VLAN) MarshalBinary() (data []byte, err error) {
	data = make([]byte, v.Len())
	binary.BigEndian.PutUint16(data[:2], v.TPID)
	var tci uint16
	tci = (tci | uint16(v.PCP)<<13) + (tci | uint16(v.DEI)<<12) + (tci | uint16(v.VID))
	binary.BigEndian.PutUint16(data[2:], tci)
	return
}

func (v *VLAN) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return errors.New("The []byte is too short to unmarshal a full VLAN header.")
	}
	v.TPID = binary.BigEndian.Uint16(data[:2])
	var tci uint16
	tci = binary.BigEndian.Uint16(data[2:])
	v.PCP = uint8(PCP_MASK & tci >> 13)
	v.DEI = uint8(DEI_MASK & tci >> 12)
	v.VID = uint8(VID_MASK & tci)
	return nil
}
