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
	expection := u.buildRequest()
	data, err := u.Repository.InvokeServerRepository().Fetch(&expection)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: aireos] buildRequest returns the expection
func (u *Usecase) buildRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "User:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "Password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
		&x.BSnd{S: "config paging disable\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: "\\(Cisco Controller\\) >"},
	}
}
