package repository

import (
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/ssh"
	"telee/pkg/telnet"
	"time"

	x "github.com/google/goexpect"
)

const (
	noSuffix string = ""
	haSuffix string = "/pri/act"
)

// Repository struct
type Repository struct {
	Config *config.Config
}

// Fetch returns stdout from telnet session
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher
	var data string
	var err error

	if r.Config.HAMode {
		expects = r.buildPrivilegedRequest(haSuffix)
	}
	if !r.Config.HAMode {
		expects = r.buildPrivilegedRequest(noSuffix)
	}

	if r.Config.SecureMode {
		data, err = ssh.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects, ssh.GenerateClientConfig(r.Config.Username, r.Config.Password))
	} else {
		data, err = telnet.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects)
	}

	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: asa] buildPrivilegedRequest returns the expects
func (r *Repository) buildPrivilegedRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: "terminal pager 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}
}
