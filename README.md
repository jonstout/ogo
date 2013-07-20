# Ogo
An OpenFlow Network Controller in Go

### [Documentation](http://godoc.org/github.com/jonstout/ogo)

## The Basic Application
All applications must implement the ogo.Application interface.
```
type Application interface {
  InitApplication(args map[string]string)
  Name() string
  Receive()
}
```
Use the `InitApplication` to recieve command line arguments. The `Name` function should return a string that will be used to identify your application. Use the `Receive` function to listen on any channels that you have subscribed to.

## Registering your Application
In order for your application to recieve OpenFlow messages from connected switches it must be registered with Ogo.
```
ctrl.RegisterApplication( new(OgoApplication) )
```

## Subscribing to OpenFlow Messages
Use `ogo.SubscribeTo(ofp10.T_*)` to get an ofp10.Msg chan.
```
echoRequestChan := ogo.SubscribeTo(ofp10.T_ECHO_REQUEST)
```

## Acting on Messages
The function `Receive()` is required for all Applications. Use this function to listen for messages on your subscription channels.
```
(app *OgoApplication) Receive() {
for {
    select {
      case msg := <-app.echoRequestChan:
        fmt.Println("Received an EchoRequest message from:", msg.DPID)
      case msg := <-app.anotherChan:
        fmt.Println("Received some other message from:", msg.DPID)
    }
  }
}
```
