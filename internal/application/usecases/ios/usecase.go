package usecase

import (
	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/infrastructure"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// Fetch returns stdout from the session
func (u *Usecase) Fetch() (string, error) {
	data, err := u.Repository.InvokeIOSRepository().Fetch()
	if err != nil {
		return "", err
	}
	return data, nil
}
