// Package repository implements Cisco AireOS-specific data access layer.
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
	promptController string = "\\(Cisco Controller\\) >"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the transcript of the commands.
func (r *Repository) Fetch() (string, error) {
	var login, commands []x.Batcher
	var data string
	var err error

	// The controller prompts User:/Password: again after SSH authentication, so one batch serves both transports.
	login, commands = r.buildRequest()

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

func (r *Repository) buildRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: "User:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: promptController},
		&x.BSnd{S: "config paging disable\n"},
		&x.BExp{R: promptController},
	}, x.Commands(promptController, "\n", r.Config.Commands)
}
