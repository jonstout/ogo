package ofp10

import (
	"io"
	"bytes"
	"encoding/binary"
)

// ofp_action_type 1.0
const (
	OFPAT_OUTPUT = iota
	OFPAT_SET_VLAN_VID
	OFPAT_SET_VLAN_PCP
	OFPAT_STRIP_VLAN
	OFPAT_SET_DL_SRC
	OFPAT_SET_DL_DST
	OFPAT_SET_NW_SRC
	OFPAT_SET_NW_DST
	OFPAT_SET_NW_TOS
	OFPAT_SET_TP_SRC
	OFPAT_SET_TP_DST
	OFPAT_ENQUEUE
	OFPAT_VENDOR = 0xffff
)

// ofp_action_header 1.0
type OfpActionHeader struct {
	Type uint16
	Length uint16
	Pad [4]uint8
}

func (a *OfpActionHeader) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (a *OfpActionHeader) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return 
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return 
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 4
	return
}

// ofp_action_output 1.0
type OfpActionOutput struct {
	Type uint16
	Length uint16
	Port uint16
	MaxLen uint16
}

func NewActionOutput() *OfpActionOutput {
	act := new(OfpActionOutput)
	act.Type = OFPAT_OUTPUT
	act.Port = OFPP_FLOOD
	act.Length = 8
	return act
}

func (a *OfpActionOutput) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionOutput) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Port); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.MaxLen); err != nil {
		return
	}
	n += 2
	return
}

// ofp_action_enqueue 1.0
type OfpActionEnqueue struct {
	Type uint16
	Length uint16
	Port uint16
	Pad [6]uint8
	QueueId uint32
}

func (a *OfpActionEnqueue) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionEnqueue) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Port); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 6
	if err = binary.Read(buf, binary.BigEndian, &a.QueueId); err != nil {
		return
	}
	n += 4
	return
}

// ofp_action_vlan_vid 1.0
type OfpActionVlanVid struct {
	Type uint16
	Length uint16
	VlanVid uint16
	Pad [2]uint8
}

func (a *OfpActionVlanVid) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionVlanVid) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.VlanVid); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 2
	return
}

// ofp_action_vlan_pcp 1.0
type OfpActionVlanPcp struct {
	Type uint16
	Length uint16
	VlanPcp uint8
	Pad [3]uint8
}

func (a *OfpActionVlanPcp) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionVlanPcp) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.VlanPcp); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 3
	return
}

// ofp_action_dl_addr 1.0
type OfpActionDLAddr struct {
	Type uint16
	Length uint16
	DLAddr [OFP_ETH_ALEN]uint8
	Pad [6]uint8
}

func (a *OfpActionDLAddr) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionDLAddr) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.DLAddr); err != nil {
		return
	}
	n += OFP_ETH_ALEN
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 6
	return
}

// ofp_action_nw_addr 1.0
type OfpActionNWAddr struct {
	Type uint16
	Length uint16
	NWAddr uint32
}

func (a *OfpActionNWAddr) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionNWAddr) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.NWAddr); err != nil {
		return
	}
	n += 4
	return
}

// ofp_action_nw_tos 1.0
type OfpActionNWTOS struct {
	Type uint16
	Length uint16
	NWTOS uint8
	Pad [3]uint8
}

func (a *OfpActionNWTOS) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionNWTOS) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.NWTOS); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 3
	return
}

// ofp_action_tp_port 1.0
type OfpActionTPPort struct {
	Type uint16
	Length uint16
	TPPort uint16
	Pad [8]uint8
}

func (a *OfpActionTPPort) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionTPPort) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.TPPort); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 8
	return
}

// ofp_action_vendor_header 1.0
type OfpActionVendorPort struct {
	Type uint16
	Length uint16
	Vendor uint32
}

func (a *OfpActionVendorPort) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *OfpActionVendorPort) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Vendor); err != nil {
		return
	}
	n += 4
	return
}
