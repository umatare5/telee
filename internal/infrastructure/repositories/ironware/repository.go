// Package repository implements Brocade IronWare-specific data access layer, reached by -x foundry.
package repository

import (
	"time"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	x "github.com/umatare5/telee/pkg/expect"
	"github.com/umatare5/telee/pkg/telnet"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one telnet session and returns the transcript of the commands.
func (r *Repository) Fetch() (string, error) {
	var login, commands []x.Batcher
	var data string
	var err error

	if r.Config.EnableMode {
		login, commands = r.buildPrivilegedRequest()
	} else {
		login, commands = r.buildUserModeRequest()
	}

	// Telnet only; checkArguments refuses --secure-mode for this platform.
	data, err = telnet.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(login, commands)
	if err != nil {
		return "", err
	}
	return data, nil
}

// IronWare is the only platform needing CRLF. A bare "\n" leaves the line unsent.
func (r *Repository) buildUserModeRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: "Please Enter Login Name:"},
		&x.BSnd{S: r.Config.Username + "\r\n"},
		&x.BExp{R: "Please Enter Password:"},
		&x.BSnd{S: r.Config.Password + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
		&x.BSnd{S: "skip-page-display\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
	}, x.Commands("telnet@"+r.Config.Hostname+">", "\r\n", r.Config.Commands)
}

func (r *Repository) buildPrivilegedRequest() (login, commands []x.Batcher) {
	return []x.Batcher{
		&x.BExp{R: "Please Enter Login Name:"},
		&x.BSnd{S: r.Config.Username + "\r\n"},
		&x.BExp{R: "Please Enter Password:"},
		&x.BSnd{S: r.Config.Password + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + ">"},
		&x.BSnd{S: "enable\r\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + "#"},
		&x.BSnd{S: "skip-page-display\r\n"},
		&x.BExp{R: "telnet@" + r.Config.Hostname + "#"},
	}, x.Commands("telnet@"+r.Config.Hostname+"#", "\r\n", r.Config.Commands)
}
