package epp

import "fmt"

// ResultCode represents a 4-digit EPP result code.
// See https://tools.ietf.org/rfcmarkup?doc=5730#section-3.
// A ResultCode can be used as an error value.
// Note: only result codes >= 2000 are considered errors.
type ResultCode uint16

// Message returns a Message representation of c.
func (c ResultCode) Message() Message {
	return Message{Lang: "en", Value: c.String()}
}

// MarshalText implements encoding.TextMarshaler to print c as a 4-digit number.
func (c ResultCode) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%04d", c)), nil
}

// IsError returns true if c represents an error code (>= 2000).
func (c ResultCode) IsError() bool {
	return c >= 2000
}

// IsFatal returns true if c represents an error code that closes the
// connection.
func (c ResultCode) IsFatal() bool {
	return c >= 2500
}

// Error returns the text representation of c if c is an error, or an empty
// string if c is a successful result code.
func (c ResultCode) Error() string {
	if c.IsError() {
		return c.String()
	}
	return ""
}

// String returns the English text representation of c.
func (c ResultCode) String() string {
	switch c {
	case Success:
		return "Command completed successfully"
	case SuccessPending:
		return "Command completed successfully; action pending"
	case SuccessNoMessages:
		return "Command completed successfully; no messages"
	case SuccessAck:
		return "Command completed successfully; ack to dequeue"
	case SuccessEnd:
		return "Command completed successfully; ending session"
	case ErrUnknownCommand:
		return "Unknown command"
	case ErrCommandSyntax:
		return "Command syntax error"
	case ErrCommandUse:
		return "Command use error"
	case ErrRequiredParameter:
		return "Required parameter missing"
	case ErrParameterRange:
		return "Parameter value range error"
	case ErrParameterSyntax:
		return "Parameter value syntax error"
	case ErrUnimplementedVersion:
		return "Unimplemented protocol version"
	case ErrUnimplementedCommand:
		return "Unimplemented command"
	case ErrUnimplementedOption:
		return "Unimplemented option"
	case ErrUnimplementedExtension:
		return "Unimplemented extension"
	case ErrBillingFailure:
		return "Billing failure"
	case ErrNotEligbleForRenewal:
		return "Object is not eligible for renewal"
	case ErrNotEligibleForTransfer:
		return "Object is not eligible for transfer"
	case ErrAuthentication:
		return "Authentication error"
	case ErrAuthorization:
		return "Authorization error"
	case ErrInvalidAuthorization:
		return "Invalid authorization information"
	case ErrPendingTransfer:
		return "Object pending transfer"
	case ErrNotPendingTransfer:
		return "Object not pending transfer"
	case ErrExists:
		return "Object exists"
	case ErrDoesNotExist:
		return "Object does not exist"
	case ErrStatus:
		return "Object status prohibits operation"
	case ErrAssociation:
		return "Object association prohibits operation"
	case ErrParameterPolicy:
		return "Parameter value policy error"
	case ErrUnimplementedObject:
		return "Unimplemented object service"
	case ErrDataManagementViolation:
		return "Data management policy violation"
	case ErrCommandFailed:
		return "Command failed"
	case ErrCommandFailedClosing:
		return "Command failed; server closing connection"
	case ErrAuthenticationClosing:
		return "Authentication error; server closing connection"
	case ErrSessionLimitExceeded:
		return "Session limit exceeded; server closing connection"
	default:
		return fmt.Sprintf("Status code %04d", c)
	}
}

const (
	ResultCodeMin ResultCode = 1000
	ResultCodeMax ResultCode = 2599

	// This should match the number of known result codes below
	KnownResultCodes = 34

	// Success result codes
	Success           ResultCode = 1000
	SuccessPending    ResultCode = 1001
	SuccessNoMessages ResultCode = 1300
	SuccessAck        ResultCode = 1301
	SuccessEnd        ResultCode = 1500

	// Error result codes
	ErrUnknownCommand          ResultCode = 2000
	ErrCommandSyntax           ResultCode = 2001
	ErrCommandUse              ResultCode = 2002
	ErrRequiredParameter       ResultCode = 2003
	ErrParameterRange          ResultCode = 2004
	ErrParameterSyntax         ResultCode = 2005
	ErrUnimplementedVersion    ResultCode = 2100
	ErrUnimplementedCommand    ResultCode = 2101
	ErrUnimplementedOption     ResultCode = 2102
	ErrUnimplementedExtension  ResultCode = 2103
	ErrBillingFailure          ResultCode = 2104
	ErrNotEligbleForRenewal    ResultCode = 2105
	ErrNotEligibleForTransfer  ResultCode = 2106
	ErrAuthentication          ResultCode = 2200
	ErrAuthorization           ResultCode = 2201
	ErrInvalidAuthorization    ResultCode = 2202
	ErrPendingTransfer         ResultCode = 2300
	ErrNotPendingTransfer      ResultCode = 2301
	ErrExists                  ResultCode = 2302
	ErrDoesNotExist            ResultCode = 2303
	ErrStatus                  ResultCode = 2304
	ErrAssociation             ResultCode = 2305
	ErrParameterPolicy         ResultCode = 2306
	ErrUnimplementedObject     ResultCode = 2307
	ErrDataManagementViolation ResultCode = 2308
	ErrCommandFailed           ResultCode = 2400
	ErrCommandFailedClosing    ResultCode = 2500
	ErrAuthenticationClosing   ResultCode = 2501
	ErrSessionLimitExceeded    ResultCode = 2502
)
