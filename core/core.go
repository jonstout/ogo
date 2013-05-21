package ogo

import (
	"time"
	"github.com/jonstout/ogo/openflow/ofp10"
)

type Core struct {
	echoRequest chan ofp10.Msg
	portStatus chan ofp10.Msg
}

func (b *Core) InitApplication(args map[string]string) {
	b.echoRequest = SubscribeTo(ofp10.T_ECHO_REQUEST)
	b.portStatus = SubscribeTo(ofp10.T_PORT_STATUS)
}

func (b *Core) Name() string {
	return "Core"
}

func (b *Core) Receive() {
	for {
		select {
		case m := <-b.echoRequest:
			go b.SendEchoReply(m.DPID)
		case m := <-b.portStatus:
			go b.UpdatePortStatus(m)
		}
	}
}

func (b *Core) SendEchoReply(dpid string) {
	if s, ok := GetSwitch(dpid); ok {
		<-time.After(time.Second * 3)
		res := ofp10.NewEchoReply()
		s.Send(res)
	}
}

func (b *Core) UpdatePortStatus(msg ofp10.Msg) {
	if _, ok := GetSwitch(msg.DPID); ok {		
	}
}
