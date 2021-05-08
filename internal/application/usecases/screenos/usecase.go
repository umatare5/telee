package usecase

import (
	"telee/internal/config"
	"telee/internal/infrastructure"

	x "github.com/google/goexpect"
)

const (
	noSuffix string = ""
	haSuffix string = "\\(M\\)"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	HAMode     bool
}

// Fetch returns stdout from telnet session
func (u *Usecase) Fetch() (string, error) {
	var expection []x.Batcher

	if u.HAMode {
		expection = u.buildRequest(haSuffix)
	} else {
		expection = u.buildRequest(noSuffix)
	}

	data, err := u.Repository.InvokeServerRepository().Fetch(&expection)
	if err != nil {
		return "", err
	}
	return data, nil
}

// [platform: ssg] buildRequest returns the expection
func (u *Usecase) buildRequest(suffix string) []x.Batcher {
	return []x.Batcher{
		&x.BExp{R: "login:"},
		&x.BSnd{S: u.Config.Username + "\n"},
		&x.BExp{R: "password:"},
		&x.BSnd{S: u.Config.Password + "\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "->"},
		&x.BSnd{S: "set console page 0\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "->"},
		&x.BSnd{S: u.Config.Command + "\n"},
		&x.BExp{R: u.Config.Hostname + suffix + "->"},
	}
}
