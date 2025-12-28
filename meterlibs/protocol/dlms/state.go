package dlms

import (
	"fmt"
	"reflect"

	"github.com/yimiliya/idis/meterlibs/protocol/dlms/exceptions"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/protocol/acse"
	"github.com/yimiliya/idis/meterlibs/protocol/dlms/protocol/xdlms"
)

// State represents a DLMS connection state
type State struct {
	name string
}

// String returns the string representation of the state
func (s *State) String() string {
	return s.name
}

// Sentinel states
var (
	NoAssociation                    = &State{name: "NO_ASSOCIATION"}
	AwaitingAssociationResponse      = &State{name: "AWAITING_ASSOCIATION_RESPONSE"}
	Ready                            = &State{name: "READY"}
	AwaitingReleaseResponse          = &State{name: "AWAITING_RELEASE_RESPONSE"}
	AwaitingActionResponse           = &State{name: "AWAITING_ACTION_RESPONSE"}
	AwaitingGetResponse              = &State{name: "AWAITING_GET_RESPONSE"}
	AwaitingGetBlockResponse         = &State{name: "AWAITING_GET_BLOCK_RESPONSE"}
	ShouldAckLastGetBlock            = &State{name: "SHOULD_ACK_LAST_GET_BLOCK"}
	AwaitingSetResponse              = &State{name: "AWAITING_SET_RESPONSE"}
	ShouldSendHlsServerChallengeResult = &State{name: "SHOULD_SEND_HLS_SEVER_CHALLENGE_RESULT"}
	AwaitingHlsClientChallengeResult  = &State{name: "AWAITING_HLS_CLIENT_CHALLENGE_RESULT"}
	HlsDone                           = &State{name: "HLS_DONE"}
	NeedData                          = &State{name: "NEED_DATA"}
)

// Flow control events
type HlsStart struct{}

type HlsSuccess struct{}

type HlsFailed struct{}

type RejectAssociation struct{}

type EndAssociation struct{}

// DlmsConnectionState handles state changes in DLMS
type DlmsConnectionState struct {
	currentState *State
}

// NewDlmsConnectionState creates a new DLMS connection state
func NewDlmsConnectionState() *DlmsConnectionState {
	return &DlmsConnectionState{
		currentState: NoAssociation,
	}
}

// NewDlmsConnectionStateWithState creates a new DLMS connection state with a specific state
func NewDlmsConnectionStateWithState(state *State) *DlmsConnectionState {
	return &DlmsConnectionState{
		currentState: state,
	}
}

// CurrentState returns the current state
func (d *DlmsConnectionState) CurrentState() *State {
	return d.currentState
}

// ProcessEvent processes an event and transitions the state machine
func (d *DlmsConnectionState) ProcessEvent(event interface{}) error {
	eventType := reflect.TypeOf(event)
	return d.transitionState(eventType)
}

// transitionState transitions the state based on event type
func (d *DlmsConnectionState) transitionState(eventType reflect.Type) error {
	transitions, ok := dlmsStateTransitions[d.currentState]
	if !ok {
		return fmt.Errorf("no transitions defined for state %s", d.currentState)
	}

	newState, ok := transitions[eventType]
	if !ok {
		return exceptions.NewLocalDlmsProtocolError(
			fmt.Sprintf("can't handle event type %s when state=%s", eventType, d.currentState),
		)
	}

	oldState := d.currentState
	d.currentState = newState
	// TODO: Add logging here if needed
	_ = oldState
	return nil
}

// dlmsStateTransitions defines the state transition table
var dlmsStateTransitions = map[*State]map[reflect.Type]*State{
	NoAssociation: {
		reflect.TypeOf((*acse.ApplicationAssociationRequest)(nil)).Elem(): AwaitingAssociationResponse,
	},
	AwaitingAssociationResponse: {
		reflect.TypeOf((*acse.ApplicationAssociationResponse)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.ExceptionResponse)(nil)).Elem(): NoAssociation,
	},
	Ready: {
		reflect.TypeOf((*acse.ReleaseRequest)(nil)).Elem(): AwaitingReleaseResponse,
		reflect.TypeOf((*xdlms.GetRequestNormal)(nil)).Elem(): AwaitingGetResponse,
		// TODO: GetRequestWithList is not yet implemented
		// reflect.TypeOf((*xdlms.GetRequestWithList)(nil)).Elem(): AwaitingGetResponse,
		reflect.TypeOf((*xdlms.SetRequestNormal)(nil)).Elem(): AwaitingSetResponse,
		reflect.TypeOf((*HlsStart)(nil)).Elem(): ShouldSendHlsServerChallengeResult,
		reflect.TypeOf((*RejectAssociation)(nil)).Elem(): NoAssociation,
		reflect.TypeOf((*xdlms.ActionRequestNormal)(nil)).Elem(): AwaitingActionResponse,
		reflect.TypeOf((*xdlms.DataNotification)(nil)).Elem(): Ready,
		reflect.TypeOf((*EndAssociation)(nil)).Elem(): NoAssociation,
	},
	ShouldSendHlsServerChallengeResult: {
		reflect.TypeOf((*xdlms.ActionRequestNormal)(nil)).Elem(): AwaitingHlsClientChallengeResult,
	},
	AwaitingHlsClientChallengeResult: {
		reflect.TypeOf((*xdlms.ActionResponseNormalWithData)(nil)).Elem(): HlsDone,
		reflect.TypeOf((*xdlms.ActionResponseNormal)(nil)).Elem(): NoAssociation,
		reflect.TypeOf((*xdlms.ActionResponseNormalWithError)(nil)).Elem(): NoAssociation,
	},
	HlsDone: {
		reflect.TypeOf((*HlsSuccess)(nil)).Elem(): Ready,
		reflect.TypeOf((*HlsFailed)(nil)).Elem(): NoAssociation,
	},
	AwaitingGetResponse: {
		reflect.TypeOf((*xdlms.GetResponseNormal)(nil)).Elem(): Ready,
		// TODO: GetResponseWithList is not yet implemented
		// reflect.TypeOf((*xdlms.GetResponseWithList)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.GetResponseWithDataBlock)(nil)).Elem(): ShouldAckLastGetBlock,
		reflect.TypeOf((*xdlms.GetResponseNormalWithError)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.ExceptionResponse)(nil)).Elem(): Ready,
	},
	AwaitingGetBlockResponse: {
		reflect.TypeOf((*xdlms.GetResponseWithDataBlock)(nil)).Elem(): ShouldAckLastGetBlock,
		reflect.TypeOf((*xdlms.GetResponseNormalWithError)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.ExceptionResponse)(nil)).Elem(): Ready,
		// TODO: Add GetResponseLastBlockWithError and GetResponseLastBlock when implemented
	},
	AwaitingSetResponse: {
		reflect.TypeOf((*xdlms.SetResponseNormal)(nil)).Elem(): Ready,
	},
	AwaitingActionResponse: {
		reflect.TypeOf((*xdlms.ActionResponseNormal)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.ActionResponseNormalWithData)(nil)).Elem(): Ready,
		reflect.TypeOf((*xdlms.ActionResponseNormalWithError)(nil)).Elem(): Ready,
	},
	ShouldAckLastGetBlock: {
		reflect.TypeOf((*xdlms.GetRequestNext)(nil)).Elem(): AwaitingGetBlockResponse,
	},
	AwaitingReleaseResponse: {
		reflect.TypeOf((*acse.ReleaseResponse)(nil)).Elem(): NoAssociation,
		reflect.TypeOf((*xdlms.ExceptionResponse)(nil)).Elem(): Ready,
	},
}

