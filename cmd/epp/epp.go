package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"

	"github.com/domainr/epp"
	"github.com/wsxiaoys/terminal/color"
)

func main() {
	var uri, addr, user, pass, proxyAddr, crtPath, caPath, keyPath string
	var useTLS, batch, verbose bool

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

	domains := make([]string, len(flag.Args()))
	for i, arg := range flag.Args() {
		domains[i] = arg // FIXME: convert unicode to Punycode?
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
	if batch {
		dc, err := c.CheckDomain(domains...)
		logif(err)
		printDCR(dc)
	} else {
		for _, domain := range domains {
			dc, err := c.CheckDomain(domain)
			logif(err)
			printDCR(dc)
		}
	}
	qdur := time.Since(start)

	color.Fprintf(os.Stderr, "@{.}Query: %s Avg: %s\n", qdur, qdur/time.Duration(len(domains)))
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
