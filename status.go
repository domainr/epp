package epp

// Status represents EPP status codes as a bitfield.
// https://www.icann.org/resources/pages/epp-status-codes-2014-06-16-en
// https://tools.ietf.org/html/std69
// https://tools.ietf.org/html/rfc3915
type Status uint32

// Status types, in order of priority, low to high.
// Status are stored in a single integer as a bit field.
const (
	StatusUnknown Status = iota

	// Server status codes set by a domain registry
	StatusOK Status = 1 << (iota - 1) // Standard status for a domain, meaning it has no pending operations or prohibitions.
	StatusLinked
	StatusAddPeriod // This grace period is provided after the initial registration of a domain name. If the registrar deletes the domain name during this period, the registry may provide credit to the registrar for the cost of the registration.
	StatusAutoRenewPeriod
	StatusInactive
	StatusPendingCreate
	StatusPendingDelete
	StatusPendingRenew
	StatusPendingRestore
	StatusPendingTransfer
	StatusPendingUpdate
	StatusRedemptionPeriod
	StatusRenewPeriod
	StatusServerDeleteProhibited
	StatusServerHold
	StatusServerRenewProhibited
	StatusServerTransferProhibited
	StatusServerUpdateProhibited
	StatusTransferPeriod
	StatusClientDeleteProhibited
	StatusClientHold
	StatusClientRenewProhibited
	StatusClientTransferProhibited
	StatusClientUpdateProhibited

	// RDAP status codes map roughly, but not exactly to EPP status codes.
	// https://tools.ietf.org/html/rfc8056#section-2
	StatusActive     = StatusOK
	StatusAssociated = StatusLinked

	// StatusClient are status codes set by a domain registrar.
	StatusClient = StatusClientDeleteProhibited | StatusClientHold | StatusClientRenewProhibited | StatusClientTransferProhibited | StatusClientUpdateProhibited
)

// stringToStatus maps EPP and RDAP status strings to Status bits.
var stringToStatus = map[string]Status{
	"ok":                         StatusOK,
	"active":                     StatusActive,
	"linked":                     StatusLinked,
	"associated":                 StatusAssociated,
	"add period":                 StatusAddPeriod,
	"addPeriod":                  StatusAddPeriod,
	"auto renew period":          StatusAutoRenewPeriod,
	"autoRenewPeriod":            StatusAutoRenewPeriod,
	"inactive":                   StatusInactive,
	"pending create":             StatusPendingCreate,
	"pendingCreate":              StatusPendingCreate,
	"pending delete":             StatusPendingDelete,
	"pendingDelete":              StatusPendingDelete,
	"pending renew":              StatusPendingRenew,
	"pendingRenew":               StatusPendingRenew,
	"pending restore":            StatusPendingRestore,
	"pendingRestore":             StatusPendingRestore,
	"pending transfer":           StatusPendingTransfer,
	"pendingTransfer":            StatusPendingTransfer,
	"pending update":             StatusPendingUpdate,
	"pendingUpdate":              StatusPendingUpdate,
	"redemption period":          StatusRedemptionPeriod,
	"redemptionPeriod":           StatusRedemptionPeriod,
	"renew period":               StatusRenewPeriod,
	"renewPeriod":                StatusRenewPeriod,
	"server delete prohibited":   StatusServerDeleteProhibited,
	"serverDeleteProhibited":     StatusServerDeleteProhibited,
	"server hold":                StatusServerHold,
	"serverHold":                 StatusServerHold,
	"server renew prohibited":    StatusServerRenewProhibited,
	"serverRenewProhibited":      StatusServerRenewProhibited,
	"server transfer prohibited": StatusServerTransferProhibited,
	"serverTransferProhibited":   StatusServerTransferProhibited,
	"server update prohibited":   StatusServerUpdateProhibited,
	"serverUpdateProhibited":     StatusServerUpdateProhibited,
	"transfer period":            StatusTransferPeriod,
	"transferPeriod":             StatusTransferPeriod,
	"client delete prohibited":   StatusClientDeleteProhibited,
	"clientDeleteProhibited":     StatusClientDeleteProhibited,
	"client hold":                StatusClientHold,
	"clientHold":                 StatusClientHold,
	"client renew prohibited":    StatusClientRenewProhibited,
	"clientRenewProhibited":      StatusClientRenewProhibited,
	"client transfer prohibited": StatusClientTransferProhibited,
	"clientTransferProhibited":   StatusClientTransferProhibited,
	"client update prohibited":   StatusClientUpdateProhibited,
	"clientUpdateProhibited":     StatusClientUpdateProhibited,
}

// ParseStatus returns a Status from one or more strings.
// It does not attempt to validate the input or resolve conflicting status bits.
func ParseStatus(in ...string) Status {
	var s Status
	for _, v := range in {
		s |= stringToStatus[v]
	}
	return s
}
