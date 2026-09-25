// Package repository implements Cisco ASA Software-specific data access layer.
package repository

import (
	"time"

	cryptossh "golang.org/x/crypto/ssh"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	x "github.com/umatare5/telee/pkg/expect"
	"github.com/umatare5/telee/pkg/ssh"
	"github.com/umatare5/telee/pkg/telnet"
)

const (
	promptUsername   string = "Username:"
	promptPassword   string = "Password:"
	cmdDisablePaging string = "terminal pager 0\n"
	noSuffix         string = ""
	haSuffix         string = "/pri/act"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the transcript of the commands.
func (r *Repository) Fetch() (string, error) {
	var login, commands []x.Batcher
	var data string
	var err error

	if r.Config.SecureMode {
		if r.Config.DefaultPrivMode && r.Config.RedundantMode {
			login, commands = r.buildDefaultPrivilegedSecureRequest(haSuffix)
		}
		if !r.Config.DefaultPrivMode && r.Config.RedundantMode {
			login, commands = r.buildPrivilegedSecureRequest(haSuffix)
		}
		if r.Config.DefaultPrivMode && !r.Config.RedundantMode {
			login, commands = r.buildDefaultPrivilegedSecureRequest(noSuffix)
		}
		if !r.Config.DefaultPrivMode && !r.Config.RedundantMode {
			login, commands = r.buildPrivilegedSecureRequest(noSuffix)
		}
	} else {
		if r.Config.DefaultPrivMode && r.Config.RedundantMode {
			login, commands = r.buildDefaultPrivilegedRequest(haSuffix)
		}
		if !r.Config.DefaultPrivMode && r.Config.RedundantMode {
			login, commands = r.buildPrivilegedRequest(haSuffix)
		}
		if r.Config.DefaultPrivMode && !r.Config.RedundantMode {
			login, commands = r.buildDefaultPrivilegedRequest(noSuffix)
		}
		if !r.Config.DefaultPrivMode && !r.Config.RedundantMode {
			login, commands = r.buildPrivilegedRequest(noSuffix)
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
		).Fetch(login, commands, clientConfig)
	} else {
		data, err = telnet.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(login, commands)
	}

	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *Repository) buildPrivilegedRequest(suffix string) (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: promptUsername},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}, x.Commands(r.Config.Hostname+suffix+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildDefaultPrivilegedRequest(suffix string) (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: promptUsername},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}, x.Commands(r.Config.Hostname+suffix+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildPrivilegedSecureRequest(suffix string) (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + suffix + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}, x.Commands(r.Config.Hostname+suffix+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildDefaultPrivilegedSecureRequest(suffix string) (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}, x.Commands(r.Config.Hostname+suffix+"#", "\n", r.Config.Commands)
}
