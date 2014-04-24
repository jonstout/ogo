package ofp10

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/jonstout/ogo/protocol/util"
)

// ofp_action_type 1.0
const (
	ActionType_Output = iota
	ActionType_SetVLAN_VID
	ActionType_SetVLAN_PCP
	ActionType_StripVLAN
	ActionType_SetDLSrc
	ActionType_SetDLDst
	ActionType_SetNWSrc
	ActionType_SetNWDst
	ActionType_SetNWTOS
	ActionType_SetTPSrc
	ActionType_SetTPDst
	ActionType_Enqueue
	ActionType_Vendor = 0xffff
)

type Action interface {
	Header() *ActionHeader
	util.Message
}

type ActionHeader struct {
	Type uint16
	Length uint16
}

func (a *ActionHeader) Header() *ActionHeader {
	return a
}

func (a *ActionHeader) Len() (n uint16) {
	return 4
}

func (a *ActionHeader) MarshalBinary() (data []byte, err error) {
	data = make([]byte, a.Len())
	binary.BigEndian.PutUint16(data[:2], a.Type)
	binary.BigEndian.PutUint16(data[2:4], a.Length)
	return
}

func (a *ActionHeader) UnmarshalBinary(data []byte) error {
	if len(data) != 4 {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionHeader message.")
	}
	a.Type = binary.BigEndian.Uint16(data[:2])
	a.Length = binary.BigEndian.Uint16(data[2:4])
	return nil
}

// TODO: Decode other Action types.
func DecodeAction(data []byte) Action {
	t := binary.BigEndian.Uint16(data[:2])
	var a Action
	switch t {
	case ActionType_Output:
		a = new(ActionOutput)
	}
	a.UnmarshalBinary(data)
	return a
}

// Action structure for OFPAT_OUTPUT, which sends packets out ’port’.
// When the ’port’ is the OFPP_CONTROLLER, ’max_len’ indicates the max
// number of bytes to send. A ’max_len’ of zero means no bytes of the
// packet should be sent.
type ActionOutput struct {
	ActionHeader
	Port   uint16
	MaxLen uint16
}

// Returns a new Action Output message which sends packets out
// port number.
func NewActionOutput(number uint16) *ActionOutput {
	act := new(ActionOutput)
	act.Type = ActionType_Output
	act.Length = 8
	act.Port = number
	act.MaxLen = 256
	return act
}

func (a *ActionOutput) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionOutput) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(a.Len()))
	b := make([]byte, 0)
	n := 0

	b, err = a.ActionHeader.MarshalBinary()
	copy(data[n:], b)
	n += len(b)
	binary.BigEndian.PutUint16(data[n:], a.Port)
	n += 2
	binary.BigEndian.PutUint16(data[n:], a.MaxLen)
	n += 2
	return
}

func (a *ActionOutput) UnmarshalBinary(data []byte) error {
	if len(data) < int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionOutput message.")
	}
	n := 0
	err := a.ActionHeader.UnmarshalBinary(data[n:])
	n += int(a.ActionHeader.Len())
	a.Port = binary.BigEndian.Uint16(data[n:])
	n += 2
	a.MaxLen = binary.BigEndian.Uint16(data[n:])
	n += 2
	return err
}

// The enqueue action maps a flow to an already-configured queue, regardless of
// the TOS and VLAN PCP bits. The packet should not change after an enqueue
// action. If the switch needs to set the TOS/PCP bits for internal handling, the
// original values should be restored before sending the packet out.
// A switch may support only queues that are tied to specific PCP/TOS bits.
// In that case, we cannot map an arbitrary flow to a specific queue, therefore the
// action ENQUEUE is not supported. The user can still use these queues and
// map flows to them by setting the relevant fields (TOS, VLAN PCP).
type ActionEnqueue struct {
	ActionHeader
	Port    uint16
	pad     []uint8
	QueueId uint32
}

func NewActionEnqueue(number uint16, queue uint32) *ActionEnqueue {
	a := new(ActionEnqueue)
	a.Type = ActionType_Enqueue
	a.Length = 16
	a.Port = number
	a.pad = make([]uint8, 6)
	a.QueueId = queue
	return a
}

func (a *ActionEnqueue) Len() (n uint16) {
	return a.ActionHeader.Len() + 12
}

func (a *ActionEnqueue) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 12)
	binary.BigEndian.PutUint16(data[:2], a.Port)
	copy(bytes[2:8], a.pad)
	binary.BigEndian.PutUint32(data[8:12], a.QueueId)

	data = append(data, bytes...)
	return
}

func (a *ActionEnqueue) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionEnqueue message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.Port = binary.BigEndian.Uint16(data[4:6])
	copy(a.pad, data[6:12])
	a.QueueId = binary.BigEndian.Uint32(data[12:16])
	return nil
}

// The vlan_vid field is 16 bits long, when an actual VLAN id is only 12 bits.
// The value 0xffff is used to indicate that no VLAN id was set.
type ActionVLANVID struct {
	ActionHeader
	VLANVID uint16
	pad     []uint8
}

// Sets a VLAN ID on tagged packets. VLAN ID may be added to
// untagged packets on some switches.
func NewActionVLANVID(vid uint16) *ActionVLANVID {
	a := new(ActionVLANVID)
	a.Type = ActionType_SetVLAN_VID
	a.Length = 8
	a.VLANVID = vid
	a.pad = make([]byte, 2)
	return a
}

func (a *ActionVLANVID) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionVLANVID) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	binary.BigEndian.PutUint16(data[:2], a.VLANVID)
	copy(bytes[2:4], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionVLANVID) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionVLANVID message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.VLANVID = binary.BigEndian.Uint16(data[4:6])
	copy(a.pad, data[6:8])
	return nil
}

// The vlan_pcp field is 8 bits long, but only the lower 3 bits have meaning.
type ActionVLANPCP struct {
	ActionHeader
	VLANPCP uint8
	pad     []uint8
}

// Modifies PCP on VLAN tagged packets.
func NewActionVLANPCP(pcp uint8) *ActionVLANPCP {
	a := new(ActionVLANPCP)
	a.Type = ActionType_SetVLAN_PCP
	a.Length = 8
	a.VLANPCP = pcp
	a.pad = make([]byte, 3)
	return a
}

func (a *ActionVLANPCP) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionVLANPCP) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	bytes[0] = a.VLANPCP
	copy(bytes[1:4], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionVLANPCP) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionVLANPCP message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.VLANPCP = data[4]
	copy(a.pad, data[5:8])
	return nil
}

// An action_strip_vlan takes no arguments and consists only of a generic
// ofp_action_header. This action strips the VLAN tag if one is present.
type ActionStripVLAN struct {
	ActionHeader
	pad    []uint8
}

// Action to strip VLAN IDs from tagged packets.
func NewActionStripVLAN() *ActionStripVLAN {
	a := new(ActionStripVLAN)
	a.Type = ActionType_StripVLAN
	a.Length = 8
	a.pad = make([]byte, 4)
	return a
}

func (a *ActionStripVLAN) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionStripVLAN) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	copy(bytes[0:4], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionStripVLAN) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionStripVLAN message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	copy(a.pad, data[4:8])
	return nil
}

// The dl_addr field is the MAC address to set.
type ActionDLAddr struct {
	ActionHeader
	DLAddr net.HardwareAddr
	pad    []uint8
}

// Sets the source MAC adddress to dlAddr
func NewActionDLSrc(dlAddr net.HardwareAddr) *ActionDLAddr {
	a := new(ActionDLAddr)
	a.Type = ActionType_SetDLSrc
	a.Length = 16
	a.DLAddr = dlAddr
	a.pad = make([]byte, 6)
	return a
}

// Sets the destination MAC adddress to dlAddr
func NewActionDLDst(dlAddr net.HardwareAddr) *ActionDLAddr {
	a := new(ActionDLAddr)
	a.Type = ActionType_SetDLDst
	a.Length = 16
	a.DLAddr = dlAddr
	a.pad = make([]byte, 6)
	return a
}

func (a *ActionDLAddr) Len() (n uint16) {
	return a.ActionHeader.Len() + 12
}

func (a *ActionDLAddr) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 12)
	copy(bytes[0:6], a.DLAddr)
	copy(bytes[6:12], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionDLAddr) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionDLAddr message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	copy(a.DLAddr, data[4:10])
	copy(a.pad, data[10:16])
	return nil
}

// The nw_addr field is the IP address to set.
type ActionNWAddr struct {
	ActionHeader
	NWAddr net.IP
}

// Sets the source IP adddress to nwAddr
func NewActionNWSrc(nwAddr net.IP) *ActionNWAddr {
	a := new(ActionNWAddr)
	a.Type = ActionType_SetNWSrc
	a.Length = 8
	a.NWAddr = nwAddr
	return a
}

// Sets the destination IP adddress to nwAddr
func NewActionNWDst(nwAddr net.IP) *ActionNWAddr {
	a := new(ActionNWAddr)
	a.Type = ActionType_SetNWDst
	a.Length = 8
	a.NWAddr = nwAddr
	return a
}

func (a *ActionNWAddr) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionNWAddr) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	copy(bytes[:4], a.NWAddr)

	data = append(data, bytes...)
	return
}

func (a *ActionNWAddr) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionDLAddr message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	copy(a.NWAddr, data[4:8])
	return nil
}

// The nw_tos field is the 6 upper bits of the ToS field to set, in the original bit
// positions (shifted to the left by 2).
type ActionNWTOS struct {
	ActionHeader
	NWTOS  uint8
	pad    []uint8
}

// Set ToS field in IP packets.
func NewActionNWTOS(tos uint8) *ActionNWTOS {
	a := new(ActionNWTOS)
	a.Type = ActionType_SetNWTOS
	a.Length = 8
	a.NWTOS = tos
	a.pad = make([]byte, 3)
	return a
}

func (a *ActionNWTOS) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionNWTOS) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	data[0] = a.NWTOS
	copy(bytes[1:4], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionNWTOS) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionDLAddr message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.NWTOS = data[4]
	copy(a.pad, data[5:8])
	return nil
}

// The tp_port field is the TCP/UDP/other port to set.
type ActionTPPort struct {
	ActionHeader
	TPPort uint16
	pad    []uint8
}

// Returns an action that sets the transport layer source port.
func NewActionTPSrc(port uint16) *ActionTPPort {
	a := new(ActionTPPort)
	a.Type = ActionType_SetTPSrc
	a.Length = 8
	a.TPPort = port
	a.pad = make([]byte, 2)
	return a
}

// Returns an action that sets the transport layer destination
// port.
func NewActionTPDst(port uint16) *ActionTPPort {
	a := new(ActionTPPort)
	a.Type = ActionType_SetTPDst
	a.Length = 8
	a.TPPort = port
	a.pad = make([]byte, 2)
	return a
}

func (a *ActionTPPort) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionTPPort) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	binary.BigEndian.PutUint16(data[:2], a.TPPort)
	copy(bytes[2:4], a.pad)

	data = append(data, bytes...)
	return
}

func (a *ActionTPPort) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionNWTOS message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.TPPort = binary.BigEndian.Uint16(data[4:6])
	copy(a.pad, data[6:8])
	return nil
}

// The Vendor field is the Vendor ID, which takes the same form as in struct
// ofp_vendor.
type ActionVendor struct {
	ActionHeader
	Vendor uint32
}

func NewActionVendor(vendor uint32) *ActionVendor {
	a := new(ActionVendor)
	a.Type = ActionType_Vendor
	a.Length = 8
	a.Vendor = vendor
	return a
}

func (a *ActionVendor) Len() (n uint16) {
	return a.ActionHeader.Len() + 4
}

func (a *ActionVendor) MarshalBinary() (data []byte, err error) {
	data, err = a.ActionHeader.MarshalBinary()

	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(data[:4], a.Vendor)

	data = append(data, bytes...)
	return
}

func (a *ActionVendor) UnmarshalBinary(data []byte) error {
	if len(data) != int(a.Len()) {
		return errors.New("The []byte the wrong size to unmarshal an " +
			"ActionVendor message.")
	}
	a.ActionHeader.UnmarshalBinary(data[:4])
	a.Vendor = binary.BigEndian.Uint32(data[4:8])
	return nil
}
