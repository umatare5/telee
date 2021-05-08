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
	var expection []x.Batcher

	if u.Config.EnableMode {
		expection = u.buildPrivilegedRequest()
	} else {
		expection = u.buildUserModeRequest()
	}

	data, err := u.Repository.InvokeServerRepository().Fetch(&expection)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: ios] buildRequest returns the expection
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
		&x.BSnd{S: "exit\n"},
	}
}

// [platform: ios] buildPrivilegedRequest returns the expection
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
