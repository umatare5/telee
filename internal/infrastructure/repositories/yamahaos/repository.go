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
	var expects []x.Batcher

	if r.Config.EnableMode {
		expects = r.buildPrivilegedRequest()
	} else {
		expects = r.buildUserModeRequest()
	}

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

// [platform: yamaha] buildRequest returns the expectation
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "console lines infinity\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
	}
}

// [platform: yamaha] buildPrivilegedRequest returns the expectation
func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "administrator\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: "console lines infinity\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
	}
}
