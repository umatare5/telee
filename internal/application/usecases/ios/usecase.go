package usecase

import (
	"telee/internal/config"
	"telee/internal/infrastructure"

	x "github.com/google/goexpect"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// Fetch returns stdout from telnet session
func (u *Usecase) Fetch() (string, error) {
	var expectation []x.Batcher

	if u.Config.EnableMode {
		expectation = u.buildPrivilegedRequest()
	} else {
		expectation = u.buildUserModeRequest()
	}

	data, err := u.Repository.InvokeServerRepository().Fetch(&expectation)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: ios] buildRequest returns the expectation
func (u *Usecase) buildUserModeRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: u.Config.Hostname + ">"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: u.Config.Hostname + ">"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: u.Config.Hostname + ">"},
	}
}

// [platform: ios] buildPrivilegedRequest returns the expectation
func (u *Usecase) buildPrivilegedRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "Username:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: u.Config.Hostname + ">"},
		&x.BSnd{S: "enable\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.PrivPassword + "\n"},
		&x.BExp{R: u.Config.Hostname + "#"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: u.Config.Hostname + "#"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: u.Config.Hostname + "#"},
	}
}
