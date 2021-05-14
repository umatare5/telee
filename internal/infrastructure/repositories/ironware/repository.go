package repository

import (
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/ssh"
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

	if r.Config.EnableMode {
		expects = r.buildPrivilegedRequest()
	} else {
		expects = r.buildUserModeRequest()
	}

	if r.Config.SecureMode {
		data, err = ssh.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects, ssh.GenerateClientConfig(r.Config.Username, r.Config.Password))
	} else {
		data, err = telnet.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects)
	}

	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: foundry] buildUserModeRequest returns the expects
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Please Enter Login Name:"},
		&x.BSnd{S: r.Config.Username + "\r\n"},
		&x.BExp{R: "Please Enter Password:"},
		&x.BSnd{S: r.Config.Password + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
		&x.BSnd{S: "skip-page-display\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
		&x.BSnd{S: r.Config.Command + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
	}
}

// [platform: foundry] buildPrivilegedRequest returns the expects
func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Please Enter Login Name:"},
		&x.BSnd{S: r.Config.Username + "\r\n"},
		&x.BExp{R: "Please Enter Password:"},
		&x.BSnd{S: r.Config.Password + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
		&x.BSnd{S: "enable\r\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\r\n"},
		&x.BSnd{S: "skip-page-display\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + "#"},
		&x.BSnd{S: r.Config.Command + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + "#"},
	}
}
