package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	epp "github.com/onasunnymorning/eppclient"
	"github.com/wsxiaoys/terminal/color"
)

var (
	profileName string
	verbose     bool
)

func main() {
	// Global flags
	flag.StringVar(&profileName, "profile", "default", "profile name in ~/.epp/credentials")
	flag.BoolVar(&verbose, "v", false, "enable verbose debug logging")

	// Capture logs
	var logBuf bytes.Buffer

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command> [arguments]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  check   Check domain availability\n")
		fmt.Fprintf(os.Stderr, "  create  Create a domain\n")
		fmt.Fprintf(os.Stderr, "  renew   Renew a domain\n")
		fmt.Fprintf(os.Stderr, "  restore Restore a domain (RGP)\n")
		fmt.Fprintf(os.Stderr, "  transfer Transfer a domain\n")
		fmt.Fprintf(os.Stderr, "  raw     Send raw XML from a file or stdin\n")
		fmt.Fprintf(os.Stderr, "  info    Get domain info\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	// We need to parse flags before subcommands to get -profile and -v
	// But flag.Parse() consumes args. So we need to be careful.
	// Actually, standard go flag pkg stops at non-flag args.
	// So `epp -profile foo check domain.com` works.
	// `epp check -profile foo` does NOT work for global flags if we do it this way.
	// But usually subcommands have their own flags.

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Manual flag parsing for global flags before subcommand
	// This is a bit hacky but allows `epp -v check ...`
	// A better way is to parse, then look at remaining args.
	flag.Parse()

	if verbose {
		epp.DebugLogger = io.MultiWriter(os.Stderr, &logBuf)
	} else {
		epp.DebugLogger = &logBuf
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	cmd := args[0]
	subArgs := args[1:]

	// Check usage before connecting
	checkUsage(cmd, subArgs)

	cfg, err := loadConfig(profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading credentials for profile %q: %v\n", profileName, err)
		os.Exit(1)
	}

	conn := connect(cfg)

	defer func() {
		if r := recover(); r != nil {
			// If it was our fatal error, we already logged it.
			// Just ensure we prompt and then exit 1.
			// We can't easily distinguish our panic from others unless we use a custom type,
			// but for this CLI it's probably fine to treat all panics as "something went wrong".
			// But we DO want to run promptRawXML.
			conn.Close()
			promptRawXML(&logBuf)
			os.Exit(1)
		}
		conn.Close()
		promptRawXML(&logBuf)
	}()

	switch cmd {
	case "check":
		runCheck(conn, subArgs)
	case "info":
		runInfo(conn, subArgs)
	case "delete":
		runDelete(conn, subArgs)
	case "create":
		runCreate(conn, subArgs)
	case "renew":
		runRenew(conn, subArgs)
	case "restore":
		runRestore(conn, subArgs)
	case "transfer":
		runTransfer(conn, subArgs)
	case "raw":
		runRaw(conn, subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func checkUsage(cmd string, args []string) {
	switch cmd {
	case "check":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp check <domain>...")
			os.Exit(1)
		}
	case "info":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp info <domain|contact> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" && sub != "contact" {
			fmt.Fprintf(os.Stderr, "Unknown info type: %s. Use 'domain' or 'contact'.\n", sub)
			os.Exit(1)
		}
	case "delete":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp delete <domain|contact> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" && sub != "contact" {
			fmt.Fprintf(os.Stderr, "Unknown delete type: %s. Use 'domain' or 'contact'.\n", sub)
			os.Exit(1)
		}
	case "create":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp create <domain|contact> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" && sub != "contact" {
			fmt.Fprintf(os.Stderr, "Unknown create type: %s. Use 'domain' or 'contact'.\n", sub)
			os.Exit(1)
		}
	case "renew":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp renew <domain> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" {
			fmt.Fprintf(os.Stderr, "Unknown renewal type: %s. Use 'domain'.\n", sub)
			os.Exit(1)
		}
	case "restore":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp restore <domain> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" {
			fmt.Fprintf(os.Stderr, "Unknown restore type: %s. Use 'domain'.\n", sub)
			os.Exit(1)
		}
	case "transfer":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp transfer <domain> [options]")
			os.Exit(1)
		}
		sub := args[0]
		if sub != "domain" {
			fmt.Fprintf(os.Stderr, "Unknown transfer type: %s. Use 'domain'.\n", sub)
			os.Exit(1)
		}
	case "raw":
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: epp raw <file>")
			os.Exit(1)
		}
	}
}

func connect(cfg *Config) *epp.Conn {
	// Set up TLS
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
	}

	host, _, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		host = cfg.Addr
	}
	tlsCfg.ServerName = host

	if cfg.CACert != "" {
		ca, err := ioutil.ReadFile(cfg.CACert)
		fatalif(err)
		tlsCfg.RootCAs = x509.NewCertPool()
		tlsCfg.RootCAs.AppendCertsFromPEM(ca)
	}

	if cfg.Cert != "" && cfg.Key != "" {
		crt, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
		fatalif(err)
		tlsCfg.Certificates = append(tlsCfg.Certificates, crt)
	}

	if !cfg.TLS {
		tlsCfg = nil
	}

	var conn net.Conn
	// TODO: Proxy support if needed from config

	color.Fprintf(os.Stderr, "Connecting to %s\n", cfg.Addr)
	conn, err = net.Dial("tcp", cfg.Addr)
	fatalif(err)

	if tlsCfg != nil {
		color.Fprintf(os.Stderr, "Establishing TLS connection\n")
		tc := tls.Client(conn, tlsCfg)
		err = tc.Handshake()
		fatalif(err)
		conn = tc
	}

	color.Fprintf(os.Stderr, "Performing EPP handshake\n")
	c, err := epp.NewConn(conn)
	fatalif(err)

	color.Fprintf(os.Stderr, "Logging in as %s...\n", cfg.User)
	err = c.Login(cfg.User, cfg.Password, "")
	fatalif(err)

	return c
}

func runCheck(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp check <domain>...")
		os.Exit(1)
	}

	start := time.Now()
	dc, err := c.CheckDomain(args...)
	logif(err)
	printDCR(dc)
	qdur := time.Since(start)
	color.Fprintf(os.Stderr, "@{.}Query: %s\n", qdur)
}

func runInfo(c *epp.Conn, args []string) {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runInfoDomain(c, subArgs)
	case "contact":
		runInfoContact(c, subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown info type: %s. Use 'domain' or 'contact'.\n", cmd)
		os.Exit(1)
	}
}

func runInfoDomain(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp info domain <domain>")
		os.Exit(1)
	}
	res, err := c.DomainInfo(args[0], nil)
	fatalif(err)

	fmt.Printf("Domain: %s\n", res.Domain)
	fmt.Printf("ROID: %s\n", res.ID)
	fmt.Printf("Status: %v\n", res.Status)
	fmt.Printf("Created: %s\n", res.CrDate)
	fmt.Printf("Expires: %s\n", res.ExDate)
}

func runInfoContact(c *epp.Conn, args []string) {
	fs := flag.NewFlagSet("info contact", flag.ExitOnError)
	auth := fs.String("auth", "", "auth info (required for contact info)")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp info contact [-auth code] <contact-id>")
		os.Exit(1)
	}

	res, err := c.ContactInfo(fs.Arg(0), *auth, nil)
	fatalif(err)

	fmt.Printf("Contact: %s\n", res.ID)
	fmt.Printf("ROID: %s\n", res.ROID)
	fmt.Printf("Status: %v\n", res.Status)
	fmt.Printf("Email: %s\n", res.Email)
	fmt.Printf("Created: %s\n", res.CrDate)
}

func runDelete(c *epp.Conn, args []string) {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runDeleteDomain(c, subArgs)
	case "contact":
		runDeleteContact(c, subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown delete type: %s. Use 'domain' or 'contact'.\n", cmd)
		os.Exit(1)
	}
}

func runDeleteDomain(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp delete domain <domain>")
		os.Exit(1)
	}
	err := c.DeleteDomain(args[0], nil)
	fatalif(err)
	color.Printf("@{g}Domain %s deleted!\n", args[0])
}

func runDeleteContact(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp delete contact <contact-id>")
		os.Exit(1)
	}
	err := c.DeleteContact(args[0], nil)
	fatalif(err)
	color.Printf("@{g}Contact %s deleted!\n", args[0])
}

func runCreate(c *epp.Conn, args []string) {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runCreateDomain(c, subArgs)
	case "contact":
		runCreateContact(c, subArgs)
	default:
		// Fallback for backward compatibility or simple "epp create domain.com"?
		// The user explicitly asked for "move to epp create domain", so enforcing subcommand is correct.
		// However, it's nice to allow "epp create domain.com" if it doesn't look like "contact"?
		// But "domain" IS the subcommand.
		// If I type `epp create example.com` -> cmd="example.com".
		// Maybe warning? Or just fail.
		// User: "lets move 'epp create' to 'epp create domain'"
		fmt.Fprintf(os.Stderr, "Unknown create type: %s. Use 'domain' or 'contact'.\n", cmd)
		os.Exit(1)
	}
}

func runCreateDomain(c *epp.Conn, args []string) {
	fs := flag.NewFlagSet("create domain", flag.ExitOnError)
	period := fs.Int("period", 1, "registration period in years")
	auth := fs.String("auth", "", "auth info")
	registrant := fs.String("registrant", "", "registrant contact ID")
	admin := fs.String("contact-admin", "", "admin contact ID")
	tech := fs.String("contact-tech", "", "tech contact ID")
	billing := fs.String("contact-billing", "", "billing contact ID")

	nsParams := fs.String("ns", "", "comma separated nameservers")
	fee := fs.String("fee", "", "fee amount (requires -currency usually)")
	currency := fs.String("currency", "", "fee currency")

	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp create domain [-period N] [-auth code] [-registrant id] ... <domain>")
		os.Exit(1)
	}

	domain := fs.Arg(0)

	contacts := make(map[string]string)
	if *admin != "" {
		contacts["admin"] = *admin
	}
	if *tech != "" {
		contacts["tech"] = *tech
	}
	if *billing != "" {
		contacts["billing"] = *billing
	}

	var ns []string
	if *nsParams != "" {
		ns = strings.Split(*nsParams, ",")
		for i, n := range ns {
			ns[i] = strings.TrimSpace(n)
		}
	}

	var extData map[string]string
	if *fee != "" {
		extData = make(map[string]string)
		extData["fee:fee"] = *fee
		if *currency != "" {
			extData["fee:currency"] = *currency
		}
	}

	res, err := c.CreateDomain(domain, *period, "y", *auth, *registrant, contacts, ns, extData)
	fatalif(err)
	color.Printf("@{g}Domain %s created!\nCreated: %s\nExpiry: %s\n", res.Domain, res.CrDate, res.ExDate)
}

func runCreateContact(c *epp.Conn, args []string) {
	fs := flag.NewFlagSet("create contact", flag.ExitOnError)
	id := fs.String("id", "", "contact ID")
	email := fs.String("email", "", "email address")
	name := fs.String("name", "", "contact name")
	org := fs.String("org", "", "organization")
	street := fs.String("street", "", "street address")
	city := fs.String("city", "", "city")
	sp := fs.String("sp", "", "state/province")
	pc := fs.String("pc", "", "postal code")
	cc := fs.String("cc", "", "country code")
	voice := fs.String("voice", "", "voice phone number")
	auth := fs.String("auth", "", "auth info")

	fs.Parse(args)

	if *id == "" || *email == "" || *name == "" || *city == "" || *cc == "" || *auth == "" {
		fmt.Fprintln(os.Stderr, "Usage: epp create contact -id <id> -email <email> -name <name> -city <city> -cc <cc> -auth <auth> [-voice number] [options]")
		fs.PrintDefaults()
		os.Exit(1)
	}

	pi := epp.PostalInfo{
		Name:   *name,
		Org:    *org,
		Street: *street,
		City:   *city,
		SP:     *sp,
		PC:     *pc,
		CC:     *cc,
	}

	res, err := c.CreateContact(*id, *email, pi, *voice, *auth, nil)
	fatalif(err)
	color.Printf("@{g}Contact %s created!\nCreated: %s\n", res.ID, res.CrDate)
}

func runRenew(c *epp.Conn, args []string) {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runRenewDomain(c, subArgs)
	default:
		// Fallback or error?
		// Since "epp renew domain" is what we want, we should enforce it?
		// But existing "epp renew" was convenient.
		// "epp renew domain.com" was the old way.
		// If cmd doesn't look like a subcommand (no "domain" keyword), maybe assume old behavior?
		// But cleaner to be strict if we are standardizing.
		// "unknown command %s".
		fmt.Fprintf(os.Stderr, "Unknown renewal type: %s. Use 'domain'.\n", cmd)
		os.Exit(1)
	}
}

func runRenewDomain(c *epp.Conn, args []string) {
	fs := flag.NewFlagSet("renew domain", flag.ExitOnError)
	period := fs.Int("period", 1, "renewal period in years")
	curExp := fs.String("exp", "", "current expiry date (YYYY-MM-DD) - optional, will be fetched if not provided")
	fee := fs.String("fee", "", "fee amount")
	currency := fs.String("currency", "", "fee currency")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp renew domain [-exp YYYY-MM-DD] [-period N] [-fee amount] [-currency code] <domain>")
		os.Exit(1)
	}

	domain := fs.Arg(0)
	var date time.Time
	var err error

	if *curExp != "" {
		date, err = time.Parse("2006-01-02", *curExp)
		fatalif(err)
	} else {
		// Auto-fetch expiry if not provided
		fmt.Printf("Fetching info for %s to determine current expiry date...\n", domain)
		infoRes, err := c.DomainInfo(domain, nil)
		fatalif(err)
		date = infoRes.ExDate
		fmt.Printf("Current expiry: %s\n", date.Format("2006-01-02"))
	}

	var extData map[string]string
	if *fee != "" {
		extData = make(map[string]string)
		extData["fee:fee"] = *fee
		if *currency != "" {
			extData["fee:currency"] = *currency
		}
	}

	res, err := c.RenewDomain(domain, date, *period, "y", extData)
	fatalif(err)
	color.Printf("@{g}Domain %s renewed!\nNew Expiry: %s\n", res.Domain, res.ExDate)
}

func runRestore(c *epp.Conn, args []string) {
	cmd := args[0]
	subArgs := args[1:]

	switch cmd {
	case "domain":
		runRestoreDomain(c, subArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown restore type: %s. Use 'domain'.\n", cmd)
		os.Exit(1)
	}
}

func runRestoreDomain(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp restore domain <domain>")
		os.Exit(1)
	}
	// Restore often requires RGP extension
	// TODO: Add support for reporting data if required by registry?
	// For now, simple RGP restore request.
	_, err := c.RestoreDomain(args[0], nil)
	fatalif(err)
	color.Printf("@{g}Domain %s restored!\n", args[0])
}

func runTransfer(c *epp.Conn, args []string) {
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
	op := fs.String("op", "query", "transfer operation (query, request, approve, reject, cancel)")
	auth := fs.String("auth", "", "auth info")
	period := fs.Int("period", 1, "registration period in years (optional for request)")
	fee := fs.String("fee", "", "fee amount")
	currency := fs.String("currency", "", "fee currency")
	fs.Parse(args)

	if fs.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp transfer domain [-op op] [-auth code] [-period N] [-fee amount] [-currency code] <domain>")
		os.Exit(1)
	}

	domain := fs.Arg(0)

	var extData map[string]string
	if *fee != "" {
		extData = make(map[string]string)
		extData["fee:fee"] = *fee
		if *currency != "" {
			extData["fee:currency"] = *currency
		}
	}

	res, err := c.TransferDomain(*op, domain, *period, "y", *auth, extData)
	fatalif(err)

	color.Printf("@{g}Domain %s %s operation successful!\n", domain, *op)
	if res != nil {
		fmt.Printf("Status: %s\n", res.Status)
		if !res.REDate.IsZero() {
			fmt.Printf("Requested: %s by %s\n", res.REDate.Format(time.RFC3339), res.REID)
		}
		if !res.ACDate.IsZero() {
			fmt.Printf("Acted: %s by %s\n", res.ACDate.Format(time.RFC3339), res.ACID)
		}
		if !res.ExDate.IsZero() {
			fmt.Printf("Expiry: %s\n", res.ExDate.Format(time.RFC3339))
		}
	}
}

func runRaw(c *epp.Conn, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: epp raw <file>")
		os.Exit(1)
	}

	var data []byte
	var err error

	filename := args[0]
	if filename == "-" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(filename)
	}
	fatalif(err)

	res, err := c.Raw(data)
	fatalif(err)

	fmt.Printf("%s\n", string(res))
}

func logif(err error) bool {
	if err != nil {
		color.Fprintf(os.Stderr, "@{r}%s\n", err)
		return true
	}
	return false
}

func fatalif(err error) {
	if logif(err) {
		// Panic ensuring we can recover in main to show logs if needed
		panic(err)
	}
}

func printDCR(dcr *epp.DomainCheckResponse) {
	if dcr == nil {
		return
	}
	av := make(map[string]bool)
	for _, c := range dcr.Checks {
		av[c.Domain] = c.Available
		if c.Available {
			color.Printf("@{g}%s\tavail=%t\treason=%q\n", c.Domain, c.Available, c.Reason)
		} else {
			color.Printf("@{y}%s\tavail=%t\treason=%q\n", c.Domain, c.Available, c.Reason)
		}
	}
	for _, c := range dcr.Charges {
		if av[c.Domain] {
			color.Printf("@{g}%s\tcategory=%s\tname=%q\tcreate=%s\trenew=%s\trestore=%s\ttransfer=%s\tcurrency=%s\n", c.Domain, c.Category, c.CategoryName, c.CreatePrice, c.RenewPrice, c.RestorePrice, c.TransferPrice, c.Currency)
		} else {
			color.Printf("@{y}%s\tcategory=%s\tname=%q\tcreate=%s\trenew=%s\trestore=%s\ttransfer=%s\tcurrency=%s\n", c.Domain, c.Category, c.CategoryName, c.CreatePrice, c.RenewPrice, c.RestorePrice, c.TransferPrice, c.Currency)
		}
	}
}
