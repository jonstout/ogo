package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/jonstout/ogo/protocol/ofpxx"
	"github.com/jonstout/ogo/protocol/util"
)

// ofp_stats_request 1.0
type StatsRequest struct {
	ofpxx.Header
	Type   uint16
	Flags  uint16
	Body   util.Message
}

func (s *StatsRequest) Len() (n uint16) {
	return s.Header.Len() + 4 + s.Body.Len()
}

func (s *StatsRequest) MarshalBinary() (data []byte, err error) {
	data, err = s.Header.MarshalBinary()

	b := make([]byte, 4)
	n := 0
	binary.BigEndian.PutUint16(b[n:], s.Type)
	n += 2
	binary.BigEndian.PutUint16(b[n:], s.Flags)
	n += 2
	data = append(data, b...)

	b, err = s.Body.MarshalBinary()
	data = append(data, b...)
	return
}

func (s *StatsRequest) UnmarshalBinary(data []byte) error {
	err := s.Header.UnmarshalBinary(data)
	n := s.Header.Len()

	s.Type = binary.BigEndian.Uint16(data[n:])
	n += 2
	s.Flags = binary.BigEndian.Uint16(data[n:])
	n += 2

	var req util.Message
	switch s.Type {
	case StatsType_Aggregate:
		req = s.Body.(*AggregateStatsRequest)
		err = req.UnmarshalBinary(data[n:])
	case StatsType_Desc:
		break
	case StatsType_Flow:
		req = s.Body.(*FlowStatsRequest)
		err = req.UnmarshalBinary(data[n:])
	case StatsType_Port:
		req = s.Body.(*PortStatsRequest)
		err = req.UnmarshalBinary(data[n:])
	case StatsType_Table:
		break
	case StatsType_Queue:
		req = s.Body.(*QueueStatsRequest)
		err = req.UnmarshalBinary(data[n:])
	case StatsType_Vendor:
		break
	}
	return err
}

// _stats_reply 1.0
type StatsReply struct {
	ofpxx.Header
	Type   uint16
	Flags  uint16
	Body   util.Message
}

func (s *StatsReply) Len() (n uint16) {
	n = s.Header.Len()
	n += 4
	n += uint16(len(s.Body))
	return
}

func (s *StatsReply) MarshalBinary() (data []byte, err error) {
	data, err = s.Header.MarshalBinary()

	b := make([]byte, 4)
	n := 0
	binary.BigEndian.PutUint16(b[n:], s.Type)
	n += 2
	binary.BigEndian.PutUint16(b[n:], s.Flags)
	n += 2
	data = append(data, b...)

	b, err = s.Body.MarshalBinary()
	data = append(data, b...)
	return
}

func (s *StatsReply) UnmarshalBinary(data []byte) error {
	err := s.Header.UnmarshalBinary(data)
	n := s.Header.Len()

	s.Type = binary.BigEndian.Uint16(data[n:])
	n += 2
	s.Flags = binary.BigEndian.Uint16(data[n:])
	n += 2

	var req util.Message
	switch s.Type {
	case StatsType_Aggregate:
		req = s.Body.(*AggregateStats)
	case StatsType_Desc:
		req = s.Body.(*DescStats)
	case StatsType_Flow:
		// Array
		req = s.Body.(*FlowStats)
	case StatsType_Port:
		req = s.Body.(*PortStats)
	case StatsType_Table:
		// Array
		req = s.Body.(*TableStats)
	case StatsType_Queue:
		// Array
		req = s.Body.(*QueueStats)
	case StatsType_Vendor:
		// Array of Group Stats
		break
	}
	err = req.UnmarshalBinary(data[n:])
	return err
}

// _stats_types
const (
	/* Description of this OpenFlow switch.
	* The request body is empty.
	* The reply body is struct ofp_desc_stats. */
	StatsType_Desc = iota
	/* Individual flow statistics.
	* The request body is struct ofp_flow_stats_request.
	* The reply body is an array of struct ofp_flow_stats. */
	StatsType_Flow
	/* Aggregate flow statistics.
	* The request body is struct ofp_aggregate_stats_request.
	* The reply body is struct ofp_aggregate_stats_reply. */
	StatsType_Aggregate
	/* Flow table statistics.
	* The request body is empty.
	* The reply body is an array of struct ofp_table_stats. */
	StatsType_Table
	/* Port statistics.
	* The request body is struct ofp_port_stats_request.
	* The reply body is an array of struct ofp_port_stats. */
	StatsType_Port
	/* Queue statistics for a port
	* The request body is struct _queue_stats_request.
	* The reply body is an array of struct ofp_queue_stats */
	StatsType_Queue
	/* Group counter statistics.
	* The request body is struct ofp_group_stats_request.
	* The reply is an array of struct ofp_group_stats. */
	StatsType_Vendor = 0xffff
)

// ofp_desc_stats 1.0
type DescStats struct {
	MfrDesc   [DESC_STR_LEN]byte
	HWDesc    [DESC_STR_LEN]byte
	SWDesc    [DESC_STR_LEN]byte
	SerialNum [SERIAL_NUM_LEN]byte
	DPDesc    [DESC_STR_LEN]byte
}

func (s *DescStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *DescStats) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.MfrDesc)
	if err != nil {
		return
	}
	n += 256
	err = binary.Read(buf, binary.BigEndian, &s.HWDesc)
	if err != nil {
		return
	}
	n += 256
	err = binary.Read(buf, binary.BigEndian, &s.SWDesc)
	if err != nil {
		return
	}
	n += 256
	err = binary.Read(buf, binary.BigEndian, &s.SerialNum)
	if err != nil {
		return
	}
	n += 32
	err = binary.Read(buf, binary.BigEndian, &s.DPDesc)
	if err != nil {
		return
	}
	n += 256
	return
}

const (
	DESC_STR_LEN   = 256
	SERIAL_NUM_LEN = 32
)

// ofp_flow_stats_request 1.0
type FlowStatsRequest struct {
	Match   Match
	TableID uint8
	Pad     uint8
	OutPort uint16
}

func (s *FlowStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *FlowStatsRequest) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = s.Match.Write(buf.Bytes())
	if n == 0 {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &s.TableID)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.OutPort)
	if err != nil {
		return
	}
	n += 2
	return
}

// ofp_flow_stats 1.0
type FlowStats struct {
	Length       uint16
	TableID      uint8
	Pad          uint8
	Match        Match
	DurationSec  uint32
	DurationNSec uint32
	Priority     uint16
	IdleTimeout  uint16
	HardTimeout  uint16
	Pad2         [6]uint8
	Cookie       uint64
	PacketCount  uint64
	ByteCount    uint64
	Actions      []Action
}

func (s *FlowStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *FlowStats) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.Length)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.TableID)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 1
	m := 0
	m, err = s.Match.Write(buf.Next(40))
	if m == 0 {
		return
	}
	n += m
	err = binary.Read(buf, binary.BigEndian, &s.DurationSec)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.DurationNSec)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.Priority)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.IdleTimeout)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.HardTimeout)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.Pad2)
	if err != nil {
		return
	}
	n += 6
	err = binary.Read(buf, binary.BigEndian, &s.Cookie)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.PacketCount)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.ByteCount)
	if err != nil {
		return
	}
	n += 8

	for buf.Len() > 2 {
		t := binary.BigEndian.Uint16(buf.Bytes()[:2])
		l := binary.BigEndian.Uint16(buf.Bytes()[2:4])
		var a Action
		m := 0

		switch t {
		case AT_OUTPUT:
			a = NewActionOutput(0)
		case AT_SET_VLAN_VID:
			a = NewActionVLANVID(0xffff)
		case AT_SET_VLAN_PCP:
			a = NewActionVLANPCP(0)
		case AT_STRIP_VLAN:
			a = NewActionStripVLAN()
		case AT_SET_DL_SRC:
			a = NewActionDLSrc(make([]byte, 6))
		case AT_SET_DL_DST:
			a = NewActionDLDst(make([]byte, 6))
		case AT_SET_NW_SRC:
			a = NewActionNWSrc(make([]byte, 4))
		case AT_SET_NW_DST:
			a = NewActionNWDst(make([]byte, 4))
		case AT_SET_NW_TOS:
			a = NewActionNWTOS(0)
		case AT_SET_TP_SRC:
			a = NewActionTPSrc(0)
		case AT_SET_TP_DST:
			a = NewActionTPDst(0)
		case AT_ENQUEUE:
			a = NewActionEnqueue(0, 0)
		case AT_VENDOR:
			a = NewActionVendor(0)
		}

		if m, err = a.Write(buf.Next(int(l))); m == 0 {
			return
		} else {
			n += m
		}
		s.Actions = append(s.Actions, a)
	}
	return
}

// ofp_aggregate_stats_request 1.0
type AggregateStatsRequest struct {
	Match   Match
	TableID uint8
	Pad     uint8
	OutPort uint16
}

func (s *AggregateStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *AggregateStatsRequest) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = s.Match.Write(buf.Next(40))
	err = binary.Read(buf, binary.BigEndian, &s.TableID)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.OutPort)
	if err != nil {
		return
	}
	n += 2
	return
}

// ofp_aggregate_stats_reply 1.0
type AggregateStatsReply struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
	Pad         []uint8 // Size 4
}

func (s *AggregateStatsReply) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *AggregateStatsReply) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.PacketCount)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.ByteCount)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.FlowCount)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 4
	return
}

// ofp_table_stats 1.0
type TableStats struct {
	TableID      uint8
	Pad          []uint8 // Size 3
	Name         []byte // Size MAX_TABLE_NAME_LEN
	Wildcards    uint32
	MaxEntries   uint32
	ActiveCount  uint32
	LookupCount  uint64
	MatchedCount uint64
}

func (s *TableStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *TableStats) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.TableID)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 3
	err = binary.Read(buf, binary.BigEndian, &s.Name)
	if err != nil {
		return
	}
	n += 32
	err = binary.Read(buf, binary.BigEndian, &s.Wildcards)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.MaxEntries)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.ActiveCount)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.LookupCount)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.MatchedCount)
	if err != nil {
		return
	}
	n += 8
	return
}

const (
	MAX_TABLE_NAME_LEN = 32
)

// ofp_port_stats_request 1.0
type PortStatsRequest struct {
	PortNo uint16
	Pad    []uint8 // Size 6
}

func (s *PortStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *PortStatsRequest) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.PortNo)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 6
	return
}

// ofp_port_stats 1.0
type PortStats struct {
	PortNo     uint16
	Pad        []uint8 // Size 6
	RxPackets  uint64
	TxPackets  uint64
	RxBytes    uint64
	TxBytes    uint64
	RxDropped  uint64
	TxDropped  uint64
	RxErrors   uint64
	TxErrors   uint64
	RxFrameErr uint64
	RxOverErr  uint64
	RxCRCErr   uint64
	Collisions uint64
}

func (p *PortStats) Len() (n uint16) {
	n = 104
	return
}

func (s *PortStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *PortStats) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.PortNo)
	if err != nil {
		return
	}
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.RxPackets)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxPackets)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxBytes)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxBytes)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxPackets)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxPackets)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxDropped)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxDropped)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxErrors)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxErrors)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxFrameErr)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxOverErr)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.RxCRCErr)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.Collisions)
	if err != nil {
		return
	}
	n += 8
	return
}

// ofp_queue_stats_request 1.0
type QueueStatsRequest struct {
	PortNo  uint16
	pad     []uint8 // Size 2
	QueueId uint32
}

func NewQueueStatsRequest() *QueueStatsRequest {
	q := new(QueueStatsRequest)
	q.pad = make([]byte, 2)
}

func (s *QueueStatsRequest) Len() (n uint16) {
	return 8
}

func (s *QueueStatsRequest) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	binary.BigEndian.PutUint16(data[n:], s.PortNo)
	n += 2
	copy(data[n:], s.pad)
	n += 2
	binary.BigEndian.PutUint32(data[n:], s.QueueId)
	n += 4
	return
}

func (s *QueueStatsRequest) UnmarshalBinary(data []byte) error {
	n := 0
	s.PortNo = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(s.pad, data[n:])
	n += 2
	s.QueueId = binary.BigEndian.Uint32(data[n:])
	return nil
}

// ofp_queue_stats 1.0
type QueueStats struct {
	PortNo    uint16
	pad       []uint8 // Size 2
	QueueId   uint32
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
}

func (s *QueueStats) Len() (n uint16) {
	return 32
}

func (s *QueueStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 32)
	n := 0

	binary.BigEndian.PutUint16(data[n:], s.PortNo)
	n += 2
	copy(data[n:], s.pad)
	n += 2
	binary.BigEndian.PutUint32(data[n:], s.QueueId)
	n += 4
	binary.BigEndian.PutUint64(data[n:], s.TxBytes)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxPackets)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxErrors)
	n += 8
	return
}

func (s *QueueStats) UnmarshalBinary(data []byte) error {
	err := s.Header.UnmarshalBinary(data)
	n := s.Header.Len()

	s.PortNo = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(s.pad, data[n:])
	n += len(s.pad)
	s.QueueId = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.TxBytes = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxPackets = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxErrors = binary.BigEndian.Uint64(data[n:])
	n += 8
	return err
}

// ofp_port_status
type PortStatus struct {
	ofpxx.Header
	Reason uint8
	pad    []uint8 // Size 7
	Desc   PhyPort
}

func NewPortStatus() *PortStatus {
	p := new(PortStatus)
	p.Header = ofpxx.NewOfp10Header()
	p.pad = make([]byte, 7)
}

func (p *PortStatus) Len() (n uint16) {
	n = p.Header.Len()
	n += 8
	n += p.Desc.Len()
	return
}

func (s *PortStatus) MarshalBinary() (data []byte, err error) {
	s.Header.Length = s.Len()
	data, err = s.Header.MarshalBinary()

	b := make([]byte, 8)
	n := 0
	b[0] = s.Reason
	n += 1
	copy(b[n:], s.pad)
	data = append(data, b...)

	b, err = s.PhyPort.MarshalBinary()
	data = append(data, b...)
	return
}

func (s *PortStatus) UnmarshalBinary(data []byte) error {
	err := s.Header.UnmarshalBinary(data)
	n := int(s.Header.Len())
	
	s.Reason = data[n]
	n += 1
	copy(s.pad, data[n:])
	n += len(s.pad)

	err = s.Desc.UnmarshalBinary(data[n:])
	return err
}

// ofp_port_reason 1.0
const (
	PR_ADD = iota
	PR_DELETE
	PR_MODIFY
)
