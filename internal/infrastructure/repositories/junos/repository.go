package repository

import (
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/ssh"
	"time"

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
	data, err = ssh.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(&expects, ssh.GenerateClientConfig(r.Config.Username, r.Config.Password))

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
