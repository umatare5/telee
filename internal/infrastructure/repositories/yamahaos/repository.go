// Package repository implements YAMAHA RT OS-specific data access layer.
package repository

import (
	"time"

	x "github.com/google/goexpect"
	cryptossh "golang.org/x/crypto/ssh"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/ssh"
	"github.com/umatare5/telee/pkg/telnet"
)

const (
	promptPassword   string = "Password:"
	cmdDisablePaging string = "console lines infinity\n"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the output the last prompt match captured.
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher
	var data string
	var err error

	if r.Config.SecureMode {
		if r.Config.EnableMode {
			expects = r.buildPrivilegedSecureRequest()
		}
		if !r.Config.EnableMode {
			expects = r.buildUserModeSecureRequest()
		}
	}
	if !r.Config.SecureMode {
		if r.Config.EnableMode {
			expects = r.buildPrivilegedRequest()
		}
		if !r.Config.EnableMode {
			expects = r.buildUserModeRequest()
		}
	}

	if r.Config.SecureMode {
		var clientConfig *cryptossh.ClientConfig
		clientConfig, err = ssh.GenerateClientConfig(r.Config.Username, r.Config.Password, r.Config.HostKeyPath, r.Config.Hostname)
		if err != nil {
			return "", err
		}
		data, err = ssh.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects, clientConfig)
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

// YAMAHA prompts for a password only, so --username reaches the device over SSH alone.
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
	}
}

func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "administrator\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
	}
}

func (r *Repository) buildUserModeSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
	}
}

func (r *Repository) buildPrivilegedSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "administrator\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
	}
}
