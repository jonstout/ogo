package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
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
	Pad       [1]uint8         /* Align to 64-bits */
	DLType    uint16           /* Ethernet frame type. */
	NWTos     uint8            /* IP ToS (actually DSCP field, 6 bits). */
	NWProto   uint8            /* IP protocol or lower 8 bits of ARP opcode. */
	Pad1      [2]uint8         /* Align to 64-bits */
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
	return m
}

func (m *Match) Len() (n uint16) {
	return 40
}

func (m *Match) Read(b []byte) (n int, err error) {
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.Wildcards)
	binary.Write(buf, binary.BigEndian, m.InPort)
	binary.Write(buf, binary.BigEndian, m.DLSrc)
	binary.Write(buf, binary.BigEndian, m.DLDst)
	binary.Write(buf, binary.BigEndian, m.DLVLAN)
	binary.Write(buf, binary.BigEndian, m.DLVLANPcp)
	binary.Write(buf, binary.BigEndian, m.Pad)
	binary.Write(buf, binary.BigEndian, m.DLType)
	binary.Write(buf, binary.BigEndian, m.NWTos)
	binary.Write(buf, binary.BigEndian, m.NWProto)
	binary.Write(buf, binary.BigEndian, m.Pad1)
	binary.Write(buf, binary.BigEndian, m.NWSrc)
	binary.Write(buf, binary.BigEndian, m.NWDst)
	binary.Write(buf, binary.BigEndian, m.TPSrc)
	binary.Write(buf, binary.BigEndian, m.TPDst)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (m *Match) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &m.Wildcards); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &m.InPort); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &m.DLSrc); err != nil {
		return
	}
	n += ETH_ALEN
	if err = binary.Read(buf, binary.BigEndian, &m.DLDst); err != nil {
		return
	}
	n += ETH_ALEN
	if err = binary.Read(buf, binary.BigEndian, &m.DLVLAN); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &m.DLVLANPcp); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &m.Pad); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &m.DLType); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &m.NWTos); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &m.NWProto); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &m.Pad1); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &m.NWSrc); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &m.NWDst); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &m.TPSrc); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &m.TPDst); err != nil {
		return
	}
	n += 2
	return
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
