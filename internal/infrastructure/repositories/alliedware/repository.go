package repository

import (
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/telnet"
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

	expects = r.buildRequest()

	// Alliedware is not supporting SSH
	data, err = telnet.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(&expects)

	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: allied] buildRequest returns the expects
func (r *Repository) buildRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: "Manager " + r.Config.Hostname + ">"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: "Manager " + r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: "Manager " + r.Config.Hostname + ">"},
	}
}
