package xdlms

import "fmt"

// InvokeIdAndPriority represents invoke ID and priority
// invoke_id: It is allowed to send several requests to the server (meter)
//     if the lower layers support it, before listening for the response. To be able to
//     correlate an answer to a request the invoke_id is used. It is copied in the
//     response from the server.
// confirmed: Indicates if the service is confirmed. Mostly it is.
// high_priority: When sending several requests to the server (meter) it is
//     possible to mark some of them as high priority. These response from the requests
//     will be sent back before the ones with normal priority. Handling of priority is
//     a negotiable feature in the Conformance block during Application Association.
//     If the server (meter) does not support priority it will treat all requests with
//     high priority as normal priority.
type InvokeIdAndPriority struct {
	InvokeID     uint8 // 0-15 (4 bits)
	Confirmed    bool
	HighPriority bool
}

const InvokeIdAndPriorityLength = 1

// NewInvokeIdAndPriority creates a new InvokeIdAndPriority
func NewInvokeIdAndPriority(invokeID uint8, confirmed, highPriority bool) (*InvokeIdAndPriority, error) {
	if invokeID > 15 {
		return nil, fmt.Errorf("invoke_id must be between 0-15, got %d", invokeID)
	}
	return &InvokeIdAndPriority{
		InvokeID:     invokeID,
		Confirmed:    confirmed,
		HighPriority: highPriority,
	}, nil
}

// FromBytes creates InvokeIdAndPriority from bytes
func (i *InvokeIdAndPriority) FromBytes(data []byte) (*InvokeIdAndPriority, error) {
	if len(data) != InvokeIdAndPriorityLength {
		return nil, fmt.Errorf("length of data does not correspond with class LENGTH. Should be %d, got %d", InvokeIdAndPriorityLength, len(data))
	}
	
	val := data[0]
	invokeID := val & 0b00001111
	confirmed := (val & 0b01000000) != 0
	highPriority := (val & 0b10000000) != 0
	
	return NewInvokeIdAndPriority(invokeID, confirmed, highPriority)
}

// ToBytes converts InvokeIdAndPriority to bytes
func (i *InvokeIdAndPriority) ToBytes() []byte {
	var out byte
	out = i.InvokeID
	if i.Confirmed {
		out |= 0b01000000
	}
	if i.HighPriority {
		out |= 0b10000000
	}
	return []byte{out}
}

