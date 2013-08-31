package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ofp_stats_request 1.0
type StatsRequest struct {
	Header Header
	Type   uint16
	Flags  uint16
	Body   interface{}
}

func (s *StatsRequest) GetHeader() *Header {
	return &s.Header
}

func (s *StatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (s *StatsRequest) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = s.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &s.Type)
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.Flags)
	n += 2
	m := 0
	switch s.Type {
	case ST_AGGREGATE:
		// empty
		break
	case ST_DESC:
		// empty
		break
	case ST_FLOW:
		// _aggregate_stats_request
		a := s.Body.(*AggregateStatsRequest)
		m, err = a.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case ST_PORT:
		// _port_stats_request
		p := s.Body.(*PortStatsRequest)
		m, err = p.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case ST_TABLE:
		// empty
		break
	case ST_QUEUE:
		// ofp_queue_stats_request
		q := s.Body.(*QueueStatsRequest)
		m, err = q.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case ST_VENDOR:
		break
	}
	return
}

// _stats_reply 1.0
type StatsReply struct {
	Header Header
	Type   uint16
	Flags  uint16
	Body   []uint8
}

func (s *StatsReply) GetHeader() *Header {
	return &s.Header
}

func (s *StatsReply) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *StatsReply) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	n, err = s.Header.Write(buf.Next(8))
	if n == 0 {
		return
	}
	err = binary.Read(buf, binary.BigEndian, &s.Type)
	n += 2
	err = binary.Read(buf, binary.BigEndian, &s.Flags)
	n += 2
	switch s.Type {
	case ST_AGGREGATE:
		a := new(AggregateStatsReply)
		m, aErr := a.Write(buf.Bytes())
		if aErr != nil {
			return
		}
		n += m
	case ST_DESC:
		d := new(DescStats)
		m, dErr := d.Write(buf.Bytes())
		if dErr != nil {
			return
		}
		n += m
	case ST_FLOW:
		for flowCount := buf.Len() / 24; flowCount > 0; flowCount-- {
			f := new(FlowStats)
			m, fErr := f.Write(buf.Next(24))
			if fErr != nil {
				return
			}
			n += m
		}
	case ST_PORT:
		for portCount := buf.Len() / 104; portCount > 0; portCount-- {
			p := new(FlowStats)
			m, pErr := p.Write(buf.Next(104))
			if pErr != nil {
				return
			}
			n += m
		}
	case ST_TABLE:
		for tableCount := buf.Len() / 32; tableCount > 0; tableCount-- {
			t := new(FlowStats)
			m, tErr := t.Write(buf.Next(32))
			if tErr != nil {
				return
			}
			n += m
		}
	case ST_QUEUE:
		for queueCount := buf.Len() / 32; queueCount > 0; queueCount-- {
			q := new(QueueStats)
			m, qErr := q.Write(buf.Next(32))
			if qErr != nil {
				return
			}
			n += m
		}
	case ST_VENDOR:
		break
	}
	return n, nil
}

// _stats_types
const (
	/* Description of this OpenFlow switch.
	* The request body is empty.
	* The reply body is struct ofp_desc_stats. */
	ST_DESC = iota
	/* Individual flow statistics.
	* The request body is struct ofp_flow_stats_request.
	* The reply body is an array of struct ofp_flow_stats. */
	ST_FLOW
	/* Aggregate flow statistics.
	* The request body is struct ofp_aggregate_stats_request.
	* The reply body is struct ofp_aggregate_stats_reply. */
	ST_AGGREGATE
	/* Flow table statistics.
	* The request body is empty.
	* The reply body is an array of struct ofp_table_stats. */
	ST_TABLE
	/* Port statistics.
	* The request body is struct ofp_port_stats_request.
	* The reply body is an array of struct ofp_port_stats. */
	ST_PORT
	/* Queue statistics for a port
	* The request body is struct _queue_stats_request.
	* The reply body is an array of struct ofp_queue_stats */
	ST_QUEUE
	/* Group counter statistics.
	* The request body is struct ofp_group_stats_request.
	* The reply is an array of struct ofp_group_stats. */
	ST_VENDOR = 0xffff
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
			a = NewActionStripVLAN()/*
		case AT_SET_DL_SRC:
			a = NewActionDLSrc()
		case AT_SET_DL_DST:
			a = NewActionDLDst()
		case AT_SET_NW_SRC:
			a = NewActionNWSrc()
		case AT_SET_NW_DST:
			a = NewActionNWDst()
		case AT_SET_NW_TOS:
			a = NewActionNWTOS()
		case AT_SET_TP_SRC:
			a = NewActionTPSrc()
		case AT_SET_TP_DST:
			a = NewActionTPDst()
		case AT_ENQUEUE:
			a = NewActionEnqueue(0, 0)
		case AT_VENDOR:
			a = NewActionVendorPort()*/
		}

		if m, err = a.Write(buf.Next(int(l) - 4)); m == 0 {
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
	Pad         [4]uint8
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
	Pad          [3]uint8
	Name         [MAX_TABLE_NAME_LEN]byte
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
	Pad    [6]uint8
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
	Pad        [6]uint8
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
	Pad     [2]uint8
	QueueID uint32
}

func (s *QueueStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *QueueStatsRequest) Write(b []byte) (n int, err error) {
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
	err = binary.Read(buf, binary.BigEndian, &s.QueueID)
	if err != nil {
		return
	}
	n += 4
	return
}

// ofp_queue_stats 1.0
type QueueStats struct {
	PortNo    uint16
	Pad       [2]uint8
	QueueID   uint32
	TxBytes   uint64
	TxPackets uint64
	TxErrors  uint64
}

func (s *QueueStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *QueueStats) Write(b []byte) (n int, err error) {
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
	err = binary.Read(buf, binary.BigEndian, &s.QueueID)
	if err != nil {
		return
	}
	n += 4
	err = binary.Read(buf, binary.BigEndian, &s.TxBytes)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxPackets)
	if err != nil {
		return
	}
	n += 8
	err = binary.Read(buf, binary.BigEndian, &s.TxErrors)
	if err != nil {
		return
	}
	n += 8
	return
}

// ofp_port_status
type PortStatus struct {
	Header Header
	Reason uint8
	Pad    [7]uint8
	Desc   PhyPort
}

func (p *PortStatus) GetHeader() *Header {
	return &p.Header
}

func (s *PortStatus) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *PortStatus) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.BigEndian, &s.Reason)
	if err != nil {
		return
	}
	n += 1
	err = binary.Read(buf, binary.BigEndian, &s.Pad)
	if err != nil {
		return
	}
	n += 7
	m := 0
	m, err = s.Desc.Write(buf.Bytes())
	if err != nil {
		return
	}
	n += m
	return
}

// ofp_port_reason 1.0
const (
	PR_ADD = iota
	PR_DELETE
	PR_MODIFY
)
