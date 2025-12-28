package exceptions

import "fmt"

// LocalDlmsProtocolError represents a protocol error
type LocalDlmsProtocolError struct {
	Message string
}

func (e *LocalDlmsProtocolError) Error() string {
	return fmt.Sprintf("DLMS protocol error: %s", e.Message)
}

// NewLocalDlmsProtocolError creates a new LocalDlmsProtocolError
func NewLocalDlmsProtocolError(message string) *LocalDlmsProtocolError {
	return &LocalDlmsProtocolError{Message: message}
}

// ApplicationAssociationError represents an error when trying to setup the application association
type ApplicationAssociationError struct {
	Message string
}

func (e *ApplicationAssociationError) Error() string {
	return fmt.Sprintf("Application association error: %s", e.Message)
}

// NewApplicationAssociationError creates a new ApplicationAssociationError
func NewApplicationAssociationError(message string) *ApplicationAssociationError {
	return &ApplicationAssociationError{Message: message}
}

// PreEstablishedAssociationError represents an error when doing illegal things to a pre-established connection
type PreEstablishedAssociationError struct {
	Message string
}

func (e *PreEstablishedAssociationError) Error() string {
	return fmt.Sprintf("Pre-established association error: %s", e.Message)
}

// NewPreEstablishedAssociationError creates a new PreEstablishedAssociationError
func NewPreEstablishedAssociationError(message string) *PreEstablishedAssociationError {
	return &PreEstablishedAssociationError{Message: message}
}

// ConformanceError represents an error when APDUs do not match connection conformance
type ConformanceError struct {
	Message string
}

func (e *ConformanceError) Error() string {
	return fmt.Sprintf("Conformance error: %s", e.Message)
}

// NewConformanceError creates a new ConformanceError
func NewConformanceError(message string) *ConformanceError {
	return &ConformanceError{Message: message}
}

// CipheringError represents an error when ciphering or deciphering an APDU
type CipheringError struct {
	Message string
}

func (e *CipheringError) Error() string {
	return fmt.Sprintf("Ciphering error: %s", e.Message)
}

// NewCipheringError creates a new CipheringError
func NewCipheringError(message string) *CipheringError {
	return &CipheringError{Message: message}
}

// DlmsClientException represents an exception relating to the client
type DlmsClientException struct {
	Message string
}

func (e *DlmsClientException) Error() string {
	return fmt.Sprintf("DLMS client exception: %s", e.Message)
}

// NewDlmsClientException creates a new DlmsClientException
func NewDlmsClientException(message string) *DlmsClientException {
	return &DlmsClientException{Message: message}
}

// CommunicationError represents an error in communication with a meter
type CommunicationError struct {
	Message string
}

func (e *CommunicationError) Error() string {
	return fmt.Sprintf("Communication error: %s", e.Message)
}

// NewCommunicationError creates a new CommunicationError
func NewCommunicationError(message string) *CommunicationError {
	return &CommunicationError{Message: message}
}

// CryptographyError represents an error when applying a cryptographic function
type CryptographyError struct {
	Message string
}

func (e *CryptographyError) Error() string {
	return fmt.Sprintf("Cryptography error: %s", e.Message)
}

// NewCryptographyError creates a new CryptographyError
func NewCryptographyError(message string) *CryptographyError {
	return &CryptographyError{Message: message}
}

// DecryptionError represents an error when unable to decrypt an APDU
// It can be due to mismatch in authentication tag because the ciphertext has changed
// or that the key, nonce or associated data is wrong
type DecryptionError struct {
	Message string
}

func (e *DecryptionError) Error() string {
	return fmt.Sprintf("Decryption error: %s", e.Message)
}

// NewDecryptionError creates a new DecryptionError
func NewDecryptionError(message string) *DecryptionError {
	return &DecryptionError{Message: message}
}

// NoRlrqRlreError is raised from connection when a ReleaseRequest is issued
// on a connection that has use_rlrq_rlre==False
// Control for client to just skip Release and disconnect the lower layer.
type NoRlrqRlreError struct {
	Message string
}

func (e *NoRlrqRlreError) Error() string {
	return fmt.Sprintf("No RLRQ/RLRE error: %s", e.Message)
}

// NewNoRlrqRlreError creates a new NoRlrqRlreError
func NewNoRlrqRlreError(message string) *NoRlrqRlreError {
	return &NoRlrqRlreError{Message: message}
}

