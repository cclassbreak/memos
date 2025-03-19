package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Profile is the configuration to start main server.
type Profile struct {
	// Mode can be "prod" or "dev" or "demo"
	Mode string
	// Addr is the binding address for server
	Addr string
	// Port is the binding port for server
	Port int
	// Data is the data directory
	Data string
	// DSN points to where memos stores its own data
	DSN string
	// Driver is the database driver
	// sqlite, mysql
	Driver string
	// Version is the current version of server
	Version string
	// InstanceURL is the url of your memos instance.
	InstanceURL string
}

func (p *Profile) IsDev() bool {
	return p.Mode != "prod"
}

func checkDataDir(dataDir string) (string, error) {
	// Convert to absolute path if relative path is supplied.
	if !filepath.IsAbs(dataDir) {
		relativeDir := filepath.Join(filepath.Dir(os.Args[0]), dataDir)
		absDir, err := filepath.Abs(relativeDir)
		if err != nil {
			return "", err
		}
		dataDir = absDir
	}

	// Trim trailing \ or / in case user supplies
	dataDir = strings.TrimRight(dataDir, "\\/")
	if _, err := os.Stat(dataDir); err != nil {
		return "", errors.Wrapf(err, "unable to access data folder %s", dataDir)
	}
	return dataDir, nil
}

func (p *Profile) Validate() error {
	// Ensure mode is valid
	if p.Mode != "demo" && p.Mode != "dev" && p.Mode != "prod" {
		p.Mode = "demo"
	}

	// Ensure address is always localhost
	if p.Addr == "" {
		p.Addr = "127.0.0.1"
	}

	// Ensure Data is explicitly set in config.yaml, default to current folder if empty
	if p.Data == "" {
		p.Data = "."
	}

	// Ensure set port to 8081 as this is required from frontend dev proxy server if mode is dev
	if p.IsDev() {
		p.Port = 8081
	}

	// Set DSN for SQLite if not provided
	if p.Driver == "sqlite" && p.DSN == "" {
		dbFile := fmt.Sprintf("memos_%s.db", p.Mode)
		p.DSN = filepath.Join(p.Data, dbFile)
	}

	return nil
}
