// Package repository implements Juniper JunOS-specific data access layer, reached by -x srx.
package repository

import (
	"time"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	x "github.com/umatare5/telee/pkg/expect"
	"github.com/umatare5/telee/pkg/ssh"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one SSH session and returns the transcript of the commands.
func (r *Repository) Fetch() (string, error) {
	var login, commands []x.Batcher
	var data string
	var err error

	login, commands = r.buildUserModeSecureRequest()

	// SSH only; checkArguments refuses a run without --secure-mode.
	clientConfig, err := ssh.GenerateClientConfig(r.Config.Username, r.Config.Password, r.Config.HostKeyPath, r.Config.Hostname)
	if err != nil {
		return "", err
	}
	data, err = ssh.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(login, commands, clientConfig)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (r *Repository) buildUserModeSecureRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: r.Config.Username + "@" + r.Config.Hostname + ">"},
	}, x.Commands(r.Config.Username+"@"+r.Config.Hostname+">", " | no-more\n", r.Config.Commands)
}
