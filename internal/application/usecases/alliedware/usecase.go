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
	expectation := u.buildRequest()
	data, err := u.Repository.InvokeServerRepository().Fetch(&expectation)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: allied] buildRequest returns the expectation
func (u *Usecase) buildRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: "Manager " + u.Config.Hostname + ">"},
		&x.BSnd{S: "terminal length 0\n"},
		&x.BExp{R: "Manager " + u.Config.Hostname + ">"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: "Manager " + u.Config.Hostname + ">"},
	}
}
