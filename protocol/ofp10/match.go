package ofp10

import (
	"encoding/binary"
	"net"
)

// ofp_match 1.0
type Match struct {
	Wildcards uint32           /* Wildcard fields. */
	InPort    uint16           /* Input switch port. */
	DLSrc     net.HardwareAddr //[ETH_ALEN]uint8 /* Ethernet source address. */
	DLDst     net.HardwareAddr //[ETH_ALEN]uint8 /* Ethernet destination address. */
	DLVLAN    uint16           /* Input VLAN id. */
	DLVLANPcp uint8            /* Input VLAN priority. */
	pad       []uint8          /* Align to 64-bits Size 1 */
	DLType    uint16           /* Ethernet frame type. */
	NWTos     uint8            /* IP ToS (actually DSCP field, 6 bits). */
	NWProto   uint8            /* IP protocol or lower 8 bits of ARP opcode. */
	pad2      []uint8          /* Align to 64-bits Size 2 */
	NWSrc     net.IP           /* IP source address. */
	NWDst     net.IP           /* IP destination address. */
	TPSrc     uint16           /* TCP/UDP source port. */
	TPDst     uint16           /* TCP/UDP destination port. */
}

func NewMatch() *Match {
	m := new(Match)
	// By default wildcard all fields
	m.Wildcards = FW_ALL
	m.DLSrc = make([]byte, ETH_ALEN)
	m.DLDst = make([]byte, ETH_ALEN)
	m.NWSrc = make([]byte, 4)
	m.NWDst = make([]byte, 4)
	m.pad = make([]byte, 1)
	m.pad2 = make([]byte, 2)
	return m
}

func (m *Match) Len() (n uint16) {
	return 40
}

func (m *Match) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(m.Len()))
	n := 0
	binary.BigEndian.PutUint32(data[n:], m.Wildcards)
	n += 4
	binary.BigEndian.PutUint16(data[n:], m.InPort)
	n += 2
	copy(data[n:], m.DLSrc)
	n += len(m.DLSrc)
	copy(data[n:], m.DLDst)
	n += len(m.DLDst)
	binary.BigEndian.PutUint16(data[n:], m.DLVLAN)
	n += 2
	data[n] = m.DLVLANPcp
	n += 1
	copy(data[n:], m.pad)
	n += len(m.pad)
	binary.BigEndian.PutUint16(data[n:], m.DLType)
	n += 2
	data[n] = m.NWTos
	n += 1
	data[n] = m.NWProto
	n += 1
	copy(data[n:], m.pad2)
	n += len(m.pad2)
	copy(data[n:], m.NWSrc)
	n += len(m.NWSrc)
	copy(data[n:], m.NWDst)
	n += len(m.NWDst)
	binary.BigEndian.PutUint16(data[n:], m.TPSrc)
	n += 2
	binary.BigEndian.PutUint16(data[n:], m.TPDst)
	n += 2
	return
}

func (m *Match) UnmarshalBinary(data []byte) error {
	// Any non-zero value fields should not be wildcarded.
	if m.InPort != 0 {
		m.Wildcards = m.Wildcards ^ FW_IN_PORT
	}
	if m.DLSrc.String() != "00:00:00:00:00:00" {
		m.Wildcards = m.Wildcards ^ FW_DL_SRC
	}
	if m.DLDst.String() != "00:00:00:00:00:00" {
		m.Wildcards = m.Wildcards ^ FW_DL_DST
	}
	if m.DLVLAN != 0 {
		m.Wildcards = m.Wildcards ^ FW_DL_VLAN
	}
	if m.DLVLANPcp != 0 {
		m.Wildcards = m.Wildcards ^ FW_DL_VLAN_PCP
	}
	if m.DLType != 0 {
		m.Wildcards = m.Wildcards ^ FW_DL_TYPE
	}
	if m.NWTos != 0 {
		m.Wildcards = m.Wildcards ^ FW_NW_TOS
	}
	if m.NWProto != 0 {
		m.Wildcards = m.Wildcards ^ FW_NW_PROTO
	}
	if m.NWSrc.String() != "0.0.0.0" {
		m.Wildcards = m.Wildcards ^ FW_NW_SRC_MASK
	}
	if m.NWDst.String() != "0.0.0.0" {
		m.Wildcards = m.Wildcards ^ FW_NW_DST_MASK
	}
	if m.TPSrc != 0 {
		m.Wildcards = m.Wildcards ^ FW_TP_SRC
	}
	if m.TPDst != 0 {
		m.Wildcards = m.Wildcards ^ FW_TP_DST
	}

	n := 0
	m.Wildcards = binary.BigEndian.Uint32(data[n:])
	n += 4
	m.InPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(m.DLSrc, data[n:])
	n += len(m.DLSrc)
	copy(m.DLDst, data[n:])
	n += len(m.DLDst)
	m.DLVLAN = binary.BigEndian.Uint16(data[n:])
	n += 2
	m.DLVLANPcp = data[n]
	n += 1
	copy(m.pad, data[n:])
	n += len(m.pad)
	m.DLType = binary.BigEndian.Uint16(data[n:])
	n += 2
	m.NWTos = data[n]
	n += 1
	m.NWProto = data[n]
	n += 1
	copy(m.pad2, data[n:])
	n += len(m.pad2)
	copy(m.NWSrc, data[n:])
	n += len(m.NWSrc)
	copy(m.NWDst, data[n:])
	n += len(m.NWDst)
	m.TPSrc = binary.BigEndian.Uint16(data[n:])
	n += 2
	m.TPDst = binary.BigEndian.Uint16(data[n:])
	n += 2
	return nil
}

// ofp_flow_wildcards 1.0
const (
	FW_IN_PORT  = 1 << 0
	FW_DL_VLAN  = 1 << 1
	FW_DL_SRC   = 1 << 2
	FW_DL_DST   = 1 << 3
	FW_DL_TYPE  = 1 << 4
	FW_NW_PROTO = 1 << 5
	FW_TP_SRC   = 1 << 6
	FW_TP_DST   = 1 << 7

	FW_NW_SRC_SHIFT = 8
	FW_NW_SRC_BITS  = 6
	FW_NW_SRC_MASK  = ((1 << FW_NW_SRC_BITS) - 1) << FW_NW_SRC_SHIFT
	FW_NW_SRC_ALL   = 32 << FW_NW_SRC_SHIFT

	FW_NW_DST_SHIFT = 14
	FW_NW_DST_BITS  = 6
	FW_NW_DST_MASK  = ((1 << FW_NW_DST_BITS) - 1) << FW_NW_DST_SHIFT
	FW_NW_DST_ALL   = 32 << FW_NW_DST_SHIFT

	FW_DL_VLAN_PCP = 1 << 20
	FW_NW_TOS      = 1 << 21

	FW_ALL = ((1 << 22) - 1)
)
