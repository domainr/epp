package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/domainr/epp"
	"github.com/wsxiaoys/terminal/color"
)

func main() {
	var cmd, obj, uri, addr, user, pass, proxyAddr, crtPath, caPath, keyPath string
	var useTLS, batch, verbose bool

	flag.StringVar(&cmd, "cmd", "check", "EPP Command")
	flag.StringVar(&obj, "obj", "domain", "EPP object type (i.e. domain, contact, host etc.")
	flag.StringVar(&uri, "url", "", "EPP server URL, e.g. epp://user:pass@api.1api.net:700")
	flag.StringVar(&addr, "addr", "", "EPP server address (HOST:PORT)")
	flag.StringVar(&user, "u", "", "EPP user name")
	flag.StringVar(&pass, "p", "", "EPP password")
	flag.BoolVar(&useTLS, "tls", true, "use TLS")
	flag.StringVar(&proxyAddr, "proxy", "", "SOCKS5 proxy address (HOST:PORT)")
	flag.StringVar(&crtPath, "cert", "", "path to SSL certificate")
	flag.StringVar(&keyPath, "key", "", "path to SSL private key")
	flag.StringVar(&caPath, "ca", "", "path to SSL certificate authority")
	flag.BoolVar(&batch, "batch", false, "check all domains in a single EPP command")
	flag.BoolVar(&verbose, "v", false, "enable verbose debug logging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <query>\n\nAvailable arguments:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
	}

	if verbose {
		epp.DebugLogger = os.Stderr
	}

	objs := make([]string, len(flag.Args()))
	for i, arg := range flag.Args() {
		objs[i] = arg // FIXME: convert unicode to Punycode?
	}

	// Parse URL
	if uri != "" {
		addr, user, pass = parseURL(uri)
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}

	// Set up TLS
	cfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Load certificates
	if caPath != "" {
		color.Fprintf(os.Stderr, "Loading CA certificate from %s\n", caPath)
		ca, err := ioutil.ReadFile(caPath)
		fatalif(err)
		cfg.RootCAs = x509.NewCertPool()
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}

	if crtPath != "" && keyPath != "" {
		color.Fprintf(os.Stderr, "Loading certificate %s and key %s\n", crtPath, keyPath)
		crt, err := tls.LoadX509KeyPair(crtPath, keyPath)
		fatalif(err)
		cfg.Certificates = append(cfg.Certificates, crt)
		// cfg.BuildNameToCertificate()
		useTLS = true
	}

	// Use TLS?
	if !useTLS {
		cfg = nil
	}

	// Dial
	start := time.Now()
	var conn net.Conn
	if proxyAddr != "" {
		color.Fprintf(os.Stderr, "Connecting to %s via proxy %s\n", addr, proxyAddr)
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, &net.Dialer{})
		fatalif(err)
		conn, err = dialer.Dial("tcp", addr)
	} else {
		color.Fprintf(os.Stderr, "Connecting to %s\n", addr)
		conn, err = net.Dial("tcp", addr)
	}
	fatalif(err)

	// TLS
	if useTLS {
		color.Fprintf(os.Stderr, "Establishing TLS connection\n")
		tc := tls.Client(conn, cfg)
		err = tc.Handshake()
		fatalif(err)
		conn = tc
	}

	// EPP
	color.Fprintf(os.Stderr, "Performing EPP handshake\n")
	c, err := epp.NewConn(conn)
	fatalif(err)
	color.Fprintf(os.Stderr, "Logging in as %s...\n", user)
	err = c.Login(user, pass, "")
	fatalif(err)

	// Check
	start = time.Now()
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "check":
		handleCheck(c, obj, objs, batch)
	case "info":
		handleInfo(c, obj, objs)
	default:
		log.Fatal("Unknown command:", cmd)
	}
	qdur := time.Since(start)

	color.Fprintf(os.Stderr, "@{.}Query: %s Avg: %s\n", qdur, qdur/time.Duration(len(objs)))
}

func handleCheck(c *epp.Conn, objKind string, objs []string, batch bool) error {
	objKind = strings.ToLower(objKind)
	switch objKind {
	case "domain":
		if batch {
			dc, err := c.CheckDomain(objs...)
			logif(err)
			printDCR(dc)
			return err
		} else {
			for _, domain := range objs {
				dc, err := c.CheckDomain(domain)
				logif(err)
				printDCR(dc)
				return err
			}
		}
	default:
		return errors.New("Unknown epp object type: " + objKind)
	}
	return nil
}

func handleInfo(c *epp.Conn, objKind string, objs []string) error {
	objKind = strings.ToLower(objKind)
	switch objKind {
	case "domain":
		for _, domain := range objs {
			res, err := c.DomainInfo(domain)
			logif(err)
			printDIR(domain, res)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("Unknown epp object type: " + objKind)
	}
	return nil
}

func parseURL(uri string) (addr, user, pass string) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	host, port, err := net.SplitHostPort(u.Host)
	if host == "" {
		host = u.Host
	}
	if port == "" {
		port = DefaultEPPPort
	}
	addr = net.JoinHostPort(host, port)
	if ui := u.User; ui != nil {
		user = ui.Username()
		pass, _ = ui.Password()
	}
	return
}

// DefaultEPPPort is the default TCP port for the EPP protocol.
const DefaultEPPPort = "700"

func logif(err error) bool {
	if err != nil {
		color.Fprintf(os.Stderr, "@{r}%s\n", err)
		return true
	}
	return false
}

func fatalif(err error) {
	if logif(err) {
		color.Fprintf(os.Stderr, "@{r}EXITING\n")
		os.Exit(1)
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
			color.Printf("@{g}%s\tcategory=%s\tname=%q\n", c.Domain, c.Category, c.CategoryName)
		} else {
			color.Printf("@{y}%s\tcategory=%s\tname=%q\n", c.Domain, c.Category, c.CategoryName)
		}
	}
}

func printDIR(name string, dir *epp.DomainInfoResponse) {
	if dir == nil {
		return
	}

	if dir.Result.Code == 1000 {
		color.Printf("@{g}domain=%s\tcreated=%s\texpires=%s\n", dir.Name, dir.CreatedDate, dir.Expiration)
	} else {
		color.Printf("@{y}domain=%s\tcode=%d\tmsg=%s\n", name, dir.Result.Code, dir.Result.Message)
	}
}
