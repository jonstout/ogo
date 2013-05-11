package ogo

import (
	"time"
	"github.com/jonstout/ogo/openflow/ofp10"
)

type BasicApplication struct {
	echoRequest chan ofp10.Msg
}

func (b *BasicApplication) InitApplication(args map[string]string) {
	b.echoRequest = SubscribeTo(ofp10.T_ECHO_REQUEST)
}

func (b *BasicApplication) Name() string {
	return "OgoCore"
}

func (b *BasicApplication) Receive() {
	for {
		select {
		case m := <-b.echoRequest:
			go b.SendEchoReply(m.DPID)
		}
	}
}

func (b *BasicApplication) SendEchoReply(dpid string) {
	req := ofp10.NewEchoReply()
	s, ok := GetSwitch(dpid)
	if ok {
		<-time.After(time.Second * 1)
		s.Send(req)
	}
}
