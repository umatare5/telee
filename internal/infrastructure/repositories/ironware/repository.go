// Package repository implements Brocade IronWare-specific data access layer, reached by -x foundry.
package repository

import (
	"time"

	x "github.com/google/goexpect"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/telnet"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one telnet session and returns the output the last prompt match captured.
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher
	var data string
	var err error

	if r.Config.EnableMode {
		expects = r.buildPrivilegedRequest()
	} else {
		expects = r.buildUserModeRequest()
	}

	// Telnet only; checkArguments refuses --secure-mode for this platform.
	data, err = telnet.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(&expects)
	if err != nil {
		return "", err
	}
	return data, nil
}

// IronWare is the only platform needing CRLF. A bare "\n" leaves the line unsent.
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
