package repository

import (
	"time"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/ssh"

	x "github.com/google/goexpect"
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

	expects = r.buildUserModeSecureRequest()

	// JunOS is not supporting Telnet
	clientConfig, err := ssh.GenerateClientConfig(r.Config.Username, r.Config.Password, r.Config.HostKeyPath)
	if err != nil {
		return "", err
	}
	data, err = ssh.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(&expects, clientConfig)

	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: junos] buildUserModeSecureRequest returns the expects
func (r *Repository) buildUserModeSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Username + "@" + r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + " | no-more\n"},
		&x.BExp{R: r.Config.Username + "@" + r.Config.Hostname + ">"},
	}
}
