package hdlc

import (
	"fmt"
)

const (
	HDLCFlag = 0x7E
	LLCCommandHeader = "\xe6\xe6\x00"
	LLCResponseHeader = "\xe7\xe7\x00"
)

// BaseHdlcFrame is the base class for HDLC frames
// HDLC frames start and end with the HDLC Frame flag 0x7E
// Frame: Flag (1 byte), Format (2 bytes), Destination Address (1-4 bytes),
// Source Address (1-4 bytes), Control (1 byte), Header check sequence (2 bytes),
// Information (n bytes), Frame check sequence (2 bytes), Flag (1 byte)
// The header check sequence field is only present when the frame has an Information field.
type BaseHdlcFrame struct {
	DestinationAddress *HdlcAddress
	SourceAddress    *HdlcAddress
	Payload          []byte
	Segmented        bool
	Final            bool
}

const FixedLengthBytes = 7

// FrameLength returns the total frame length
func (b *BaseHdlcFrame) FrameLength() int {
	return FixedLengthBytes +
		b.DestinationAddress.Length() +
		b.SourceAddress.Length() +
		len(b.Information())
}

// HCS returns the Header Check Sequence
func (b *BaseHdlcFrame) HCS() []byte {
	headerContent := b.HeaderContent()
	if len(headerContent) == 0 {
		return []byte{}
	}
	return HCS.CalculateFor(headerContent, false)
}

// FCS returns the Frame Check Sequence
func (b *BaseHdlcFrame) FCS() []byte {
	frameContent := b.FrameContent()
	return FCS.CalculateFor(frameContent, false)
}

// Information returns the information field
func (b *BaseHdlcFrame) Information() []byte {
	return b.Payload
}

// HeaderContent returns the header content for HCS calculation
func (b *BaseHdlcFrame) HeaderContent() []byte {
	formatField := &DlmsHdlcFrameFormatField{
		Length:    uint16(b.FrameLength()),
		Segmented: b.Segmented,
	}
	formatBytes := formatField.ToBytes()
	
	controlField := b.GetControlField()
	controlBytes := controlField.ToBytes()
	
	result := make([]byte, 0)
	result = append(result, formatBytes...)
	result = append(result, b.DestinationAddress.ToBytes()...)
	result = append(result, b.SourceAddress.ToBytes()...)
	result = append(result, controlBytes...)
	
	return result
}

// FrameContent returns the frame content for FCS calculation
func (b *BaseHdlcFrame) FrameContent() []byte {
	result := make([]byte, 0)
	result = append(result, b.HeaderContent()...)
	hcs := b.HCS()
	if len(hcs) > 0 {
		result = append(result, hcs...)
	}
	result = append(result, b.Information()...)
	return result
}

// ToBytes converts the frame to bytes
func (b *BaseHdlcFrame) ToBytes() []byte {
	result := make([]byte, 0)
	result = append(result, byte(HDLCFlag))
	result = append(result, b.FrameContent()...)
	result = append(result, b.FCS()...)
	result = append(result, byte(HDLCFlag))
	return result
}

// GetControlField returns the control field (to be implemented by specific frame types)
func (b *BaseHdlcFrame) GetControlField() HdlcControlField {
	panic("GetControlField must be implemented by specific frame type")
}

// ExtractFormatFieldFromBytes extracts the format field from frame bytes
func ExtractFormatFieldFromBytes(frameBytes []byte) (*DlmsHdlcFrameFormatField, error) {
	if len(frameBytes) < 3 {
		return nil, fmt.Errorf("frame too short for format field")
	}
	formatField := &DlmsHdlcFrameFormatField{}
	return formatField.FromBytes(frameBytes[1:3])
}

// FrameIsEnclosedByHdlcFlags checks if frame is enclosed by HDLC flags
func FrameIsEnclosedByHdlcFlags(frameBytes []byte) bool {
	if len(frameBytes) < 2 {
		return false
	}
	first := frameBytes[0]
	last := frameBytes[len(frameBytes)-1]
	return first == HDLCFlag && last == HDLCFlag && first == last
}

// FrameHasCorrectLength checks if frame has correct length
func FrameHasCorrectLength(controlFieldLength int, frameBytes []byte) bool {
	return (controlFieldLength + 2) == len(frameBytes)
}

// SetNormalResponseModeFrame (SNRM-frame) is used to start a new HDLC connection
type SetNormalResponseModeFrame struct {
	*BaseHdlcFrame
}

// NewSetNormalResponseModeFrame creates a new SNRM frame
func NewSetNormalResponseModeFrame(destinationAddress, sourceAddress *HdlcAddress) *SetNormalResponseModeFrame {
	return &SetNormalResponseModeFrame{
		BaseHdlcFrame: &BaseHdlcFrame{
			DestinationAddress: destinationAddress,
			SourceAddress:      sourceAddress,
			Final:              true,
		},
	}
}

// HCS returns empty bytes (SNRM is an S-frame without information field)
func (s *SetNormalResponseModeFrame) HCS() []byte {
	return []byte{}
}

// Information returns empty bytes (no information field for SNRM)
func (s *SetNormalResponseModeFrame) Information() []byte {
	return []byte{}
}

// GetControlField returns the SNRM control field
func (s *SetNormalResponseModeFrame) GetControlField() HdlcControlField {
	return NewSnrmControlField()
}

// FrameLength returns the frame length for SNRM
func (s *SetNormalResponseModeFrame) FrameLength() int {
	return 5 + // fixed length without HCS
		s.DestinationAddress.Length() +
		s.SourceAddress.Length()
}

// UnNumberedAcknowledgmentFrame (UA-frame) is used to acknowledge SNRM
type UnNumberedAcknowledgmentFrame struct {
	*BaseHdlcFrame
}

// NewUnNumberedAcknowledgmentFrame creates a new UA frame
func NewUnNumberedAcknowledgmentFrame(destinationAddress, sourceAddress *HdlcAddress, payload []byte) *UnNumberedAcknowledgmentFrame {
	return &UnNumberedAcknowledgmentFrame{
		BaseHdlcFrame: &BaseHdlcFrame{
			DestinationAddress: destinationAddress,
			SourceAddress:      sourceAddress,
			Payload:            payload,
			Final:              true,
		},
	}
}

// FrameLength returns the frame length for UA
func (u *UnNumberedAcknowledgmentFrame) FrameLength() int {
	fixed := 7
	if len(u.Information()) == 0 {
		fixed = 5 // without HCS
	}
	return fixed +
		u.DestinationAddress.Length() +
		u.SourceAddress.Length() +
		len(u.Information())
}

// HCS returns HCS if information field is present
func (u *UnNumberedAcknowledgmentFrame) HCS() []byte {
	if len(u.Payload) > 0 {
		return u.BaseHdlcFrame.HCS()
	}
	return []byte{}
}

// GetControlField returns the UA control field
func (u *UnNumberedAcknowledgmentFrame) GetControlField() HdlcControlField {
	return NewUaControlField()
}

// FromBytes creates a UA frame from bytes
func (u *UnNumberedAcknowledgmentFrame) FromBytes(frameBytes []byte) (*UnNumberedAcknowledgmentFrame, error) {
	if !FrameIsEnclosedByHdlcFlags(frameBytes) {
		return nil, NewMissingHdlcFlags()
	}

	formatField, err := ExtractFormatFieldFromBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	if !FrameHasCorrectLength(int(formatField.Length), frameBytes) {
		return nil, NewHdlcParsingError(fmt.Sprintf(
			"frame data is not of length specified in frame format field. Should be %d but is %d",
			formatField.Length, len(frameBytes)))
	}

	destinationAddress, err := DestinationFromBytes(frameBytes, AddressTypeClient)
	if err != nil {
		return nil, err
	}
	sourceAddress, err := SourceFromBytes(frameBytes, AddressTypeServer)
	if err != nil {
		return nil, err
	}

	hcsPosition := 1 + 2 + destinationAddress.Length() + sourceAddress.Length() + 1
	hcs := frameBytes[hcsPosition : hcsPosition+2]
	fcs := frameBytes[len(frameBytes)-3 : len(frameBytes)-1]
	information := frameBytes[hcsPosition+2 : len(frameBytes)-3]

	frame := NewUnNumberedAcknowledgmentFrame(destinationAddress, sourceAddress, information)

	if len(frame.HCS()) > 0 {
		calculatedHCS := frame.HCS()
		if len(hcs) != len(calculatedHCS) {
			return nil, NewHdlcParsingError("HCS length mismatch")
		}
		for i := range hcs {
			if hcs[i] != calculatedHCS[i] {
				return nil, NewHdlcParsingError(fmt.Sprintf("HCS is not correct. Calculated: %v, in data: %v", calculatedHCS, hcs))
			}
		}
	}

	calculatedFCS := frame.FCS()
	if len(fcs) != len(calculatedFCS) {
		return nil, NewHdlcParsingError("FCS length mismatch")
	}
	for i := range fcs {
		if fcs[i] != calculatedFCS[i] {
			return nil, NewHdlcParsingError(fmt.Sprintf("FCS is not correct. Calculated: %v, in data: %v", calculatedFCS, fcs))
		}
	}

	return frame, nil
}

// ReceiveReadyFrame (RR-frame) is used for acknowledgment
type ReceiveReadyFrame struct {
	*BaseHdlcFrame
	ReceiveSequenceNumber uint8
}

// NewReceiveReadyFrame creates a new RR frame
func NewReceiveReadyFrame(destinationAddress, sourceAddress *HdlcAddress, receiveSequenceNumber uint8) (*ReceiveReadyFrame, error) {
	rr := &ReceiveReadyFrame{
		BaseHdlcFrame: &BaseHdlcFrame{
			DestinationAddress: destinationAddress,
			SourceAddress:      sourceAddress,
			Final:              true,
		},
		ReceiveSequenceNumber: receiveSequenceNumber,
	}
	return rr, nil
}

// HCS returns empty bytes (no information field)
func (r *ReceiveReadyFrame) HCS() []byte {
	return []byte{}
}

// Information returns empty bytes
func (r *ReceiveReadyFrame) Information() []byte {
	return []byte{}
}

// GetControlField returns the RR control field
func (r *ReceiveReadyFrame) GetControlField() HdlcControlField {
	control, _ := NewReceiveReadyControlField(r.ReceiveSequenceNumber)
	return control
}

// FromBytes creates a RR frame from bytes
func (r *ReceiveReadyFrame) FromBytes(frameBytes []byte) (*ReceiveReadyFrame, error) {
	if !FrameIsEnclosedByHdlcFlags(frameBytes) {
		return nil, NewMissingHdlcFlags()
	}

	formatField, err := ExtractFormatFieldFromBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	if !FrameHasCorrectLength(int(formatField.Length), frameBytes) {
		return nil, NewHdlcParsingError(fmt.Sprintf(
			"frame data is not of length specified in frame format field. Should be %d but is %d",
			formatField.Length, len(frameBytes)))
	}

	destinationAddress, err := DestinationFromBytes(frameBytes, AddressTypeClient)
	if err != nil {
		return nil, err
	}
	sourceAddress, err := SourceFromBytes(frameBytes, AddressTypeServer)
	if err != nil {
		return nil, err
	}

	controlBytePosition := 1 + 2 + destinationAddress.Length() + sourceAddress.Length()
	controlByte := frameBytes[controlBytePosition : controlBytePosition+1]
	controlField := &ReceiveReadyControlField{}
	control, err := controlField.FromBytes(controlByte)
	if err != nil {
		return nil, err
	}

	fcs := frameBytes[len(frameBytes)-3 : len(frameBytes)-1]

	frame, err := NewReceiveReadyFrame(destinationAddress, sourceAddress, control.ReceiveSequenceNumber)
	if err != nil {
		return nil, err
	}

	calculatedFCS := frame.FCS()
	if len(fcs) != len(calculatedFCS) {
		return nil, NewHdlcParsingError("FCS length mismatch")
	}
	for i := range fcs {
		if fcs[i] != calculatedFCS[i] {
			return nil, NewHdlcParsingError("FCS is not correct")
		}
	}

	return frame, nil
}

// InformationFrame is used for data transmission
type InformationFrame struct {
	*BaseHdlcFrame
	SendSequenceNumber    uint8
	ReceiveSequenceNumber uint8
}

// NewInformationFrame creates a new Information frame
func NewInformationFrame(
	destinationAddress, sourceAddress *HdlcAddress,
	payload []byte,
	sendSequenceNumber, receiveSequenceNumber uint8,
	segmented, final bool,
) (*InformationFrame, error) {
	return &InformationFrame{
		BaseHdlcFrame: &BaseHdlcFrame{
			DestinationAddress: destinationAddress,
			SourceAddress:      sourceAddress,
			Payload:            payload,
			Segmented:          segmented,
			Final:              final,
		},
		SendSequenceNumber:    sendSequenceNumber,
		ReceiveSequenceNumber: receiveSequenceNumber,
	}, nil
}

// Information returns the information field with LLC header
func (i *InformationFrame) Information() []byte {
	if len(i.Payload) == 0 {
		return []byte{}
	}
	result := make([]byte, 0)
	result = append(result, []byte(LLCCommandHeader)...)
	result = append(result, i.Payload...)
	return result
}

// GetControlField returns the Information control field
func (i *InformationFrame) GetControlField() HdlcControlField {
	control, _ := NewInformationControlField(i.SendSequenceNumber, i.ReceiveSequenceNumber, i.Final)
	return control
}

// FromBytes creates an Information frame from bytes
func (i *InformationFrame) FromBytes(frameBytes []byte) (*InformationFrame, error) {
	if !FrameIsEnclosedByHdlcFlags(frameBytes) {
		return nil, NewMissingHdlcFlags()
	}

	formatField, err := ExtractFormatFieldFromBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	if !FrameHasCorrectLength(int(formatField.Length), frameBytes) {
		return nil, NewHdlcParsingError(fmt.Sprintf(
			"frame data is not of length specified in frame format field. Should be %d but is %d",
			formatField.Length, len(frameBytes)))
	}

	destinationAddress, err := DestinationFromBytes(frameBytes, AddressTypeClient)
	if err != nil {
		return nil, err
	}
	sourceAddress, err := SourceFromBytes(frameBytes, AddressTypeServer)
	if err != nil {
		return nil, err
	}

	informationControlBytePosition := 1 + 2 + destinationAddress.Length() + sourceAddress.Length()
	informationControlByte := frameBytes[informationControlBytePosition : informationControlBytePosition+1]
	controlField := &InformationControlField{}
	informationControl, err := controlField.FromBytes(informationControlByte)
	if err != nil {
		return nil, err
	}

	hcsPosition := 1 + 2 + destinationAddress.Length() + sourceAddress.Length() + 1
	hcs := frameBytes[hcsPosition : hcsPosition+2]
	fcs := frameBytes[len(frameBytes)-3 : len(frameBytes)-1]
	information := frameBytes[hcsPosition+2 : len(frameBytes)-3]

	// Remove LLC header if present
	payload := information
	if len(information) >= 3 && string(information[:3]) == LLCCommandHeader {
		payload = information[3:]
	} else if len(information) >= 3 && string(information[:3]) == LLCResponseHeader {
		payload = information[3:]
	}

	frame, err := NewInformationFrame(
		destinationAddress,
		sourceAddress,
		payload,
		informationControl.SendSequenceNumber,
		informationControl.ReceiveSequenceNumber,
		formatField.Segmented,
		informationControl.Final,
	)
	if err != nil {
		return nil, err
	}

	calculatedHCS := frame.HCS()
	if len(hcs) != len(calculatedHCS) {
		return nil, NewHdlcParsingError("HCS length mismatch")
	}
	for i := range hcs {
		if hcs[i] != calculatedHCS[i] {
			return nil, NewHdlcParsingError(fmt.Sprintf("HCS is not correct. Calculated: %v, in data: %v", calculatedHCS, hcs))
		}
	}

	calculatedFCS := frame.FCS()
	if len(fcs) != len(calculatedFCS) {
		return nil, NewHdlcParsingError("FCS length mismatch")
	}
	for i := range fcs {
		if fcs[i] != calculatedFCS[i] {
			return nil, NewHdlcParsingError(fmt.Sprintf("FCS is not correct. Calculated: %v, in data: %v", calculatedFCS, fcs))
		}
	}

	return frame, nil
}

// DisconnectFrame is used to disconnect HDLC connection
type DisconnectFrame struct {
	*BaseHdlcFrame
}

// NewDisconnectFrame creates a new Disconnect frame
func NewDisconnectFrame(destinationAddress, sourceAddress *HdlcAddress) *DisconnectFrame {
	return &DisconnectFrame{
		BaseHdlcFrame: &BaseHdlcFrame{
			DestinationAddress: destinationAddress,
			SourceAddress:      sourceAddress,
			Final:              true,
		},
	}
}

// HCS returns empty bytes (no information field)
func (d *DisconnectFrame) HCS() []byte {
	return []byte{}
}

// Information returns empty bytes
func (d *DisconnectFrame) Information() []byte {
	return []byte{}
}

// GetControlField returns the Disconnect control field
func (d *DisconnectFrame) GetControlField() HdlcControlField {
	return NewDisconnectControlField()
}

// FromBytes creates a Disconnect frame from bytes
func (d *DisconnectFrame) FromBytes(frameBytes []byte) (*DisconnectFrame, error) {
	if !FrameIsEnclosedByHdlcFlags(frameBytes) {
		return nil, NewMissingHdlcFlags()
	}

	formatField, err := ExtractFormatFieldFromBytes(frameBytes)
	if err != nil {
		return nil, err
	}

	if !FrameHasCorrectLength(int(formatField.Length), frameBytes) {
		return nil, NewHdlcParsingError(fmt.Sprintf(
			"frame data is not of length specified in frame format field. Should be %d but is %d",
			formatField.Length, len(frameBytes)))
	}

	destinationAddress, err := DestinationFromBytes(frameBytes, AddressTypeServer)
	if err != nil {
		return nil, err
	}
	sourceAddress, err := SourceFromBytes(frameBytes, AddressTypeClient)
	if err != nil {
		return nil, err
	}

	fcs := frameBytes[len(frameBytes)-3 : len(frameBytes)-1]

	frame := NewDisconnectFrame(destinationAddress, sourceAddress)

	calculatedFCS := frame.FCS()
	if len(fcs) != len(calculatedFCS) {
		return nil, NewHdlcParsingError("FCS length mismatch")
	}
	for i := range fcs {
		if fcs[i] != calculatedFCS[i] {
			return nil, NewHdlcParsingError("FCS is not correct")
		}
	}

	return frame, nil
}

