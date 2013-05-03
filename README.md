# Ogo
An OpenFlow Network Controller in Go

## The Basic Application
All applications must implement the ogo.Application interface.
```
type Application interface {
  InitApplication(args map[string]string)
  GetName() string
  Receive()
}
```
Use the `InitApplication` to recieve command line arguments. The `GetName` function should return a string that will be used to identify your application. Use the `Receive` function to listen on any channels that you have subscribed to.

## Registering your Application
In order for your application to recieve OpenFlow messages from connected switches it must be registered with Ogo.
```
ctrl.RegisterApplication( new(OgoApplication) )
```

## Subscribing to OpenFlow Messages
Use `ogo.SubscribeTo(ofp10.OFPM_*)` to get an ofp10.OfpMsg chan.
```
echoRequestChan := ogo.SubscribeTo(ofp10.OFPT_ECHO_REQUEST)
```

## Acting on Messages
The function `Receive()` is required for all Applications. Use this function to listen for messages on your subscription channels.
```
(app *OgoApplication) Receive() {
for {
    select {
      case msg := <-app.echoRequestChan:
        fmt.Println("Received an EchoRequest message from:", msg.DPID)
      case msg := <-app.anotherOfpChan:
        fmt.Println("Received some other message from:", msg.DPID)
    }
  }
}
```
