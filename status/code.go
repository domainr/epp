package epp

// Code represents EPP status codes as a bitfield.
// See https://tools.ietf.org/html/std69, https://tools.ietf.org/html/rfc3915,
// and https://www.icann.org/resources/pages/epp-status-codes-2014-06-16-en.
type Code uint32

// EPP status codes, in order of priority, from low to high. Codes are stored in
// a single integer as a bit field.
const (
	Unknown Code = iota

	// OK is the default status code for a domain, meaning it has no pending operations or
	// prohibitions.
	OK Code = 1 << (iota - 1)
	Linked

	// This grace period is provided after the initial registration of a
	// domain name. If the registrar deletes the domain name during this
	// period, the registry may provide credit to the registrar for the cost
	// of the registration.
	AddPeriod
	AutoRenewPeriod
	Inactive
	PendingCreate
	PendingDelete
	PendingRenew
	PendingRestore
	PendingTransfer
	PendingUpdate
	RedemptionPeriod
	RenewPeriod
	ServerDeleteProhibited
	ServerHold
	ServerRenewProhibited
	ServerTransferProhibited
	ServerUpdateProhibited
	TransferPeriod
	ClientDeleteProhibited
	ClientHold
	ClientRenewProhibited
	ClientTransferProhibited
	ClientUpdateProhibited

	// RDAP status codes loosely map to EPP status codes.
	// See https://tools.ietf.org/html/rfc8056#section-2.
	Active     = OK
	Associated = Linked

	// ClientCodes are status codes set by a domain registrar.
	ClientCodes = ClientDeleteProhibited |
		ClientHold |
		ClientRenewProhibited |
		ClientTransferProhibited |
		ClientUpdateProhibited
)

// stringToCode maps EPP and RDAP status strings to Codes.
var stringToCode = map[string]Code{
	"ok":                         OK,
	"active":                     Active,
	"linked":                     Linked,
	"associated":                 Associated,
	"add period":                 AddPeriod,
	"addPeriod":                  AddPeriod,
	"auto renew period":          AutoRenewPeriod,
	"autoRenewPeriod":            AutoRenewPeriod,
	"inactive":                   Inactive,
	"pending create":             PendingCreate,
	"pendingCreate":              PendingCreate,
	"pending delete":             PendingDelete,
	"pendingDelete":              PendingDelete,
	"pending renew":              PendingRenew,
	"pendingRenew":               PendingRenew,
	"pending restore":            PendingRestore,
	"pendingRestore":             PendingRestore,
	"pending transfer":           PendingTransfer,
	"pendingTransfer":            PendingTransfer,
	"pending update":             PendingUpdate,
	"pendingUpdate":              PendingUpdate,
	"redemption period":          RedemptionPeriod,
	"redemptionPeriod":           RedemptionPeriod,
	"renew period":               RenewPeriod,
	"renewPeriod":                RenewPeriod,
	"server delete prohibited":   ServerDeleteProhibited,
	"serverDeleteProhibited":     ServerDeleteProhibited,
	"server hold":                ServerHold,
	"serverHold":                 ServerHold,
	"server renew prohibited":    ServerRenewProhibited,
	"serverRenewProhibited":      ServerRenewProhibited,
	"server transfer prohibited": ServerTransferProhibited,
	"serverTransferProhibited":   ServerTransferProhibited,
	"server update prohibited":   ServerUpdateProhibited,
	"serverUpdateProhibited":     ServerUpdateProhibited,
	"transfer period":            TransferPeriod,
	"transferPeriod":             TransferPeriod,
	"client delete prohibited":   ClientDeleteProhibited,
	"clientDeleteProhibited":     ClientDeleteProhibited,
	"client hold":                ClientHold,
	"clientHold":                 ClientHold,
	"client renew prohibited":    ClientRenewProhibited,
	"clientRenewProhibited":      ClientRenewProhibited,
	"client transfer prohibited": ClientTransferProhibited,
	"clientTransferProhibited":   ClientTransferProhibited,
	"client update prohibited":   ClientUpdateProhibited,
	"clientUpdateProhibited":     ClientUpdateProhibited,
}

// Parse returns a status code from one or more strings. It does not attempt to
// validate the input or resolve conflicting status bits.
func Parse(in ...string) Code {
	var s Code
	for _, v := range in {
		s |= stringToCode[v]
	}
	return s
}
