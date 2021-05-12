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

// [platform: aireos] buildRequest returns the expects
func (r *Repository) buildRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "User:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
		&x.BSnd{S: "config paging disable\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
	}
}
