package repository

import (
	"time"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/ssh"
	"github.com/umatare5/telee/pkg/telnet"

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

	// AireOS needs the interactive authentication when also use SSH
	expects = r.buildRequest()

	if r.Config.SecureMode {
		clientConfig, err := ssh.GenerateClientConfig(r.Config.Username, r.Config.Password, r.Config.HostKeyPath)
		if err != nil {
			return "", err
		}
		data, err = ssh.New(
			r.Config.Hostname, r.Config.Port, domain.ProtocolTCP, time.Duration(r.Config.Timeout)*time.Second,
		).Fetch(&expects, clientConfig)
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
