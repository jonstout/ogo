package ofp10

import (
	"encoding/binary"

	"github.com/jonstout/ogo/protocol/ofpxx"
)

// ofp_flow_mod
type FlowMod struct {
	ofpxx.Header
	Match Match
	Cookie uint64

	Command     uint16
	IdleTimeout uint16
	HardTimeout uint16
	Priority    uint16
	BufferId    uint32
	OutPort     uint16
	Flags       uint16
	Actions     []Action
}

func NewFlowMod() *FlowMod {
	f := new(FlowMod)
	f.Header = ofpxx.NewOfp10Header()
	f.Header.Type = Type_FlowMod
	f.Match = *NewMatch()
	// Add a generator for f.Cookie here
	f.Cookie = 0

	f.Command = FC_ADD
	f.IdleTimeout = 0
	f.HardTimeout = 0
	// Add a priority gen here
	f.Priority = 1000
	f.BufferId = 0xffffffff
	f.OutPort = P_NONE
	f.Flags = 0
	f.Actions = make([]Action, 0)
	return f
}

func (f *FlowMod) AddAction(a Action) {
	f.Actions = append(f.Actions, a)
}

func (f *FlowMod) Len() (n uint16) {
	n = 72
	if f.Command == FC_DELETE || f.Command == FC_DELETE_STRICT {
		return
	}
	for _, v := range f.Actions {
		n += v.Len()
	}
	return
}

func (f *FlowMod) MarshalBinary() (data []byte, err error) {
	f.Header.Length = f.Len()
	data, err = f.Header.MarshalBinary()
	bytes, err := f.Match.MarshalBinary()
	data = append(data, bytes...)

	bytes = make([]byte, 24)
	n := 0
	binary.BigEndian.PutUint64(bytes[n:], f.Cookie)
	n += 8
	binary.BigEndian.PutUint16(bytes[n:], f.Command)
	n += 2
	binary.BigEndian.PutUint16(bytes[n:], f.IdleTimeout)
	n += 2
	binary.BigEndian.PutUint16(bytes[n:], f.HardTimeout)
	n += 2
	binary.BigEndian.PutUint16(bytes[n:], f.Priority)
	n += 2
	binary.BigEndian.PutUint32(bytes[n:], f.BufferId)
	n += 2
	binary.BigEndian.PutUint16(bytes[n:], f.OutPort)
	n += 2
	binary.BigEndian.PutUint16(bytes[n:], f.Flags)
	n += 2
	data = append(data, bytes...)

	for _, a := range f.Actions {
		bytes, err = a.MarshalBinary()
		data = append(data, bytes...)
	}
	return
}

func (f *FlowMod) UnmarshalBinary(data []byte) error {
	n := 0
	f.Header.UnmarshalBinary(data[n:])
	n += int(f.Header.Len())
	f.Match.UnmarshalBinary(data[n:])
	n += int(f.Match.Len())
	f.Cookie = binary.BigEndian.Uint64(data[n:])
	n += 8
	f.Command = binary.BigEndian.Uint16(data[n:])
	n += 2
	f.IdleTimeout = binary.BigEndian.Uint16(data[n:])
	n += 2
	f.HardTimeout = binary.BigEndian.Uint16(data[n:])
	n += 2
	f.Priority = binary.BigEndian.Uint16(data[n:])
	n += 2
	f.BufferId = binary.BigEndian.Uint32(data[n:])
	n += 4
	f.OutPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	f.Flags = binary.BigEndian.Uint16(data[n:])
	n += 2

	for n < int(f.Header.Length) {
		a := DecodeAction(data[n:])
		f.Actions = append(f.Actions, a)
		n += int(a.Len())
	}
	return nil
}

// ofp_flow_mod_command 1.0
const (
	FC_ADD = iota // OFPFC_ADD == 0
	FC_MODIFY
	FC_MODIFY_STRICT
	FC_DELETE
	FC_DELETE_STRICT
)

// ofp_flow_mod_flags 1.0
const (
	FF_SEND_FLOW_REM = 1 << 0
	FF_CHECK_OVERLAP = 1 << 1
	FF_EMERG         = 1 << 2
)

// BEGIN: ofp10 - 5.4.2
type FlowRemoved struct {
	ofpxx.Header
	Match Match
	Cookie   uint64
	Priority uint16
	Reason   uint8
	pad      []uint8 // Size 1

	DurationSec  uint32
	DurationNSec uint32

	IdleTimeout uint16
	pad2        []uint8 // Size 2
	PacketCount uint64
	ByteCount   uint64
}

func NewFlowRemoved() *FlowRemoved {
	f := new(FlowRemoved)
	f.Header = ofpxx.NewOfp10Header()
	f.Match = *NewMatch()
	f.pad = make([]byte, 1)
	f.pad2 = make([]byte, 2)
	return f
}

func (f *FlowRemoved) Len() (n uint16) {
	n = f.Header.Len()
	n += f.Match.Len()
	n += 42
	return
}

func (f *FlowRemoved) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(f.Len()))
	bytes := make([]byte, 0)
	next := 0

	bytes, err = f.Header.MarshalBinary()
	copy(data[next:], bytes)
	next += int(f.Header.Len())
	bytes, err = f.Match.MarshalBinary()
	copy(data[next:], bytes)
	next += int(f.Match.Len())
	binary.BigEndian.PutUint64(data[next:], f.Cookie)
	next += 8
	binary.BigEndian.PutUint16(data[next:], f.Priority)
	next += 2
	data[next] = f.Reason
	next += 1
	copy(data[next:], f.pad)
	next += len(f.pad)
	binary.BigEndian.PutUint32(data[next:], f.DurationSec)
	next += 4
	binary.BigEndian.PutUint32(data[next:], f.DurationNSec)
	next += 4
	binary.BigEndian.PutUint16(data[next:], f.IdleTimeout)
	next += 2
	copy(data[next:], f.pad2)
	next += len(f.pad2)
	binary.BigEndian.PutUint64(data[next:], f.PacketCount)
	next += 8
	binary.BigEndian.PutUint64(data[next:], f.ByteCount)
	next += 8
	return
}

func (f *FlowRemoved) UnmarshalBinary(data []byte) error {
	next := 0
	var err error
	err = f.Header.UnmarshalBinary(data[next:])
	next += int(f.Header.Len())
	err = f.Match.UnmarshalBinary(data[next:])
	next += int(f.Match.Len())
	f.Cookie = binary.BigEndian.Uint64(data[next:])
	next += 8
	f.Priority = binary.BigEndian.Uint16(data[next:])
	next += 2
	f.Reason = data[next]
	next += 1
	copy(f.pad, data[next:])
	next += len(f.pad)
	f.DurationSec = binary.BigEndian.Uint32(data[next:])
	next += 4
	f.DurationNSec = binary.BigEndian.Uint32(data[next:])
	next += 4
	f.IdleTimeout = binary.BigEndian.Uint16(data[next:])
	next += 2
	copy(f.pad2, data[next:])
	next += len(f.pad2)
	f.PacketCount = binary.BigEndian.Uint64(data[next:])
	next += 8
	f.ByteCount = binary.BigEndian.Uint64(data[next:])
	next += 8
	return err
}

// ofp_flow_removed_reason 1.0
const (
	RR_IDLE_TIMEOUT = iota
	RR_HARD_TIMEOUT
	RR_DELETE
)
