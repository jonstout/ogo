// Messages should be created from this package, unless you desire to
// use a specific version. This package will use the latest supported
// protocol version to generate messages. Hopefully message types don't
// change too much in forthcoming OpenFlow specifications. Unsupported
// message fields will be taken care of by MessageStream.
package ofp

var SupportedVersions []uint16 = []uint16{1}
