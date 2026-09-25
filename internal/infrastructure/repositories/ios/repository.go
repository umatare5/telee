// Package repository implements Cisco IOS/IOS-XE-specific data access layer.
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
	cmdDisablePaging string = "terminal length 0\n"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the transcript of the commands.
func (r *Repository) Fetch() (string, error) {
	var login, commands []x.Batcher
	var data string
	var err error

	if r.Config.SecureMode && r.Config.DefaultPrivMode {
		login, commands = r.buildDefaultPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && r.Config.EnableMode {
		login, commands = r.buildPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
		login, commands = r.buildUserModeSecureRequest()
	}
	if !r.Config.SecureMode && r.Config.DefaultPrivMode {
		login, commands = r.buildDefaultPrivilegedRequest()
	}
	if !r.Config.SecureMode && r.Config.EnableMode {
		login, commands = r.buildPrivilegedRequest()
	}
	if !r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
		login, commands = r.buildUserModeRequest()
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

func (r *Repository) buildUserModeRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: promptUsername},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + ">"},
	}, x.Commands(r.Config.Hostname+">", "\n", r.Config.Commands)
}

func (r *Repository) buildPrivilegedRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: promptUsername},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
	}, x.Commands(r.Config.Hostname+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildDefaultPrivilegedRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: promptUsername},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
	}, x.Commands(r.Config.Hostname+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildUserModeSecureRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + ">"},
	}, x.Commands(r.Config.Hostname+">", "\n", r.Config.Commands)
}

func (r *Repository) buildPrivilegedSecureRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: promptPassword},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
	}, x.Commands(r.Config.Hostname+"#", "\n", r.Config.Commands)
}

func (r *Repository) buildDefaultPrivilegedSecureRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: cmdDisablePaging},
		&x.BExp{R: r.Config.Hostname + "#"},
	}, x.Commands(r.Config.Hostname+"#", "\n", r.Config.Commands)
}
