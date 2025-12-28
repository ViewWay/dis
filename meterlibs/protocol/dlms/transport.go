package dlms

import "log"

// DataChannel is a channel for receiving data from the transport layer
type DataChannel chan []byte

// Transport specifies the transport layer interface
type Transport interface {
	Close()
	Connect() (err error)
	Disconnect() (err error)
	IsConnected() bool
	SetAddress(client int, server int)
	SetReception(dc DataChannel)
	Send(src []byte) error
	SetLogger(logger *log.Logger)
}

// TransportWithBroadcast are optional methods for transport layers with broadcast capabilities
type TransportWithBroadcast interface {
	SendBroadcast(src []byte) error
}
