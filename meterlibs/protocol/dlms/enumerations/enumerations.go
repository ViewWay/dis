package enumerations

// DataAccessResult represents the result of data access operations
type DataAccessResult uint8

const (
	DataAccessSuccess                DataAccessResult = 0
	DataAccessHardwareFault          DataAccessResult = 1
	DataAccessTemporaryFailure      DataAccessResult = 2
	DataAccessReadWriteDenied        DataAccessResult = 3
	DataAccessObjectUndefined        DataAccessResult = 4
	DataAccessObjectClassInconsistent DataAccessResult = 9
	DataAccessObjectUnavailable      DataAccessResult = 11
	DataAccessTypeUnmatched          DataAccessResult = 12
	DataAccessScopeOfAccessViolated  DataAccessResult = 13
	DataAccessDataBlockUnavailable   DataAccessResult = 14
	DataAccessLongGetAborted         DataAccessResult = 15
	DataAccessNoLongGetInProgress    DataAccessResult = 16
	DataAccessLongSetAborted         DataAccessResult = 17
	DataAccessNoLongSetInProgress    DataAccessResult = 18
	DataAccessDataBlockNumberInvalid DataAccessResult = 19
	DataAccessOtherReason            DataAccessResult = 250
)

// GetRequestType represents the type of GET request
type GetRequestType uint8

const (
	GetRequestNormal   GetRequestType = 1
	GetRequestNext     GetRequestType = 2
	GetRequestWithList GetRequestType = 3
)

// GetResponseType represents the type of GET response
type GetResponseType uint8

const (
	GetResponseNormal            GetResponseType = 1
	GetResponseWithBlock         GetResponseType = 2
	GetResponseWithList          GetResponseType = 3
	GetResponseLastBlock         GetResponseType = 4
	GetResponseLastBlockWithError GetResponseType = 5
)

// SetRequestType represents the type of SET request
type SetRequestType uint8

const (
	SetRequestNormal          SetRequestType = 1
	SetRequestWithFirstBlock   SetRequestType = 2
	SetRequestWithBlock        SetRequestType = 3
	SetRequestWithList         SetRequestType = 4
	SetRequestFirstBlockWithList SetRequestType = 5
)

// SetResponseType represents the type of SET response
type SetResponseType uint8

const (
	SetResponseNormal         SetResponseType = 1
	SetResponseWithBlock       SetResponseType = 2
	SetResponseWithLastBlock   SetResponseType = 3
	SetResponseLastBlockWithList SetResponseType = 4
	SetResponseWithList        SetResponseType = 5
)

// ActionType represents the type of ACTION request
type ActionType uint8

const (
	ActionNormal              ActionType = 1
	ActionNextPBlock          ActionType = 2
	ActionWithList            ActionType = 3
	ActionWithFirstPBlock     ActionType = 4
	ActionWithListAndFirstPBlock ActionType = 5
	ActionWithPBlock          ActionType = 6
)

// StateException represents state exception types
type StateException uint8

const (
	StateExceptionServiceNotAllowed StateException = 1
	StateExceptionServiceUnknown    StateException = 2
)

// ServiceException represents service exception types
type ServiceException uint8

const (
	ServiceExceptionOperationNotPossible ServiceException = 1
	ServiceExceptionServiceNotSupported  ServiceException = 2
	ServiceExceptionOtherReason          ServiceException = 3
	ServiceExceptionPDUTooLong           ServiceException = 4
	ServiceExceptionDecipheringError     ServiceException = 5
	ServiceExceptionInvocationCounterError ServiceException = 6
)

// ApplicationReferenceError represents application reference error types
type ApplicationReferenceError uint8

const (
	ApplicationReferenceErrorOther                    ApplicationReferenceError = 0
	ApplicationReferenceErrorTimeElapsed            ApplicationReferenceError = 1
	ApplicationReferenceErrorApplicationUnreachable ApplicationReferenceError = 2
	ApplicationReferenceErrorApplicationReferenceInvalid ApplicationReferenceError = 3
	ApplicationReferenceErrorApplicationContextUnsupported ApplicationReferenceError = 4
	ApplicationReferenceErrorProviderCommunicationError ApplicationReferenceError = 5
	ApplicationReferenceErrorDecipheringError        ApplicationReferenceError = 6
)

// HardwareResourceError represents hardware resource error types
type HardwareResourceError uint8

const (
	HardwareResourceErrorOther                  HardwareResourceError = 0
	HardwareResourceErrorMemoryUnavailable      HardwareResourceError = 1
	HardwareResourceErrorProcessorResourceUnavailable HardwareResourceError = 2
	HardwareResourceErrorMassStorageUnavailable HardwareResourceError = 3
	HardwareResourceErrorOtherResourceUnavailable HardwareResourceError = 4
)

// VdeStateError represents VDE state error types
type VdeStateError uint8

const (
	VdeStateErrorOther        VdeStateError = 0
	VdeStateErrorNoDlmsContext VdeStateError = 1
	VdeStateErrorLoadingDataset VdeStateError = 2
	VdeStateErrorStatusNoChange VdeStateError = 3
	VdeStateErrorStatusInoperable VdeStateError = 4
)

// ServiceError represents service error types
type ServiceError uint8

const (
	ServiceErrorOther            ServiceError = 0
	ServiceErrorPDUSize          ServiceError = 1
	ServiceErrorServiceUnsupported ServiceError = 2
)

// DefinitionError represents definition error types
type DefinitionError uint8

const (
	DefinitionErrorOther                  DefinitionError = 0
	DefinitionErrorObjectUndefined        DefinitionError = 1
	DefinitionErrorObjectClassInconsistent DefinitionError = 2
	DefinitionErrorObjectAttributeInconsistent DefinitionError = 3
)

// AccessError represents access error types
type AccessError uint8

const (
	AccessErrorOther                AccessError = 0
	AccessErrorScopeOfAccessViolated AccessError = 1
	AccessErrorObjectAccessViolated  AccessError = 2
	AccessErrorHardwareFault         AccessError = 3
	AccessErrorObjectUnavailable     AccessError = 4
)

// InitiateError represents initiate error types
type InitiateError uint8

const (
	InitiateErrorOther              InitiateError = 0
	InitiateErrorDlmsVersionTooLow  InitiateError = 1
	InitiateErrorIncompatibleConformance InitiateError = 2
	InitiateErrorPDUSizeTooShort    InitiateError = 3
	InitiateErrorRefusedByVdeHandler InitiateError = 4
)

// LoadDataError represents load data error types
type LoadDataError uint8

const (
	LoadDataErrorOther              LoadDataError = 0
	LoadDataErrorPrimitiveOutOfSequence LoadDataError = 1
	LoadDataErrorNotLoadable        LoadDataError = 2
	LoadDataErrorDatasetSizeTooLarge LoadDataError = 3
	LoadDataErrorNotAwaitedSegment  LoadDataError = 4
	LoadDataErrorInterpretationFailure LoadDataError = 5
	LoadDataErrorStorageFailure     LoadDataError = 6
	LoadDataErrorDatasetNotReady    LoadDataError = 7
)

// DataScopeError represents data scope error types
type DataScopeError uint8

const (
	DataScopeErrorOther DataScopeError = 0
)

// TaskError represents task error types
type TaskError uint8

const (
	TaskErrorOther         TaskError = 0
	TaskErrorNoRemoteControl TaskError = 1
	TaskErrorTIStopped     TaskError = 2
	TaskErrorTIRunning     TaskError = 3
	TaskErrorTIUnusable    TaskError = 4
)

// OtherError represents other error types
type OtherError uint8

const (
	OtherErrorOther OtherError = 0
)

// CosemInterface represents COSEM interface class identifiers
type CosemInterface uint8

const (
	// Parameters and measurement data
	CosemInterfaceData            CosemInterface = 1
	CosemInterfaceRegister        CosemInterface = 3
	CosemInterfaceExtendedRegister CosemInterface = 4
	CosemInterfaceDemandRegister CosemInterface = 5
	CosemInterfaceRegisterActivation CosemInterface = 6
	CosemInterfaceProfileGeneric  CosemInterface = 7
	CosemInterfaceUtilityTables   CosemInterface = 26
	CosemInterfaceRegisterTable   CosemInterface = 61
	CosemInterfaceCompactData    CosemInterface = 62
	CosemInterfaceStatusMapping   CosemInterface = 63

	// Access control and management
	CosemInterfaceAssociationSN   CosemInterface = 12
	CosemInterfaceAssociationLN   CosemInterface = 15
	CosemInterfaceSAPAssignment  CosemInterface = 17
	CosemInterfaceImageTransfer  CosemInterface = 18
	CosemInterfaceSecuritySetup  CosemInterface = 64
	CosemInterfacePush            CosemInterface = 40
	CosemInterfaceCosemDataProtection CosemInterface = 30
	CosemInterfaceFunctionControl CosemInterface = 122
	CosemInterfaceArrayManager    CosemInterface = 123
	CosemInterfaceCommunicationPortProtection CosemInterface = 124

	// Time and event bound control
	CosemInterfaceClock            CosemInterface = 8
	CosemInterfaceScriptTable     CosemInterface = 9
	CosemInterfaceSchedule        CosemInterface = 10
	CosemInterfaceSpecialDaysTable CosemInterface = 11
	CosemInterfaceActivityCalendar CosemInterface = 20
	CosemInterfaceRegisterMonitor  CosemInterface = 21
	CosemInterfaceSingleActionSchedule CosemInterface = 22
	CosemInterfaceDisconnectControl CosemInterface = 70
	CosemInterfaceLimiter         CosemInterface = 71
	CosemInterfaceParameterMonitor CosemInterface = 65
	CosemInterfaceSensorManager   CosemInterface = 67
	CosemInterfaceArbitrator      CosemInterface = 68

	// Payment related interfaces
	CosemInterfaceAccount      CosemInterface = 111
	CosemInterfaceCredit      CosemInterface = 112
	CosemInterfaceCharge      CosemInterface = 113
	CosemInterfaceTokenGateway CosemInterface = 115

	// Data exchange over local ports and modems
	CosemInterfaceIECLocalPortSetup CosemInterface = 19
	CosemInterfaceIECHDLCSetup     CosemInterface = 23
	CosemInterfaceIECTwistedPairSetup CosemInterface = 24
	CosemInterfaceModemConfiguration CosemInterface = 27
	CosemInterfaceAutoAnswer        CosemInterface = 28
	CosemInterfaceAutoConnect       CosemInterface = 29
	CosemInterfaceGPRSModemSetup    CosemInterface = 45
	CosemInterfaceGSMDiagnostics    CosemInterface = 47
	CosemInterfaceLTEMonitoring     CosemInterface = 151

	// Data exchange over M-Bus
	CosemInterfaceMBusSlavePortSetup CosemInterface = 25
	CosemInterfaceMBusClient         CosemInterface = 72
	CosemInterfaceMBusWirelessModeQChannel CosemInterface = 73
	CosemInterfaceMBusMasterPortSetup CosemInterface = 74
	CosemInterfaceMBusPortSetupDlmsCosemServer CosemInterface = 76
	CosemInterfaceMBusDiagnostics    CosemInterface = 77

	// Data exchange over Internet
	CosemInterfaceTCPUDPSetup    CosemInterface = 41
	CosemInterfaceIPv4Setup      CosemInterface = 42
	CosemInterfaceIPv6Setup      CosemInterface = 48
	CosemInterfaceMACAddressSetup CosemInterface = 43
	CosemInterfacePPPSetup       CosemInterface = 44
	CosemInterfaceSMTPSetup      CosemInterface = 46
	CosemInterfaceNTPSetup       CosemInterface = 100

	// Data exchange using S-FSK PLC
	CosemInterfaceSFSKPhyMacSetup CosemInterface = 50
	CosemInterfaceSFSKActiveInitiator CosemInterface = 51
	CosemInterfaceSFSKMacSynchronisationTimeouts CosemInterface = 52
	CosemInterfaceSFSKMacCounters CosemInterface = 53
	CosemInterfaceSFSKIEC61334432LLCSetup CosemInterface = 55
	CosemInterfaceSFSKReportingSystemList CosemInterface = 56

	// LLC layers for IEC 8802-2
	CosemInterfaceIEC88022LLCType1Setup CosemInterface = 57
	CosemInterfaceIEC88022LLCType2Setup CosemInterface = 58
	CosemInterfaceIEC88022LLCType3Setup CosemInterface = 59

	// Narrowband OFDM PLC profile for PRIME networks
	CosemInterfacePrime61344432LLCSSCSSetup CosemInterface = 80
	CosemInterfacePrimeOFDMPLCPhysicalLayerCounters CosemInterface = 81
	CosemInterfacePrimeOFDMPLCMACSetup CosemInterface = 82
	CosemInterfacePrimeOFDMPLCMACFunctionalParameters CosemInterface = 83
	CosemInterfacePrimeOFDMPLCMACCounters CosemInterface = 84
	CosemInterfacePrimeOFDMPLCMACNetworkAdministrationData CosemInterface = 85
	CosemInterfacePrimeOFDMPLCMACApplicationIdentification CosemInterface = 86

	// Narrowband OFDM PLC profile for G3-PLC network
	CosemInterfaceG3PLCMACLayerCounters CosemInterface = 90
	CosemInterfaceG3PLCMACSetup        CosemInterface = 91
	CosemInterfaceG3PLC6LowpanAdaptationLayerSetup CosemInterface = 92

	// HS-PLC IEC 12139-1
	CosemInterfaceHSPLCIEC121391MACSetup CosemInterface = 140
	CosemInterfaceHSPLCIEC121391CPASSetup CosemInterface = 141
	CosemInterfaceHSPLCIEC121391IPSSASSetup CosemInterface = 142
	CosemInterfaceHSPLCIEC121391HDLCSSASSetup CosemInterface = 143

	// Zigbee
	CosemInterfaceZigbeeSASStartup CosemInterface = 101
	CosemInterfaceZigbeeSASJoin   CosemInterface = 102
	CosemInterfaceZigbeeSASAPSFragmentation CosemInterface = 103
	CosemInterfaceZigbeeNetworkControl CosemInterface = 104
	CosemInterfaceZigbeeTunnelSetup CosemInterface = 105

	// LPWAN networks
	CosemInterfaceSCHCLPWAN        CosemInterface = 126
	CosemInterfaceSCHCLPWANDiagnostics CosemInterface = 127
	CosemInterfaceLoRaWANSetup    CosemInterface = 128
	CosemInterfaceLoRaWANDiagnostics CosemInterface = 129

	// Wi-SUN
	CosemInterfaceWiSUNSetup       CosemInterface = 95
	CosemInterfaceWiSUMDiagnostics  CosemInterface = 96
	CosemInterfaceRPLDiagnostics   CosemInterface = 97
	CosemInterfaceMPLDiagnostics   CosemInterface = 98

	// IEC 14908 PLC
	CosemInterfaceIEC14908Identification CosemInterface = 130
	CosemInterfaceIEC14908ProtocolSetup  CosemInterface = 131
	CosemInterfaceIEC14908ProtocolStatus CosemInterface = 132
	CosemInterfaceIEC14908Diagnostics    CosemInterface = 133
)

// ReleaseRequestReason represents release request reason
type ReleaseRequestReason uint8

const (
	ReleaseRequestReasonNormal     ReleaseRequestReason = 0
	ReleaseRequestReasonUrgent     ReleaseRequestReason = 1
	ReleaseRequestReasonUserDefined ReleaseRequestReason = 30
)

// ReleaseResponseReason represents release response reason
type ReleaseResponseReason uint8

const (
	ReleaseResponseReasonNormal     ReleaseResponseReason = 0
	ReleaseResponseReasonNotFinished ReleaseResponseReason = 1
	ReleaseResponseReasonUserDefined ReleaseResponseReason = 30
)

// AuthenticationMechanism represents authentication mechanism types
type AuthenticationMechanism uint8

const (
	AuthenticationMechanismNone    AuthenticationMechanism = 0
	AuthenticationMechanismLLS    AuthenticationMechanism = 1
	AuthenticationMechanismHLS     AuthenticationMechanism = 2
	AuthenticationMechanismHLSMD5  AuthenticationMechanism = 3 // Insecure. Don't use with new meters
	AuthenticationMechanismHLSSHA1 AuthenticationMechanism = 4 // Insecure. Don't use with new meters
	AuthenticationMechanismHLSGMAC AuthenticationMechanism = 5
	AuthenticationMechanismHLSSHA256 AuthenticationMechanism = 6
	AuthenticationMechanismHLSECDSA AuthenticationMechanism = 7
)

// AcseServiceUserDiagnostics represents ACSE service user diagnostics
type AcseServiceUserDiagnostics uint8

const (
	AcseServiceUserDiagnosticsNull AcseServiceUserDiagnostics = 0
	AcseServiceUserDiagnosticsNoReasonGiven AcseServiceUserDiagnostics = 1
	AcseServiceUserDiagnosticsApplicationContextNameNotSupported AcseServiceUserDiagnostics = 2
	AcseServiceUserDiagnosticsCallingAPTitleNotRecognized AcseServiceUserDiagnostics = 3
	AcseServiceUserDiagnosticsCallingAPInvocationIdentifierNotRecognized AcseServiceUserDiagnostics = 4
	AcseServiceUserDiagnosticsCallingAEQualifierNotRecognized AcseServiceUserDiagnostics = 5
	AcseServiceUserDiagnosticsCallingAEInvocationIdentifierNotRecognized AcseServiceUserDiagnostics = 6
	AcseServiceUserDiagnosticsCalledAPTitleNotRecognized AcseServiceUserDiagnostics = 7
	AcseServiceUserDiagnosticsCalledAPInvocationIdentifierNotRecognized AcseServiceUserDiagnostics = 8
	AcseServiceUserDiagnosticsCalledAEQualifierNotRecognized AcseServiceUserDiagnostics = 9
	AcseServiceUserDiagnosticsCalledAEInvocationIdentifierNotRecognized AcseServiceUserDiagnostics = 10
	AcseServiceUserDiagnosticsAuthenticationMechanismNameNotRecognized AcseServiceUserDiagnostics = 11
	AcseServiceUserDiagnosticsAuthenticationMechanismNameRequired AcseServiceUserDiagnostics = 12
	AcseServiceUserDiagnosticsAuthenticationFailed AcseServiceUserDiagnostics = 13
	AcseServiceUserDiagnosticsAuthenticationRequired AcseServiceUserDiagnostics = 14
)

// AcseServiceProviderDiagnostics represents ACSE service provider diagnostics
type AcseServiceProviderDiagnostics uint8

const (
	AcseServiceProviderDiagnosticsNull AcseServiceProviderDiagnostics = 0
	AcseServiceProviderDiagnosticsNoReasonGiven AcseServiceProviderDiagnostics = 1
	AcseServiceProviderDiagnosticsNoCommonACSEVersion AcseServiceProviderDiagnostics = 2
)

// AssociationResult represents association result
type AssociationResult uint8

const (
	AssociationResultAccepted         AssociationResult = 0
	AssociationResultRejectedPermanent AssociationResult = 1
	AssociationResultRejectedTransient AssociationResult = 2
)

// ActionResultStatus represents action result status
type ActionResultStatus uint8

const (
	ActionResultStatusSuccess                ActionResultStatus = 0
	ActionResultStatusHardwareFault         ActionResultStatus = 1
	ActionResultStatusTemporaryFailure       ActionResultStatus = 2
	ActionResultStatusReadWriteDenied        ActionResultStatus = 3
	ActionResultStatusObjectUndefined        ActionResultStatus = 4
	ActionResultStatusObjectClassInconsistent ActionResultStatus = 9
	ActionResultStatusObjectUnavailable      ActionResultStatus = 11
	ActionResultStatusTypeUnmatched          ActionResultStatus = 12
	ActionResultStatusScopeOfAccessViolated  ActionResultStatus = 13
	ActionResultStatusDataBlockUnavailable   ActionResultStatus = 14
	ActionResultStatusLongActionAborted      ActionResultStatus = 15
	ActionResultStatusNoLongActionInProgress ActionResultStatus = 16
	ActionResultStatusOtherReason            ActionResultStatus = 250
)

