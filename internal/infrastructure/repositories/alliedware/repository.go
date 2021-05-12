package repository

import (
	"fmt"
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/telnet"
	"time"

	x "github.com/google/goexpect"
)

const (
	protocol string = "tcp"
)

// Repository struct
type Repository struct {
	Config *config.Config
}

// Fetch returns stdout from telnet session
func (r *Repository) Fetch() (string, error) {
	expects := r.buildRequest()

	client := telnet.New(
		r.Config.Hostname, r.Config.Port, protocol, time.Duration(r.Config.Timeout)*time.Second,
	)
	data, err := client.Fetch(&expects)
	if err != nil {
		fmt.Println(domain.HintTelnetFailed)
		return "", err
	}
	return data, nil
}

// [platform: allied] buildRequest returns the expectation
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
