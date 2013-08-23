package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)


// ofp_action_type 1.0
const (
	AT_OUTPUT = iota
	AT_SET_VLAN_VID
	AT_SET_VLAN_PCP
	AT_STRIP_VLAN
	AT_SET_DL_SRC
	AT_SET_DL_DST
	AT_SET_NW_SRC
	AT_SET_NW_DST
	AT_SET_NW_TOS
	AT_SET_TP_SRC
	AT_SET_TP_DST
	AT_ENQUEUE
	AT_VENDOR = 0xffff
)


type Action interface {
	Packetish
	ActionType() uint16
}


// ofp_action_output 1.0
// Action structure for OFPAT_OUTPUT, which sends packets out ’port’.
// When the ’port’ is the OFPP_CONTROLLER, ’max_len’ indicates the max
// number of bytes to send. A ’max_len’ of zero means no bytes of the
// packet should be sent.
type ActionOutput struct {
	Type   uint16
	Length uint16
	Port   uint16
	MaxLen uint16
}


// Returns a new Action Output message which sends packets out
// port number.
func NewActionOutput(number uint16) *ActionOutput {
	act := new(ActionOutput)
	act.Type = AT_OUTPUT
	act.Length = 8
	act.Port = number
	act.MaxLen = 0
	return act
}


func (a *ActionOutput) ActionType() (n uint16) {
	return AT_OUTPUT
}


func (a *ActionOutput) Len() (n uint16) {
	return 8
}


func (a *ActionOutput) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Write(buf, binary.BigEndian, a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Write(buf, binary.BigEndian, a.Port); err != nil {
		return
	}
	n += 2
	if err = binary.Write(buf, binary.BigEndian, a.MaxLen); err != nil {
		return
	}
	n += 2
	if n, err = buf.Read(b); n == 0 {
		return
	}
	return n, io.EOF
}


func (a *ActionOutput) Write(b []byte) (n int, err error) {
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
type ActionEnqueue struct {
	Type    uint16
	Length  uint16
	Port    uint16
	Pad     [6]uint8
	QueueId uint32
}

func NewActionEnqueue() *ActionEnqueue {
	a := new(ActionEnqueue)
	a.Type = AT_ENQUEUE
	a.Length = 16
	return a
}

func (a *ActionEnqueue) Len() (n uint16) {
	return a.Length
}

func (a *ActionEnqueue) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionEnqueue) Write(b []byte) (n int, err error) {
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
type ActionVLANVID struct {
	Type    uint16
	Length  uint16
	VLANVID uint16
	Pad     [2]uint8
}

func NewActionVLANVID() *ActionVLANVID {
	a := new(ActionVLANVID)
	a.Type = AT_SET_VLAN_VID
	a.Length = 8
	a.VLANVID = 0xffff
	return a
}

func (a *ActionVLANVID) Len() (n uint16) {
	return a.Length
}

func (a *ActionVLANVID) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionVLANVID) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.VLANVID); err != nil {
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
type ActionVLANPCP struct {
	Type    uint16
	Length  uint16
	VLANPCP uint8
	Pad     [3]uint8
}

func NewActionVLANPCP() *ActionVLANPCP {
	a := new(ActionVLANPCP)
	a.Type = AT_SET_VLAN_PCP
	a.Length = 8
	return a
}

func (a *ActionVLANPCP) Len() (n uint16) {
	return 8
}

func (a *ActionVLANPCP) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionVLANPCP) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.VLANPCP); err != nil {
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
type ActionDLAddr struct {
	Type   uint16
	Length uint16
	DLAddr net.HardwareAddr
	Pad    [6]uint8
}

func NewActionDLSrc() *ActionDLAddr {
	a := new(ActionDLAddr)
	a.Type = AT_SET_DL_SRC
	a.Length = 16
	a.DLAddr = make([]byte, ETH_ALEN)
	return a
}

func NewActionDLDst() *ActionDLAddr {
	a := new(ActionDLAddr)
	a.Type = AT_SET_DL_DST
	a.Length = 16
	a.DLAddr = make([]byte, ETH_ALEN)
	return a
}

func (a *ActionDLAddr) Len() (n uint16) {
	return a.Length
}

func (a *ActionDLAddr) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a.Type)
	err = binary.Write(buf, binary.BigEndian, a.Length)
	err = binary.Write(buf, binary.BigEndian, a.DLAddr)
	err = binary.Write(buf, binary.BigEndian, a.Pad)
	n, err = buf.Read(b)
	return
}

func (a *ActionDLAddr) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	a.DLAddr = make([]byte, ETH_ALEN)
	if err = binary.Read(buf, binary.BigEndian, &a.DLAddr); err != nil {
		return
	}
	n += ETH_ALEN
	if err = binary.Read(buf, binary.BigEndian, &a.Pad); err != nil {
		return
	}
	n += 6
	return
}

// ofp_action_nw_addr 1.0
type ActionNWAddr struct {
	Type   uint16
	Length uint16
	NWAddr net.IP
}

func NewActionNWSrc() *ActionNWAddr {
	a := new(ActionNWAddr)
	a.Type = AT_SET_NW_SRC
	a.Length = 8
	a.NWAddr = make([]byte, 4)
	return a
}

func NewActionNWDst() *ActionNWAddr {
	a := new(ActionNWAddr)
	a.Type = AT_SET_NW_DST
	a.Length = 8
	a.NWAddr = make([]byte, 4)
	return a
}

func (a *ActionNWAddr) Len() (n uint16) {
	return 8
}

func (a *ActionNWAddr) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a.Type)
	err = binary.Write(buf, binary.BigEndian, a.Length)
	err = binary.Write(buf, binary.BigEndian, a.NWAddr)
	n, err = buf.Read(b)
	return
}

func (a *ActionNWAddr) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	if err = binary.Read(buf, binary.BigEndian, &a.Type); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &a.Length); err != nil {
		return
	}
	n += 2
	a.NWAddr = make([]byte, 4)
	if err = binary.Read(buf, binary.BigEndian, &a.NWAddr); err != nil {
		return
	}
	n += 4
	return
}

// ofp_action_nw_tos 1.0
type ActionNWTOS struct {
	Type   uint16
	Length uint16
	NWTOS  uint8
	Pad    [3]uint8
}

func NewActionNWTOS() *ActionNWTOS {
	a := new(ActionNWTOS)
	a.Type = AT_SET_NW_TOS
	a.Length = 8
	return a
}

func (a *ActionNWTOS) Len() (n uint16) {
	return a.Length
}

func (a *ActionNWTOS) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionNWTOS) Write(b []byte) (n int, err error) {
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
type ActionTPPort struct {
	Type   uint16
	Length uint16
	TPPort uint16
	Pad    [2]uint8
}

func NewActionTPSrc() *ActionTPPort {
	a := new(ActionTPPort)
	a.Type = AT_SET_TP_SRC
	a.Length = 8
	return a
}

func NewActionTPDst() *ActionTPPort {
	a := new(ActionTPPort)
	a.Type = AT_SET_TP_DST
	a.Length = 8
	return a
}

func (a *ActionTPPort) Len() (n uint16) {
	return a.Length
}

func (a *ActionTPPort) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionTPPort) Write(b []byte) (n int, err error) {
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
type ActionVendorPort struct {
	Type   uint16
	Length uint16
	Vendor uint32
}

func NewActionVendorPort() *ActionVendorPort {
	a := new(ActionVendorPort)
	a.Type = AT_VENDOR
	a.Length = 8
	return a
}

func (a *ActionVendorPort) Len() (n uint16) {
	return a.Length
}

func (a *ActionVendorPort) Read(b []byte) (n int, err error) {
	a.Length = a.Len()
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, a)
	n, err = buf.Read(b)
	return
}

func (a *ActionVendorPort) Write(b []byte) (n int, err error) {
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
