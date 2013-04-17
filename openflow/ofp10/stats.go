package ofp10

import (
	"bytes"
	"encoding/binary"
	"io"
)

// ofp_stats_request 1.0
type OfpStatsRequest struct {
	Header OfpHeader
	Type uint16
	Flags uint16
	Body interface{}
}

func (s *OfpStatsRequest) GetHeader() *OfpHeader {
	return &s.Header
}

func (s *OfpStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if err != nil {
		return
	}
	return n, io.EOF
}

func (s *OfpStatsRequest) Write(b []byte) (n int, err error) {
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
	case OFPST_AGGREGATE:
		// empty
		break
	case OFPST_DESC:
		// empty
		break
	case OFPST_FLOW:
		// ofp_aggregate_stats_request
		a := s.Body.(*OfpAggregateStatsRequest)
		m, err = a.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case OFPST_PORT:
		// ofp_port_stats_request
		p := s.Body.(*OfpPortStatsRequest)
		m, err = p.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case OFPST_TABLE:
		// empty
		break
	case OFPST_QUEUE:
		// ofp_queue_stats_request
		q := s.Body.(*OfpQueueStatsRequest)
		m, err = q.Write(buf.Bytes())
		if m == 0 {
			return
		}
		n += m
	case OFPST_VENDOR:
		break
	}
	return
}

// ofp_stats_reply 1.0
type OfpStatsReply struct {
	Header OfpHeader
	Type uint16
	Flags uint16
	Body []uint8
}

func (s *OfpStatsReply) GetHeader() *OfpHeader {
	return &s.Header
}

func (s *OfpStatsReply) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpStatsReply) Write(b []byte) (n int, err error) {
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
		case OFPST_AGGREGATE:
			a := new(OfpAggregateStatsReply)
			m, aErr := a.Write(buf.Bytes())
			if aErr != nil {
				return
			}
			n += m
		case OFPST_DESC:
			d := new(OfpDescStats)
			m, dErr := d.Write(buf.Bytes())
			if dErr != nil {
				return
			}
			n += m
		case OFPST_FLOW:
			for flowCount := buf.Len() / 24; flowCount > 0; flowCount-- {
				f := new(OfpFlowStats)
				m, fErr := f.Write(buf.Next(24))
				if fErr != nil {
					return
				}
				n += m
			}
		case OFPST_PORT:
			for portCount := buf.Len() / 104; portCount > 0; portCount-- {
				p := new(OfpFlowStats)
				m, pErr := p.Write(buf.Next(104))
				if pErr != nil {
					return
				}
				n += m
			}
		case OFPST_TABLE:
			for tableCount := buf.Len() / 32; tableCount > 0; tableCount-- {
				t := new(OfpFlowStats)
				m, tErr := t.Write(buf.Next(32))
				if tErr != nil {
					return
				}
				n += m
			}
		case OFPST_QUEUE:
			for queueCount := buf.Len() / 32; queueCount > 0; queueCount-- {
				q := new(OfpQueueStats)
				m, qErr := q.Write(buf.Next(32))
				if qErr != nil {
					return
				}
				n += m
			}
		case OFPST_VENDOR:
			break
	}
	return n, nil
}

// ofp_stats_types
const (
	/* Description of this OpenFlow switch.
	* The request body is empty.
	* The reply body is struct ofp_desc_stats. */
	OFPST_DESC = iota
	/* Individual flow statistics.
	* The request body is struct ofp_flow_stats_request.
	* The reply body is an array of struct ofp_flow_stats. */
	OFPST_FLOW
	/* Aggregate flow statistics.
	* The request body is struct ofp_aggregate_stats_request.
	* The reply body is struct ofp_aggregate_stats_reply. */
	OFPST_AGGREGATE
	/* Flow table statistics.
	* The request body is empty.
	* The reply body is an array of struct ofp_table_stats. */
	OFPST_TABLE
	/* Port statistics.
	* The request body is struct ofp_port_stats_request.
	* The reply body is an array of struct ofp_port_stats. */
	OFPST_PORT
	/* Queue statistics for a port
	* The request body is struct ofp_queue_stats_request.
	* The reply body is an array of struct ofp_queue_stats */
	OFPST_QUEUE
	/* Group counter statistics.
	* The request body is struct ofp_group_stats_request.
	* The reply is an array of struct ofp_group_stats. */
	OFPST_VENDOR = 0xffff
)

// ofp_desc_stats 1.0
type OfpDescStats struct {
	MfrDesc [DESC_STR_LEN]byte
	HWDesc [DESC_STR_LEN]byte
	SWDesc [DESC_STR_LEN]byte
	SerialNum [SERIAL_NUM_LEN]byte
	DPDesc [DESC_STR_LEN]byte
}

func (s *OfpDescStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpDescStats) Write(b []byte) (n int, err error) {
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
	DESC_STR_LEN = 256
	SERIAL_NUM_LEN = 32
)

// ofp_flow_stats_request 1.0
type OfpFlowStatsRequest struct {
	Match OfpMatch
	TableID uint8
	Pad uint8
	OutPort uint16
}

func (s *OfpFlowStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpFlowStatsRequest) Write(b []byte) (n int, err error) {
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
type OfpFlowStats struct {
	Length uint16
	TableID uint8
	Pad uint8
	Match OfpMatch
	DurationSec uint32
	DurationNSec uint32
	Priority uint16
	IdleTimeout uint16
	HardTimeout uint16
	Pad2 [6]uint8
	Cookie uint64
	PacketCount uint64
	ByteCount uint64
	Actions []OfpActionHeader
}

func (s *OfpFlowStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpFlowStats) Write(b []byte) (n int, err error) {
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
	for actionCount := buf.Len() / 8; actionCount > 0; actionCount-- {
		a := new(OfpActionHeader)
		a.Write(buf.Next(8))
	}
	return
}

// ofp_aggregate_stats_request 1.0
type OfpAggregateStatsRequest struct {
	Match OfpMatch
	TableID uint8
	Pad uint8
	OutPort uint16
}

func (s *OfpAggregateStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpAggregateStatsRequest) Write(b []byte) (n int, err error) {
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
type OfpAggregateStatsReply struct {
	PacketCount uint64
	ByteCount uint64
	FlowCount uint32
	Pad [4]uint8
}

func (s *OfpAggregateStatsReply) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpAggregateStatsReply) Write(b []byte) (n int, err error) {
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
type OfpTableStats struct {
	TableID uint8
	Pad [3]uint8
	Name [OFP_MAX_TABLE_NAME_LEN]byte
	Wildcards uint32
	MaxEntries uint32
	ActiveCount uint32
	LookupCount uint64
	MatchedCount uint64
}

func (s *OfpTableStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpTableStats) Write(b []byte) (n int, err error) {
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
	OFP_MAX_TABLE_NAME_LEN = 32
)

// ofp_port_stats_request 1.0
type OfpPortStatsRequest struct {
	PortNo uint16
	Pad [6]uint8
}

func (s *OfpPortStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpPortStatsRequest) Write(b []byte) (n int, err error) {
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
type OfpPortStats struct {
	PortNo uint16
	Pad [6]uint8
	RxPackets uint64
	TxPackets uint64
	RxBytes uint64
	TxBytes uint64
	RxDropped uint64
	TxDropped uint64
	RxErrors uint64
	TxErrors uint64
	RxFrameErr uint64
	RxOverErr uint64
	RxCRCErr uint64
	Collisions uint64
}

func (s *OfpPortStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpPortStats) Write(b []byte) (n int, err error) {
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
type OfpQueueStatsRequest struct {
	PortNo uint16
	Pad [2]uint8
	QueueID uint32
}

func (s *OfpQueueStatsRequest) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpQueueStatsRequest) Write(b []byte) (n int, err error) {
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
type OfpQueueStats struct {
	PortNo uint16
	Pad [2]uint8
	QueueID uint32
	TxBytes uint64
	TxPackets uint64
	TxErrors uint64
}

func (s *OfpQueueStats) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpQueueStats) Write(b []byte) (n int, err error) {
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
type OfpPortStatus struct {
	Header OfpHeader
	Reason uint8
	Pad [7]uint8
	Desc OfpPhyPort
}

func (p *OfpPortStatus) GetHeader() *OfpHeader {
	return &p.Header
}

func (s *OfpPortStatus) Read(b []byte) (n int, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, s)
	n, err = buf.Read(b)
	if n == 0 {
		return
	}
	return n, io.EOF
}

func (s *OfpPortStatus) Write(b []byte) (n int, err error) {
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
	OFPPR_ADD = iota
	OFPPR_DELETE
	OFPPR_MODIFY
)
