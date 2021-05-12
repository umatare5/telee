package repository

import (
	"telee/internal/config"
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
		return "", err
	}
	return data, nil
}

// [platform: ios] buildRequest returns the expects
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
	}
}

// [platform: ios] buildPrivilegedRequest returns the expects
func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "#"},
	}
}
