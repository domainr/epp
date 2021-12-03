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
	case SuccessEndingSession:
		return "Command completed successfully; ending session"
	case ErrUnknownCommand:
		return "Unknown command"
	case ErrCommandSyntaxError:
		return "Command syntax error"
	case ErrCommandUseError:
		return "Command use error"
	case ErrRequiredParameterMissing:
		return "Required parameter missing"
	case ErrParameterValueRangeError:
		return "Parameter value range error"
	case ErrParameterValueSyntaxError:
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
	case ErrObjectNotEligbleForRenewal:
		return "Object is not eligible for renewal"
	case ErrObjectNotEligibleForTransfer:
		return "Object is not eligible for transfer"
	case ErrAuthenticationError:
		return "Authentication error"
	case ErrAuthorizationError:
		return "Authorization error"
	case ErrInvalidAuthorization:
		return "Invalid authorization information"
	case ErrObjectPendingTransfer:
		return "Object pending transfer"
	case ErrObjectNotPendingTransfer:
		return "Object not pending transfer"
	case ErrObjectExists:
		return "Object exists"
	case ErrObjectDoesNotExist:
		return "Object does not exist"
	case ErrObjectStatusProhibitsOperation:
		return "Object status prohibits operation"
	case ErrObjectAssociationProhibitsOperation:
		return "Object association prohibits operation"
	case ErrParameterValuePolicyError:
		return "Parameter value policy error"
	case ErrUnimplementedObject:
		return "Unimplemented object service"
	case ErrDataManagementViolation:
		return "Data management policy violation"
	case ErrCommandFailed:
		return "Command failed"
	case ErrCommandFailedClosing:
		return "Command failed; server closing connection"
	case ErrAuthenticationErrorClosing:
		return "Authentication error; server closing connection"
	case ErrSessionLimitExceeded:
		return "Session limit exceeded; server closing connection"
	default:
		return fmt.Sprintf("Err code %04d", c)
	}
}

const (
	ResultCodeMin ResultCode = 1000
	ResultCodeMax ResultCode = 2599

	// This should match the number of known result codes below
	KnownResultCodes = 34

	// Success result codes
	Success              ResultCode = 1000
	SuccessPending       ResultCode = 1001
	SuccessNoMessages    ResultCode = 1300
	SuccessAck           ResultCode = 1301
	SuccessEndingSession ResultCode = 1500

	// Error result codes
	ErrUnknownCommand                      ResultCode = 2000
	ErrCommandSyntaxError                  ResultCode = 2001
	ErrCommandUseError                     ResultCode = 2002
	ErrRequiredParameterMissing            ResultCode = 2003
	ErrParameterValueRangeError            ResultCode = 2004
	ErrParameterValueSyntaxError           ResultCode = 2005
	ErrUnimplementedVersion                ResultCode = 2100
	ErrUnimplementedCommand                ResultCode = 2101
	ErrUnimplementedOption                 ResultCode = 2102
	ErrUnimplementedExtension              ResultCode = 2103
	ErrBillingFailure                      ResultCode = 2104
	ErrObjectNotEligbleForRenewal          ResultCode = 2105
	ErrObjectNotEligibleForTransfer        ResultCode = 2106
	ErrAuthenticationError                 ResultCode = 2200
	ErrAuthorizationError                  ResultCode = 2201
	ErrInvalidAuthorization                ResultCode = 2202
	ErrObjectPendingTransfer               ResultCode = 2300
	ErrObjectNotPendingTransfer            ResultCode = 2301
	ErrObjectExists                        ResultCode = 2302
	ErrObjectDoesNotExist                  ResultCode = 2303
	ErrObjectStatusProhibitsOperation      ResultCode = 2304
	ErrObjectAssociationProhibitsOperation ResultCode = 2305
	ErrParameterValuePolicyError           ResultCode = 2306
	ErrUnimplementedObject                 ResultCode = 2307
	ErrDataManagementViolation             ResultCode = 2308
	ErrCommandFailed                       ResultCode = 2400
	ErrCommandFailedClosing                ResultCode = 2500
	ErrAuthenticationErrorClosing          ResultCode = 2501
	ErrSessionLimitExceeded                ResultCode = 2502
)
