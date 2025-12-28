package hdlc

import (
	"fmt"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/encoding"
)

// HdlcControlField is the interface for HDLC control fields
// Control field is represented by 1 byte of data.
// Indicates the type of commands or responses, and contains HDLC
// sequence numbers, where appropriate. The last bits of the control field (CTRL)
// identify the type of the HDLC frame
// The is_final represents the poll/final bit. Indicates if the frame is the last one
// of a sequence. Setting the bit releases the control to the other party.
type HdlcControlField interface {
	IsFinal() bool
	ToBytes() []byte
}

// SnrmControlField is an S-frame for SNRM request
type SnrmControlField struct{}

// NewSnrmControlField creates a new SnrmControlField
func NewSnrmControlField() *SnrmControlField {
	return &SnrmControlField{}
}

// IsFinal returns true (almost all the time a SNRM frame is contained in single frame)
func (s *SnrmControlField) IsFinal() bool {
	return true
}

// ToBytes converts SnrmControlField to bytes
func (s *SnrmControlField) ToBytes() []byte {
	out := byte(0b10000011)
	if s.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// UaControlField is an S-frame for Unacknowledge Answer
type UaControlField struct{}

// NewUaControlField creates a new UaControlField
func NewUaControlField() *UaControlField {
	return &UaControlField{}
}

// IsFinal returns true (most UA is only one frame)
func (u *UaControlField) IsFinal() bool {
	return true
}

// ToBytes converts UaControlField to bytes
func (u *UaControlField) ToBytes() []byte {
	out := byte(0b01100011)
	if u.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// DisconnectControlField is an S-frame for disconnect
type DisconnectControlField struct{}

// NewDisconnectControlField creates a new DisconnectControlField
func NewDisconnectControlField() *DisconnectControlField {
	return &DisconnectControlField{}
}

// IsFinal returns true (always final)
func (d *DisconnectControlField) IsFinal() bool {
	return true
}

// ToBytes converts DisconnectControlField to bytes
func (d *DisconnectControlField) ToBytes() []byte {
	out := byte(0b01000011)
	if d.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// ReceiveReadyControlField is an RR-frame for ack
type ReceiveReadyControlField struct {
	ReceiveSequenceNumber uint8 // 0-7
}

// NewReceiveReadyControlField creates a new ReceiveReadyControlField
func NewReceiveReadyControlField(receiveSequenceNumber uint8) (*ReceiveReadyControlField, error) {
	if receiveSequenceNumber > 7 {
		return nil, fmt.Errorf("sequence number can only be between 0-7, got %d", receiveSequenceNumber)
	}
	return &ReceiveReadyControlField{
		ReceiveSequenceNumber: receiveSequenceNumber,
	}, nil
}

// IsFinal returns true (always final)
func (r *ReceiveReadyControlField) IsFinal() bool {
	return true
}

// ToBytes converts ReceiveReadyControlField to bytes
func (r *ReceiveReadyControlField) ToBytes() []byte {
	out := byte(0b00000001)
	out += r.ReceiveSequenceNumber << 5
	if r.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// FromBytes creates a ReceiveReadyControlField from bytes
func (r *ReceiveReadyControlField) FromBytes(inByte []byte) (*ReceiveReadyControlField, error) {
	if len(inByte) != 1 {
		return nil, fmt.Errorf("ReceiveReadyControlField can only be 1 byte, got %d", len(inByte))
	}
	value := inByte[0]
	controlFrame := value&0b00000001 != 0
	if !controlFrame {
		return nil, fmt.Errorf("frame is an information frame not a ReceiveReadyFrame")
	}
	rsn := (value & 0b11100000) >> 5
	return NewReceiveReadyControlField(rsn)
}

// InformationControlField contains information about the acknowledge frames
// sent between the client and server.
// The send_sequence_number holds information about the enumeration of the current
// frame in transit.
// The receive_sequence_number holds information about the enumeration of the next
// frame the sender is expecting to be delivered.
// send_sequence_number and receive_sequence_number are in DLMS limited to 3 bits
// and can take the value of 0-7.
type InformationControlField struct {
	SendSequenceNumber    uint8 // 0-7
	ReceiveSequenceNumber uint8 // 0-7
	Final                 bool
}

// NewInformationControlField creates a new InformationControlField
func NewInformationControlField(sendSequenceNumber, receiveSequenceNumber uint8, final bool) (*InformationControlField, error) {
	if sendSequenceNumber > 7 {
		return nil, fmt.Errorf("send sequence number can only be between 0-7, got %d", sendSequenceNumber)
	}
	if receiveSequenceNumber > 7 {
		return nil, fmt.Errorf("receive sequence number can only be between 0-7, got %d", receiveSequenceNumber)
	}
	return &InformationControlField{
		SendSequenceNumber:    sendSequenceNumber,
		ReceiveSequenceNumber: receiveSequenceNumber,
		Final:                 final,
	}, nil
}

// IsFinal returns the final flag
func (i *InformationControlField) IsFinal() bool {
	return i.Final
}

// ToBytes converts InformationControlField to bytes
func (i *InformationControlField) ToBytes() []byte {
	out := byte(0b00000000)
	out += i.SendSequenceNumber << 1
	out += i.ReceiveSequenceNumber << 5
	if i.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// FromBytes creates an InformationControlField from bytes
func (i *InformationControlField) FromBytes(inByte []byte) (*InformationControlField, error) {
	if len(inByte) != 1 {
		return nil, fmt.Errorf("InformationControlField can only be 1 byte, got %d", len(inByte))
	}
	value := inByte[0]
	notInfoFrame := value&0b00000001 != 0
	if notInfoFrame {
		return nil, fmt.Errorf("byte is not representing an InformationControlField. LSB is 1, should be 0")
	}
	ssn := (value & 0b00001110) >> 1
	rsn := (value & 0b11100000) >> 5
	final := value&0b00010000 != 0
	return NewInformationControlField(ssn, rsn, final)
}

// UnnumberedInformationControlField is used for UnnumberedInformationFrames
type UnnumberedInformationControlField struct {
	Final bool
}

// NewUnnumberedInformationControlField creates a new UnnumberedInformationControlField
func NewUnnumberedInformationControlField(final bool) *UnnumberedInformationControlField {
	return &UnnumberedInformationControlField{Final: final}
}

// IsFinal returns the final flag
func (u *UnnumberedInformationControlField) IsFinal() bool {
	return u.Final
}

// ToBytes converts UnnumberedInformationControlField to bytes
func (u *UnnumberedInformationControlField) ToBytes() []byte {
	out := byte(0b00000011)
	if u.IsFinal() {
		out |= 0b00010000
	}
	return []byte{out}
}

// FromBytes creates an UnnumberedInformationControlField from bytes
func (u *UnnumberedInformationControlField) FromBytes(inByte []byte) (*UnnumberedInformationControlField, error) {
	if len(inByte) != 1 {
		return nil, fmt.Errorf("UnnumberedInformationControlField can only be 1 byte, got %d", len(inByte))
	}
	value := inByte[0]
	isUnnumberedInfoFrame := value&0b00000011 == 0b00000011
	if !isUnnumberedInfoFrame {
		return nil, fmt.Errorf("byte is not representing an UnnumberedInformationControlField")
	}
	final := value&0b00010000 != 0
	return NewUnnumberedInformationControlField(final), nil
}

// DlmsHdlcFrameFormatField represents the HDLC frame format field (2 bytes)
// The 4 leftmost bits represents the HDLC frame format.
// DLMS used HDLC frame format 3 -> 0b1010 -> 0xA. This is always the same in all frames
// The bit 11 is the segmentation bit. If set it indicates that the
// data the frame consists of is not complete and has been segmented into several frames.
// The bit 0-10 rightmost bits represents the frame length.
// Length of a frame is calculated excluding the frame tags (0x7e).
type DlmsHdlcFrameFormatField struct {
	Length    uint16 // 0-2047
	Segmented bool
}

// NewDlmsHdlcFrameFormatField creates a new DlmsHdlcFrameFormatField
func NewDlmsHdlcFrameFormatField(length uint16, segmented bool) (*DlmsHdlcFrameFormatField, error) {
	if length > 0b11111111111 {
		return nil, fmt.Errorf("frame length is too long: %d", length)
	}
	return &DlmsHdlcFrameFormatField{
		Length:    length,
		Segmented: segmented,
	}, nil
}

// FromBytes creates a DlmsHdlcFrameFormatField from bytes
func (d *DlmsHdlcFrameFormatField) FromBytes(inBytes []byte) (*DlmsHdlcFrameFormatField, error) {
	if len(inBytes) != 2 {
		return nil, NewHdlcParsingError(fmt.Sprintf("HDLC frame format length is %d, should be 2", len(inBytes)))
	}
	if !d.correctFrameFormat(inBytes) {
		return nil, NewHdlcParsingError(fmt.Sprintf("received a HDLC frame of the incorrect format: %v", inBytes))
	}
	segmented := inBytes[0]&0b00001000 != 0
	length := d.getLengthFromBytes(inBytes)
	return NewDlmsHdlcFrameFormatField(length, segmented)
}

// correctFrameFormat checks if the frame format is correct
func (d *DlmsHdlcFrameFormatField) correctFrameFormat(bytes []byte) bool {
	leftmost := bytes[0]
	maskedLeftmost := leftmost & 0b11110000
	return maskedLeftmost == 0b10100000
}

// getLengthFromBytes extracts length from bytes (rightmost 11 bits)
func (d *DlmsHdlcFrameFormatField) getLengthFromBytes(bytes []byte) uint16 {
	total := uint16(bytes[0])<<8 | uint16(bytes[1])
	return total & 0b0000011111111111
}

// ToBytes converts DlmsHdlcFrameFormatField to bytes
func (d *DlmsHdlcFrameFormatField) ToBytes() []byte {
	total := uint16(0b1010000000000000) | d.Length
	if d.Segmented {
		total |= 0b0000100000000000
	}
	return []byte{byte(total >> 8), byte(total & 0xFF)}
}

// HCS is the Header Check Sequence calculator
var HCS = encoding.NewCRCCCITT()

// FCS is the Frame Check Sequence calculator
var FCS = encoding.NewCRCCCITT()

