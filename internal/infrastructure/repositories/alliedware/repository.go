// Package repository implements AlliedTelesis AlliedWare-specific data access layer, reached by -x allied.
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

// Fetch runs one telnet session and returns the output the last prompt match captured.
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher
	var data string
	var err error

	expects = r.buildRequest()

	// Telnet only; checkArguments refuses --secure-mode for this platform.
	data, err = telnet.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	).Fetch(&expects)
	if err != nil {
		return "", err
	}
	return data, nil
}

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
