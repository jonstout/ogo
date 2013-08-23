package ofp10

import (
	"io"
	"bytes"
	"encoding/binary"
)

// ofp_flow_mod
type FlowMod struct {
	Header Header
	Match Match
	Cookie uint64

	Command uint16
	IdleTimeout uint16
	HardTimeout uint16
	Priority uint16
	BufferID uint32
	OutPort uint16
	Flags uint16
	Actions []Action
}

func NewFlowMod() *FlowMod {
	f := new(FlowMod)
	f.Header = *NewHeader()
	f.Header.Type = T_FLOW_MOD
	f.Match = *NewMatch()
	// Add a generator for f.Cookie here
	f.Cookie = 0

	f.Command = FC_ADD
	f.IdleTimeout = 0
	f.HardTimeout = 0
	// Add a priority gen here
	f.Priority = 1000
	f.BufferID = 0xffffffff
	f.OutPort = P_NONE
	f.Flags = 0
	f.Actions = make([]Action, 0)
	return f
}


func (f *FlowMod) AddAction(a Action) {
	f.Actions = append(f.Actions, a)
}


func (f *FlowMod) GetHeader() *Header {
	return &f.Header
}

func (f *FlowMod) Len() (n uint16) {
	for _, v := range f.Actions {
		n += v.Len()
	}
	n += 72
	return
}

func (f *FlowMod) Read(b []byte) (n int, err error) {
	f.Header.Length = f.Len()
	buf := new(bytes.Buffer)
	buf.ReadFrom(&f.Header)
	buf.ReadFrom(&f.Match)
	binary.Write(buf, binary.BigEndian, f.Cookie)
	binary.Write(buf, binary.BigEndian, f.Command)
	binary.Write(buf, binary.BigEndian, f.IdleTimeout)
	binary.Write(buf, binary.BigEndian, f.HardTimeout)
	binary.Write(buf, binary.BigEndian, f.Priority)
	binary.Write(buf, binary.BigEndian, f.BufferID)
	binary.Write(buf, binary.BigEndian, f.OutPort)
	binary.Write(buf, binary.BigEndian, f.Flags)
	for _, a := range f.Actions {
		buf.ReadFrom(a)
	}
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (f *FlowMod) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = f.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	m := 0
	m, err = f.Match.Write(buf.Next(40))
	if m == 0 {
		return
	}
	n += m
	if err = binary.Read(buf, binary.BigEndian, &f.Cookie); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.Command); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.IdleTimeout); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.HardTimeout); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.Priority); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.BufferID); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &f.OutPort); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.Flags); err != nil {
		return
	}
	n += 2
	actionCount := buf.Len() / 8
	f.Actions = make([]Action, actionCount)
	for i := 0; i < actionCount; i++ {
		a := new(ActionOutput)
		m, err = a.Write(buf.Next(8))
		if m == 0 {
			return
		}
		n += m
		f.Actions[i] = a
	}
	return
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
	FF_EMERG = 1 << 2
)

// BEGIN: ofp10 - 5.4.2
type FlowRemoved struct {
	Header Header
	Match Match
	Cookie uint64
	Priority uint16
	Reason uint8
	Pad [1]uint8

	DurationSec uint32
	DurationNSec uint32

	IdleTimeout uint16
	Pad2 [2]uint8
	PacketCount uint64
	ByteCount uint64
}

func (f *FlowRemoved) GetHeader() *Header {
	return &f.Header
}

func (f *FlowRemoved) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, f)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (f *FlowRemoved) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = f.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	m := 0
	m, err = f.Match.Write(buf.Next(40))
	if m == 0 {
		return
	}
	n += m
	if err = binary.Read(buf, binary.BigEndian, &f.Cookie); err != nil {
		return
	}
	n += 8
	if err = binary.Read(buf, binary.BigEndian, &f.Priority); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.Reason); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &f.Pad); err != nil {
		return
	}
	n += 1
	if err = binary.Read(buf, binary.BigEndian, &f.DurationSec); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &f.DurationNSec); err != nil {
		return
	}
	n += 4
	if err = binary.Read(buf, binary.BigEndian, &f.IdleTimeout); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.Pad2); err != nil {
		return
	}
	n += 2
	if err = binary.Read(buf, binary.BigEndian, &f.PacketCount); err != nil {
		return
	}
	n += 8
	if err = binary.Read(buf, binary.BigEndian, &f.ByteCount); err != nil {
		return
	}
	n += 8
	return
}

// ofp_flow_removed_reason 1.0
const (
	RR_IDLE_TIMEOUT = iota
	RR_HARD_TIMEOUT
	RR_DELETE
)
