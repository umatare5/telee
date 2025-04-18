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

	if r.Config.SecureMode && r.Config.DefaultPrivMode {
		expects = r.buildDefaultPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && r.Config.EnableMode {
		expects = r.buildPrivilegedSecureRequest()
	}
	if r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
		expects = r.buildUserModeSecureRequest()
	}
	if !r.Config.SecureMode && r.Config.DefaultPrivMode {
		expects = r.buildDefaultPrivilegedRequest()
	}
	if !r.Config.SecureMode && r.Config.EnableMode {
		expects = r.buildPrivilegedRequest()
	}
	if !r.Config.SecureMode && !r.Config.DefaultPrivMode && !r.Config.EnableMode {
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

// [platform: nxos] buildUserModeRequest returns the expects
func (r *Repository) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
	}
}

// [platform: nxos] buildPrivilegedRequest returns the expects
func (r *Repository) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

// [platform: nxos] buildDefaultPrivilegedRequest returns the expects
func (r *Repository) buildDefaultPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: r.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.Password + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

// [platform: nxos] buildUserModeSecureRequest returns the expects
func (r *Repository) buildUserModeSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "> "},
	}
}

// [platform: nxos] buildPrivilegedSecureRequest returns the expects
func (r *Repository) buildPrivilegedSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "> "},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: r.Config.PrivPassword + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}

// [platform: nxos] buildDefaultPrivilegedSecureRequest returns the expects
func (r *Repository) buildDefaultPrivilegedSecureRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
		&x.BSnd{S: r.Config.Command + "\n"},
		&x.BExp{R: r.Config.Hostname + "# "},
	}
}
