package usecase

import (
	"telee/internal/config"
	"telee/internal/infrastructure"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// Fetch returns stdout from ssh session
func (u *Usecase) Fetch() (string, error) {
	data, err := u.Repository.InvokeJunOSRepository().Fetch()
	if err != nil {
		return "", err
	}
	return data, nil
}