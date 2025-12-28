package hdlc

import (
	"fmt"
)

// HdlcState represents the state of HDLC connection
type HdlcState int

const (
	// NotConnected is when we have created a session but not actually set up HDLC
	// connection with the server (meter). We used a SNRM frame to set up the connection
	HdlcStateNotConnected HdlcState = iota
	// Idle State is when we are connected but we have not started a data exchange or we
	// just finished a data exchange
	HdlcStateIdle
	HdlcStateAwaitingResponse
	HdlcStateAwaitingConnection
	HdlcStateAwaitingDisconnect
	HdlcStateClosed
	HdlcStateNeedData
)

// String returns the string representation of the state
func (s HdlcState) String() string {
	switch s {
	case HdlcStateNotConnected:
		return "NOT_CONNECTED"
	case HdlcStateIdle:
		return "IDLE"
	case HdlcStateAwaitingResponse:
		return "AWAITING_RESPONSE"
	case HdlcStateAwaitingConnection:
		return "AWAITING_CONNECTION"
	case HdlcStateAwaitingDisconnect:
		return "AWAITING_DISCONNECT"
	case HdlcStateClosed:
		return "CLOSED"
	case HdlcStateNeedData:
		return "NEED_DATA"
	default:
		return "UNKNOWN"
	}
}

// HdlcConnectionState handles state changes in HDLC
// A HDLC frame is passed to ProcessFrame and it moves the state machine to the
// correct state. If a frame is processed that is not set to be able to transition
// the state in the current state a LocalProtocolError is raised.
type HdlcConnectionState struct {
	CurrentState HdlcState
}

// NewHdlcConnectionState creates a new HdlcConnectionState
func NewHdlcConnectionState() *HdlcConnectionState {
	return &HdlcConnectionState{
		CurrentState: HdlcStateNotConnected,
	}
}

// ProcessFrame processes a frame and transitions the state
func (h *HdlcConnectionState) ProcessFrame(frame interface{}) error {
	frameType := getFrameType(frame)
	newState, ok := hdlcStateTransitions[h.CurrentState][frameType]
	if !ok {
		return NewLocalProtocolError(fmt.Sprintf(
			"can't handle frame type %s when state=%s",
			frameType, h.CurrentState))
	}
	h.CurrentState = newState
	return nil
}

// FrameType represents the type of HDLC frame
type FrameType string

const (
	FrameTypeSetNormalResponseMode FrameType = "SetNormalResponseModeFrame"
	FrameTypeUnNumberedAcknowledgment FrameType = "UnNumberedAcknowledgmentFrame"
	FrameTypeInformation FrameType = "InformationFrame"
	FrameTypeReceiveReady FrameType = "ReceiveReadyFrame"
	FrameTypeDisconnect FrameType = "DisconnectFrame"
)

// getFrameType returns the type of frame
func getFrameType(frame interface{}) FrameType {
	switch frame.(type) {
	case *SetNormalResponseModeFrame:
		return FrameTypeSetNormalResponseMode
	case *UnNumberedAcknowledgmentFrame:
		return FrameTypeUnNumberedAcknowledgment
	case *InformationFrame:
		return FrameTypeInformation
	case *ReceiveReadyFrame:
		return FrameTypeReceiveReady
	case *DisconnectFrame:
		return FrameTypeDisconnect
	default:
		return ""
	}
}

// hdlcStateTransitions defines the state transition table
var hdlcStateTransitions = map[HdlcState]map[FrameType]HdlcState{
	HdlcStateNotConnected: {
		FrameTypeSetNormalResponseMode: HdlcStateAwaitingConnection,
	},
	HdlcStateAwaitingConnection: {
		FrameTypeUnNumberedAcknowledgment: HdlcStateIdle,
	},
	HdlcStateIdle: {
		FrameTypeInformation:        HdlcStateAwaitingResponse,
		FrameTypeDisconnect:         HdlcStateAwaitingDisconnect,
		FrameTypeReceiveReady:       HdlcStateAwaitingResponse,
	},
	HdlcStateAwaitingResponse: {
		FrameTypeInformation:  HdlcStateIdle,
		FrameTypeReceiveReady: HdlcStateIdle,
	},
	HdlcStateAwaitingDisconnect: {
		FrameTypeUnNumberedAcknowledgment: HdlcStateNotConnected,
	},
}

// IsSendState returns true if the current state allows sending
func (h *HdlcConnectionState) IsSendState() bool {
	return h.CurrentState == HdlcStateNotConnected || h.CurrentState == HdlcStateIdle
}

// IsReceiveState returns true if the current state allows receiving
func (h *HdlcConnectionState) IsReceiveState() bool {
	return h.CurrentState == HdlcStateAwaitingConnection ||
		h.CurrentState == HdlcStateAwaitingResponse ||
		h.CurrentState == HdlcStateAwaitingDisconnect
}

