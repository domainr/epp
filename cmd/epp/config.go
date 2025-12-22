package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	User     string
	Password string
	Addr     string
	TLS      bool
	Cert     string
	Key      string
	CACert   string
}

func loadConfig(profile string) (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(usr.HomeDir, ".epp", "credentials")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open credentials file: %w", err)
	}
	defer f.Close()

	cfg := &Config{}
	scanner := bufio.NewScanner(f)
	currentSection := ""

	// Basic INI parser
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			continue
		}

		if currentSection != profile {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "user":
			cfg.User = val
		case "password":
			cfg.Password = val
		case "addr":
			cfg.Addr = val
		case "tls":
			cfg.TLS = (val == "true")
		case "cert":
			cfg.Cert = val
		case "key":
			cfg.Key = val
		case "ca":
			cfg.CACert = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if cfg.User == "" {
		return nil, fmt.Errorf("user not found in profile %s", profile)
	}

	return cfg, nil
}
