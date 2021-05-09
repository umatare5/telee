package usecase

import (
	"telee/internal/config"
	"telee/internal/infrastructure"

	x "github.com/google/goexpect"
)

const (
	noSuffix string = ""
	haSuffix string = "/pri/act"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	HAMode     bool
}

// Fetch returns stdout from telnet session
func (u *Usecase) Fetch() (string, error) {
	var expectation []x.Batcher

	if u.HAMode {
		expectation = u.buildPrivilegedRequest(haSuffix)
	}
	if !u.HAMode {
		expectation = u.buildPrivilegedRequest(noSuffix)
	}

	data, err := u.Repository.InvokeServerRepository().Fetch(&expectation)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: asa] buildPrivilegedRequest returns the expectation
func (u *Usecase) buildPrivilegedRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: u.Config.Hostname + suffix + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.PrivPassword + "\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "#"},
		&x.BSnd{S: "terminal pager 0\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "#"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "#"},
	}
}
