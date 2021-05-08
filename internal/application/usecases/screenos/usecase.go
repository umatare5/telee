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

// [platform: ssg] buildRequest returns the expection
func (u *Usecase) buildRequest() []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: u.Config.Hostname + "\\(M\\)->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: u.Config.Hostname + "\\(M\\)->"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: u.Config.Hostname + "\\(M\\)->"},
	}
}
