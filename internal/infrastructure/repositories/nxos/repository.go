// Package repository implements Cisco NX-OS-specific data access layer.
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
	promptLogin      string = "login:"
	promptPassword   string = "Password:"
	cmdDisablePaging string = "terminal length 0\n"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the output the last prompt match captured.
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher
	var data string
	var err error

	if r.Config.SecureMode && r.Config.DefaultPrivMode {
		expects = r.buildDefaultPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && r.Config.EnableMode {
		expects = r.buildPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
		expects = r.buildUserModeSecureRequest()
	}
	if !r.Config.SecureMode && r.Config.DefaultPrivMode {
		expects = r.buildDefaultPrivilegedRequest()
	}
	if !r.Config.SecureMode && r.Config.EnableMode {
		expects = r.buildPrivilegedRequest()
	}
	if !r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
		expects = r.buildUserModeRequest()
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

// NX-OS echoes one space after the prompt character, so every BExp here ends with it.
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: promptLogin},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
	}
}

func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: promptLogin},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

func (r *Repository) buildDefaultPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: promptLogin},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

func (r *Repository) buildUserModeSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
	}
}

func (r *Repository) buildPrivilegedSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

func (r *Repository) buildDefaultPrivilegedSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}
