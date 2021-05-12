package repository

import (
	"telee/internal/config"
	"telee/internal/domain"
	"telee/pkg/telnet"
	"time"

	x "github.com/google/goexpect"
)

const (
	noSuffix string = ""
	haSuffix string = "\\(M\\)"
)

// Repository struct
type Repository struct {
	Config *config.Config
}

// Fetch returns stdout from telnet session
func (r *Repository) Fetch() (string, error) {
	var expects []x.Batcher

	if r.Config.HAMode {
		expects = r.buildRequest(haSuffix)
	} else {
		expects = r.buildRequest(noSuffix)
	}

	client := telnet.New(
		r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
	)
	data, err := client.Fetch(&expects)
	if err != nil {
		return "", err
	}

	return data, nil
}

// [platform: ssg] buildRequest returns the expects
func (r *Repository) buildRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
	}
}
