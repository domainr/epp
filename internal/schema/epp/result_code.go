package epp

import "fmt"

// ResultCode represents a 4-digit EPP result code.
// See https://tools.ietf.org/rfcmarkup?doc=5730#section-3.
// A ResultCode can be used as an error value. Note: only result codes >= 2000
// are considered errors.
type ResultCode uint16

// MarshalText implements encoding.TextMarshaler to print c as a 4-digit number.
func (c ResultCode) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%04d", c)), nil
}

// IsFatal returns true if c represents an error code that closes the
// connection.
func (c ResultCode) IsFatal() bool {
	return c >= 2500
}

// IsError returns true if c represents an error code (>= 2000).
func (c ResultCode) IsError() bool {
	return c >= 2000
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
	case ResultSuccess:
		return "Command completed successfully"
	case ResultSuccessPending:
		return "Command completed successfully; action pending"
	case ResultSuccessNoMessages:
		return "Command completed successfully; no messages"
	case ResultSuccessAck:
		return "Command completed successfully; ack to dequeue"
	case ResultSuccessEndingSession:
		return "Command completed successfully; ending session"
	case ResultUnknownCommand:
		return "Unknown command"
	case ResultCommandSyntaxError:
		return "Command syntax error"
	case ResultCommandUseError:
		return "Command use error"
	case ResultRequiredParameterMissing:
		return "Required parameter missing"
	case ResultParameterValueRangeError:
		return "Parameter value range error"
	case ResultParameterValueSyntaxError:
		return "Parameter value syntax error"
	case ResultUnimplementedVersion:
		return "Unimplemented protocol version"
	case ResultUnimplementedCommand:
		return "Unimplemented command"
	case ResultUnimplementedOption:
		return "Unimplemented option"
	case ResultUnimplementedExtension:
		return "Unimplemented extension"
	case ResultBillingFailure:
		return "Billing failure"
	case ResultObjectNotEligbleForRenewal:
		return "Object is not eligible for renewal"
	case ResultObjectNotEligibleForTransfer:
		return "Object is not eligible for transfer"
	case ResultAuthenticationError:
		return "Authentication error"
	case ResultAuthorizationError:
		return "Authorization error"
	case ResultInvalidAuthorization:
		return "Invalid authorization information"
	case ResultObjectPendingTransfer:
		return "Object pending transfer"
	case ResultObjectNotPendingTransfer:
		return "Object not pending transfer"
	case ResultObjectExists:
		return "Object exists"
	case ResultObjectDoesNotExist:
		return "Object does not exist"
	case ResultObjectStatusProhibitsOperation:
		return "Object status prohibits operation"
	case ResultObjectAssociationProhibitsOperation:
		return "Object association prohibits operation"
	case ResultParameterValuePolicyError:
		return "Parameter value policy error"
	case ResultUnimplementedObject:
		return "Unimplemented object service"
	case ResultDataManagementViolation:
		return "Data management policy violation"
	case ResultCommandFailed:
		return "Command failed"
	case ResultCommandFailedClosing:
		return "Command failed; server closing connection"
	case ResultAuthenticationErrorClosing:
		return "Authentication error; server closing connection"
	case ResultSessionLimitExceeded:
		return "Session limit exceeded; server closing connection"
	default:
		return fmt.Sprintf("Result code %04d", c)
	}
}

const (
	// Success result codes
	ResultSuccess              ResultCode = 1000
	ResultSuccessPending       ResultCode = 1001
	ResultSuccessNoMessages    ResultCode = 1300
	ResultSuccessAck           ResultCode = 1301
	ResultSuccessEndingSession ResultCode = 1500

	// Error result codes
	ResultUnknownCommand                      ResultCode = 2000
	ResultCommandSyntaxError                  ResultCode = 2001
	ResultCommandUseError                     ResultCode = 2002
	ResultRequiredParameterMissing            ResultCode = 2003
	ResultParameterValueRangeError            ResultCode = 2004
	ResultParameterValueSyntaxError           ResultCode = 2005
	ResultUnimplementedVersion                ResultCode = 2100
	ResultUnimplementedCommand                ResultCode = 2101
	ResultUnimplementedOption                 ResultCode = 2102
	ResultUnimplementedExtension              ResultCode = 2103
	ResultBillingFailure                      ResultCode = 2104
	ResultObjectNotEligbleForRenewal          ResultCode = 2105
	ResultObjectNotEligibleForTransfer        ResultCode = 2106
	ResultAuthenticationError                 ResultCode = 2200
	ResultAuthorizationError                  ResultCode = 2201
	ResultInvalidAuthorization                ResultCode = 2202
	ResultObjectPendingTransfer               ResultCode = 2300
	ResultObjectNotPendingTransfer            ResultCode = 2301
	ResultObjectExists                        ResultCode = 2302
	ResultObjectDoesNotExist                  ResultCode = 2303
	ResultObjectStatusProhibitsOperation      ResultCode = 2304
	ResultObjectAssociationProhibitsOperation ResultCode = 2305
	ResultParameterValuePolicyError           ResultCode = 2306
	ResultUnimplementedObject                 ResultCode = 2307
	ResultDataManagementViolation             ResultCode = 2308
	ResultCommandFailed                       ResultCode = 2400
	ResultCommandFailedClosing                ResultCode = 2500
	ResultAuthenticationErrorClosing          ResultCode = 2501
	ResultSessionLimitExceeded                ResultCode = 2502
)
