package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/domainr/epp"
)

func runTransfer(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp transfer domain <request|query|approve|reject|cancel> [options] <domain>")
		os.Exit(1)
	}

	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runTransferDomain(c, subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown transfer type: %s. Use 'domain'.\n", cmd)
		os.Exit(1)
	}
}

func runTransferDomain(c *epp.Conn, args []string) {
	fs := flag.NewFlagSet("transfer domain", flag.ExitOnError)
	auth := fs.String("auth", "", "auth info (often required for transfer request)")
	periodYears := fs.Int("period", 0, "transfer period in years (optional; usually for request/approve)")
	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: epp transfer domain <request|query|approve|reject|cancel> [-auth code] [-period N] <domain>")
		os.Exit(1)
	}

	opStr := fs.Arg(0)
	domain := fs.Arg(1)

	var op epp.TransferOp
	switch opStr {
	case "request":
		op = epp.TransferRequest
	case "query":
		op = epp.TransferQuery
	case "approve":
		op = epp.TransferApprove
	case "reject":
		op = epp.TransferReject
	case "cancel":
		op = epp.TransferCancel
	default:
		fmt.Fprintf(os.Stderr, "Unknown transfer op: %s. Use request|query|approve|reject|cancel.\n", opStr)
		os.Exit(1)
	}

	var period *epp.Period
	if *periodYears > 0 {
		period = &epp.Period{Value: *periodYears, Unit: "y"}
	}

	res, err := c.TransferDomain(op, domain, *auth, period, nil)
	fatalif(err)

	fmt.Printf("Domain: %s\n", res.Name)
	fmt.Printf("Status: %s\n", res.TrStatus)
	if res.ReID != "" {
		fmt.Printf("RequestedBy: %s\n", res.ReID)
	}
	if !res.ReDate.IsZero() {
		fmt.Printf("RequestedAt: %s\n", res.ReDate)
	}
	if res.AcID != "" {
		fmt.Printf("ActionBy: %s\n", res.AcID)
	}
	if !res.AcDate.IsZero() {
		fmt.Printf("ActionAt: %s\n", res.AcDate)
	}
	if res.ExDate != nil && !res.ExDate.IsZero() {
		fmt.Printf("Expires: %s\n", res.ExDate)
	}
}



