package ofp10

import (
	"io"
	"net"
	"encoding/binary"
	"bytes"
)
// ofp_match 1.0
type OfpMatch struct {
	Wildcards uint32 /* Wildcard fields. */
	InPort uint16 /* Input switch port. */
	DLSrc net.HardwareAddr//[OFP_ETH_ALEN]uint8 /* Ethernet source address. */
	DLDst net.HardwareAddr//[OFP_ETH_ALEN]uint8 /* Ethernet destination address. */
	DLVLAN uint16 /* Input VLAN id. */
	DLVLANPcp uint8 /* Input VLAN priority. */
	Pad [1]uint8 /* Align to 64-bits */
	DLType uint16 /* Ethernet frame type. */
	NWTos uint8 /* IP ToS (actually DSCP field, 6 bits). */
	NWProto uint8 /* IP protocol or lower 8 bits of ARP opcode. */
	Pad1 [2]uint8 /* Align to 64-bits */
	NWSrc net.IP /* IP source address. */
	NWDst net.IP /* IP destination address. */
	TPSrc uint16 /* TCP/UDP source port. */
	TPDst uint16 /* TCP/UDP destination port. */
}

func NewMatch() *OfpMatch {
	m := new(OfpMatch)
	m.Wildcards = 0xffffffff
	m.DLSrc = make([]byte, OFP_ETH_ALEN)
	m.DLDst = make([]byte, OFP_ETH_ALEN)
	m.NWSrc = make([]byte, 4)
	m.NWDst = make([]byte, 4)
	return m
}

func (m *OfpMatch) Len() (n uint16) {
	return 40
}

func (m *OfpMatch) Read(b []byte) (n int, err error) {
	// Any non-zero value fields should not be wildcarded.
	if m.InPort != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_IN_PORT
	}
	if m.DLSrc.String() != "00:00:00:00:00:00" {
		m.Wildcards = m.Wildcards ^ OFPFW_DL_SRC
	}
	if m.DLDst.String() != "00:00:00:00:00:00" {
		m.Wildcards = m.Wildcards ^ OFPFW_DL_DST
	}
	if m.DLVLAN != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_DL_VLAN
	}
	if m.DLVLANPcp != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_DL_VLAN_PCP
	}
	if m.DLType != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_DL_TYPE
	}
	if m.NWTos != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_NW_TOS
	}
	if m.NWProto != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_NW_PROTO
	}
	if m.NWSrc.String() != "0.0.0.0" {
		m.Wildcards = m.Wildcards ^ OFPFW_NW_SRC_ALL
	}
	if m.NWDst.String() != "0.0.0.0" {
		m.Wildcards = m.Wildcards ^ OFPFW_NW_DST_ALL
	}
	if m.TPSrc != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_TP_SRC
	}
	if m.TPDst != 0 {
		m.Wildcards = m.Wildcards ^ OFPFW_TP_DST
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

func (m *OfpMatch) Write(b []byte) (n int, err error) {
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
	n += OFP_ETH_ALEN
	if err = binary.Read(buf, binary.BigEndian, &m.DLDst); err != nil {
		return
	}
	n += OFP_ETH_ALEN
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
	OFPFW_IN_PORT = 1 << 0
	OFPFW_DL_VLAN = 1 << 1
	OFPFW_DL_SRC = 1 << 2
	OFPFW_DL_DST = 1 << 3
	OFPFW_DL_TYPE = 1 << 4
	OFPFW_NW_PROTO = 1 << 5
	OFPFW_TP_SRC = 1 << 6
	OFPFW_TP_DST = 1 << 7

	OFPFW_NW_SRC_SHIFT = 8
	OFPFW_NW_SRC_BITS = 6
	OFPFW_NW_SRC_MASK = ((1 << OFPFW_NW_SRC_BITS) - 1) << OFPFW_NW_SRC_SHIFT
	OFPFW_NW_SRC_ALL = 32 << OFPFW_NW_SRC_SHIFT

	OFPFW_NW_DST_SHIFT = 14
	OFPFW_NW_DST_BITS = 6
	OFPFW_NW_DST_MASK = ((1 << OFPFW_NW_DST_BITS) - 1) << OFPFW_NW_DST_SHIFT
	OFPFW_NW_DST_ALL = 32 << OFPFW_NW_DST_SHIFT

	OFPFW_DL_VLAN_PCP = 1 << 20
	OFPFW_NW_TOS = 1 << 21

	OFPFW_ALL = ((1 << 22) - 1)
)
