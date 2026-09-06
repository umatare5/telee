// Package repository implements Juniper ScreenOS-specific data access layer, reached by -x ssg.
package repository

import (
	"time"

	x "github.com/google/goexpect"
	cryptossh "golang.org/x/crypto/ssh"

	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/domain"
	"github.com/umatare5/telee/pkg/ssh"
	"github.com/umatare5/telee/pkg/telnet"
)

const (
	noSuffix string = ""
	// BExp is a regular expression, so the master-unit suffix keeps its parentheses escaped.
	haSuffix string = "\\(M\\)"
)

type Repository struct {
	Config *config.Config
}

// Fetch runs one session and returns the output the last prompt match captured.
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
		var clientConfig *cryptossh.ClientConfig
		clientConfig, err = ssh.GenerateClientConfig(r.Config.Username, r.Config.Password, r.Config.HostKeyPath, r.Config.Hostname)
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

func (r *Repository) buildRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		// ScreenOS is the only platform prompting lower-case "password:".
		&x.BExp{R: "password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
	}
}

func (r *Repository) buildSecureRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + suffix + "->"},
	}
}
