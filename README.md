[![Build Status](https://drone.io/github.com/jonstout/ogo/status.png)](https://drone.io/github.com/jonstout/ogo/latest) [![GoDoc](https://godoc.org/github.com/jonstout/ogo?status.png)](https://godoc.org/github.com/jonstout/ogo)
# Ogo
An OpenFlow Network Controller written in Go.

## A Basic Application
### Register
To process OpenFlow messages register a function that returns a pointer to an
existing or new Application struct.
```
func NewDemoInstance() interface{} {
  return &DemoInstance{}
}
controller.RegisterApplication(NewDemoInstance)
```

### Receive
To receive OpenFlow messages, applications should implement the interfaces
found in `protocol/ofp10/interface.go` or `protocol/ofp13/interface.go`.
```
func (b *DemoInstance) ConnectionUp(dpid net.HardwareAddr) {
  log.Println("Switch connected:", dpid)
}

func (b *DemoInstance) ConnectionDown(dpid net.HardwareAddr) {
  log.Println("Switch disconnected:", dpid)
}

func (b *DemoInstance) PacketIn(dpid net.HardwareAddr, pkt *ofp10.PacketIn) {
  log.Println("PacketIn message received from:", dpid)
}
```

### Send
Any struct that implements `util.Message` can be sent to the switch. Only
OpenFlow messages should be sent using `OFSwitch.Send(m util.Message)`.
```
req := ofp10.NewEchoRequest()

// If switch dpid is known, returns its OFPSwitch struct. The
// switch is not guaranteed to have an active connection.
if sw, ok := ogo.Switch(dpid string); ok {
  sw.Send(req)
}
```
