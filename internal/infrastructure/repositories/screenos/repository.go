package repository

import (
	"time"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/ssh"
	"github.com/umatare5/telee/pkg/telnet"

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
	var data string
	var err error

	if r.Config.SecureMode && r.Config.RedundantMode {
		expects = r.buildSecureRequest(haSuffix)
	}
	if r.Config.SecureMode && !r.Config.RedundantMode {
		expects = r.buildSecureRequest(noSuffix)
	}
	if !r.Config.SecureMode && r.Config.RedundantMode {
		expects = r.buildRequest(haSuffix)
	}
	if !r.Config.SecureMode && !r.Config.RedundantMode {
		expects = r.buildRequest(noSuffix)
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

// [platform: ssg] buildSecureRequest returns the expects
func (r *Repository) buildSecureRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
	}
}
