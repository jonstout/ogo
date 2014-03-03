package ofp10

import (
	"encoding/binary"

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
	n += uint16(s.Body.Len())
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
	MfrDesc   []byte // Size DESC_STR_LEN
	HWDesc    []byte // Size DESC_STR_LEN
	SWDesc    []byte // Size DESC_STR_LEN
	SerialNum []byte // Size SERIAL_NUM_LEN
	DPDesc    []byte // Size DESC_STR_LEN
}

func NewDescStats() *DescStats {
	s := new(DescStats)
	s.MfrDesc = make([]byte, DESC_STR_LEN)
	s.HWDesc = make([]byte, DESC_STR_LEN)
	s.SWDesc = make([]byte, DESC_STR_LEN)
	s.SerialNum = make([]byte, SERIAL_NUM_LEN)
	s.DPDesc = make([]byte, DESC_STR_LEN)
	return s
}

func (s *DescStats) Len() (n uint16) {
	return uint16(DESC_STR_LEN * 4 + SERIAL_NUM_LEN)
}

func (s *DescStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	copy(data[n:], s.MfrDesc)
	n += len(s.MfrDesc)
	copy(data[n:], s.HWDesc)
	n += len(s.HWDesc)
	copy(data[n:], s.SWDesc)
	n += len(s.SWDesc)
	copy(data[n:], s.SerialNum)
	n += len(s.SerialNum)
	copy(data[n:], s.DPDesc)
	n += len(s.DPDesc)
	return
}

func (s *DescStats) UnmarshalBinary(data []byte) error {
	n := 0
	copy(s.MfrDesc, data[n:])
	n += len(s.MfrDesc)
	copy(s.HWDesc, data[n:])
	n += len(s.HWDesc)
	copy(s.SWDesc, data[n:])
	n += len(s.SWDesc)
	copy(s.SerialNum, data[n:])
	n += len(s.SerialNum)
	copy(s.DPDesc, data[n:])
	n += len(s.DPDesc)
	return nil
}

const (
	DESC_STR_LEN   = 256
	SERIAL_NUM_LEN = 32
)

// ofp_flow_stats_request 1.0
type FlowStatsRequest struct {
	Match
	TableId uint8
	pad     uint8
	OutPort uint16
}

func NewFlowStatsRequest() *FlowStatsRequest {
	s := new(FlowStatsRequest)
	s.Match = *NewMatch()
	return s
}

func (s *FlowStatsRequest) Len() (n uint16) {
	return s.Match.Len() + 4
}

func (s *FlowStatsRequest) MarshalBinary() (data []byte, err error) {
	data, err = s.Match.MarshalBinary()
	
	b := make([]byte, 4)
	n := 0
	b[n] = s.TableId
	n += 1
	b[n] = s.pad
	n += 1
	binary.BigEndian.PutUint16(data[n:], s.OutPort)
	n += 2
	data = append(data, b...)
	return
}

func (s *FlowStatsRequest) UnmarshalBinary(data []byte) error {
	err := s.Match.UnmarshalBinary(data)
	n := s.Match.Len()

	s.TableId = data[n]
	n += 1
	s.pad = data[n]
	n += 1
	s.OutPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	return err
}

// ofp_flow_stats 1.0
type FlowStats struct {
	Length       uint16
	TableId      uint8
	pad          uint8
	Match
	DurationSec  uint32
	DurationNSec uint32
	Priority     uint16
	IdleTimeout  uint16
	HardTimeout  uint16
	pad2         []uint8 // Size 6
	Cookie       uint64
	PacketCount  uint64
	ByteCount    uint64
	Actions      []Action
}

func NewFlowStats() *FlowStats {
	f := new(FlowStats)
	f.Match = *NewMatch()
	f.pad2 = make([]byte, 6)
	return f
}

func (s *FlowStats) Len() (n uint16) {
	n = 24 + s.Match.Len() + 24
	for _, a := range s.Actions {
		n += a.Len()
	}
	return
}

func (s *FlowStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0

	binary.BigEndian.PutUint16(data[n:], s.Length)
	n += 2
	data[n] = s.TableId
	n += 1
	data[n] = s.pad
	n += 1
	b, err := s.Match.MarshalBinary()
	data = append(data, b...)
	n += len(b)
	binary.BigEndian.PutUint32(data[n:], s.DurationSec)
	n += 4
	binary.BigEndian.PutUint32(data[n:], s.DurationNSec)
	n += 4
	binary.BigEndian.PutUint16(data[n:], s.Priority)
	n += 2
	binary.BigEndian.PutUint16(data[n:], s.IdleTimeout)
	n += 2
	binary.BigEndian.PutUint16(data[n:], s.HardTimeout)
	n += 2
	copy(data[n:], s.pad2)
	n += len(s.pad2)
	binary.BigEndian.PutUint64(data[n:], s.Cookie)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.PacketCount)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.ByteCount)
	n += 8

	for _, a := range s.Actions {
		b, err = a.MarshalBinary()
		data = append(data, b...)
		n += len(b)
	}
	return
}

func (s *FlowStats) UnmarshalBinary(data []byte) error {
	n := 0
	s.Length = binary.BigEndian.Uint16(data[n:])
	n += 2
	s.TableId = data[n]
	n += 1
	s.pad = data[n]
	n += 1	
	err := s.Match.UnmarshalBinary(data[n:])
	n += int(s.Match.Len())
	s.DurationSec = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.DurationNSec = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.Priority = binary.BigEndian.Uint16(data[n:])
	n += 2
	s.IdleTimeout = binary.BigEndian.Uint16(data[n:])
	n += 2
	s.HardTimeout = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(s.pad2, data[n:])
	n += len(s.pad2)
	s.Cookie = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.PacketCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.ByteCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	for n < int(s.Length) {
		t := binary.BigEndian.Uint16(data[n:])
		var a Action
		switch t {
		case ActionType_Output:
			a = NewActionOutput(0)
		case ActionType_SetVLAN_VID:
			a = NewActionVLANVID(0xffff)
		case ActionType_SetVLAN_PCP:
			a = NewActionVLANPCP(0)
		case ActionType_StripVLAN:
			a = NewActionStripVLAN()
		case ActionType_SetDLSrc:
			a = NewActionDLSrc(make([]byte, 6))
		case ActionType_SetDLDst:
			a = NewActionDLDst(make([]byte, 6))
		case ActionType_SetNWSrc:
			a = NewActionNWSrc(make([]byte, 4))
		case ActionType_SetNWDst:
			a = NewActionNWDst(make([]byte, 4))
		case ActionType_SetNWTOS:
			a = NewActionNWTOS(0)
		case ActionType_SetTPSrc:
			a = NewActionTPSrc(0)
		case ActionType_SetTPDst:
			a = NewActionTPDst(0)
		case ActionType_Enqueue:
			a = NewActionEnqueue(0, 0)
		case ActionType_Vendor:
			a = NewActionVendor(0)
		}
		s.Actions = append(s.Actions, a)
		n += int(a.Len())
	}
	return err
}

// ofp_aggregate_stats_request 1.0
type AggregateStatsRequest struct {
	Match
	TableId uint8
	pad     uint8
	OutPort uint16
}

func NewAggregateStatsRequest() *AggregateStatsRequest {
	return new(AggregateStatsRequest)
}

func (s *AggregateStatsRequest) Len() (n uint16) {
	return s.Match.Len() + 4
}

func (s *AggregateStatsRequest) MarshalBinary() (data []byte, err error) {
	data, err = s.Match.MarshalBinary()

	b := make([]byte, 4)
	n := 0
	b[n] = s.TableId
	n += 1
	b[n] = s.pad
	n += 1
	binary.BigEndian.PutUint16(data[n:], s.OutPort)
	n += 2
	data = append(data, b...)
	return
}

func (s *AggregateStatsRequest) UnmarshalBinary(data []byte) error {
	n := 0
	s.Match.UnmarshalBinary(data[n:])
	n += int(s.Match.Len())
	s.TableId = data[n]
	n += 1
	s.pad = data[n]
	n += 1
	s.OutPort = binary.BigEndian.Uint16(data[n:])
	n += 2
	return nil
}

// ofp_aggregate_stats_reply 1.0
type AggregateStats struct {
	PacketCount uint64
	ByteCount   uint64
	FlowCount   uint32
	pad         []uint8 // Size 4
}

func NewAggregateStats() *AggregateStats {
	s := new(AggregateStats)
	s.pad = make([]byte, 4)
	return s
}

func (s *AggregateStats) Len() (n uint16) {
	return 24
}

func (s *AggregateStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	binary.BigEndian.PutUint64(data[n:], s.PacketCount)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.ByteCount)
	n += 8
	binary.BigEndian.PutUint32(data[n:], s.FlowCount)
	n += 4
	copy(data[n:], s.pad)
	n += 4
	return
}

func (s *AggregateStats) UnmarshalBinary(data []byte) error {
	n := 0
	s.PacketCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.ByteCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.FlowCount = binary.BigEndian.Uint32(data[n:])
	n += 4
	copy(s.pad, data[n:])
	return nil
}

// ofp_table_stats 1.0
type TableStats struct {
	TableId      uint8
	pad          []uint8 // Size 3
	Name         []byte // Size MAX_TABLE_NAME_LEN
	Wildcards    uint32
	MaxEntries   uint32
	ActiveCount  uint32
	LookupCount  uint64
	MatchedCount uint64
}

func NewTableStats() *TableStats {
	s := new(TableStats)
	s.pad = make([]byte, 3)
	s.Name = make([]byte, MAX_TABLE_NAME_LEN)
	return s
}

func (s *TableStats) Len() (n uint16) {
	return 4 + MAX_TABLE_NAME_LEN + 28
}

func (s *TableStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	data[n] = s.TableId
	n += 1
	copy(data[n:], s.pad)
	n += len(s.pad)
	copy(data[n:], s.Name)
	n += len(s.Name)
	binary.BigEndian.PutUint32(data[n:], s.Wildcards)
	n += 4
	binary.BigEndian.PutUint32(data[n:], s.MaxEntries)
	n += 4
	binary.BigEndian.PutUint32(data[n:], s.ActiveCount)
	n += 4
	binary.BigEndian.PutUint64(data[n:], s.LookupCount)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.MatchedCount)
	n += 8
	return
}

func (s *TableStats) UnmarshalBinary(data []byte) error {
	n := 0
	s.TableId = data[0]
	n += 1
	copy(s.pad, data[n:])
	n += len(s.pad)
	copy(s.Name, data[n:])
	n += len(s.Name)
	s.Wildcards = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.MaxEntries = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.ActiveCount = binary.BigEndian.Uint32(data[n:])
	n += 4
	s.LookupCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.MatchedCount = binary.BigEndian.Uint64(data[n:])
	n += 8
	return nil
}

const (
	MAX_TABLE_NAME_LEN = 32
)

// ofp_port_stats_request 1.0
type PortStatsRequest struct {
	PortNo uint16
	pad    []uint8 // Size 6
}

func NewPortStatsRequest() *PortStatsRequest {
	p := new(PortStatsRequest)
	p.pad = make([]byte, 6)
	return p
}

func (s *PortStatsRequest) Len() (n uint16) {
	return 8
}

func (s *PortStatsRequest) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	binary.BigEndian.PutUint16(data[n:], s.PortNo)
	n += 2
	copy(data[n:], s.pad)
	n += len(s.pad)
	return
}

func (s *PortStatsRequest) UnmarshalBinary(data []byte) error {
	n := 0
	s.PortNo = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(s.pad, data[n:])
	n += len(s.pad)
	return nil
}

// ofp_port_stats 1.0
type PortStats struct {
	PortNo     uint16
	pad        []uint8 // Size 6
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

func NewPortStats() *PortStats {
	p := new(PortStats)
	p.pad = make([]byte, 6)
	return p
}

func (s *PortStats) Len() (n uint16) {
	return 104
}

func (s *PortStats) MarshalBinary() (data []byte, err error) {
	data = make([]byte, int(s.Len()))
	n := 0
	binary.BigEndian.PutUint16(data[n:], s.PortNo)
	n += 2
	copy(data[n:], s.pad)
	n += len(s.pad)
	binary.BigEndian.PutUint64(data[n:], s.RxPackets)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxPackets)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxBytes)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxBytes)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxDropped)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxDropped)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxErrors)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.TxErrors)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxFrameErr)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxOverErr)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.RxCRCErr)
	n += 8
	binary.BigEndian.PutUint64(data[n:], s.Collisions)
	n += 8
	return
}

func (s *PortStats) UnmarshalBinary(data []byte) error {
	n := 0
	s.PortNo = binary.BigEndian.Uint16(data[n:])
	n += 2
	copy(s.pad, data[n:])
	n += len(s.pad)
	s.RxPackets = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxPackets = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxBytes = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxBytes = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxDropped = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxDropped = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxErrors = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.TxErrors = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxFrameErr = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxOverErr = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.RxCRCErr = binary.BigEndian.Uint64(data[n:])
	n += 8
	s.Collisions = binary.BigEndian.Uint64(data[n:])
	n += 8
	return nil
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
	return q
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
	n := 0
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
	return nil
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
	return p
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

	b, err = s.Desc.MarshalBinary()
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
