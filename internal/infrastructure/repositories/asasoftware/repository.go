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

const (
	noSuffix string = ""
	haSuffix string = "/pri/act"
)

// Repository struct
type Repository struct {
	Config *config.Config
}

// Fetch returns stdout from telnet session
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher

	if r.Config.HAMode {
		expects = r.buildPrivilegedRequest(haSuffix)
	}
	if !r.Config.HAMode {
		expects = r.buildPrivilegedRequest(noSuffix)
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

// [platform: asa] buildPrivilegedRequest returns the expects
func (r *Repository) buildPrivilegedRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: "terminal pager 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "#"},
	}
}
