// Package usecase implements Juniper ScreenOS-specific business logic.
package usecase

import (
	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/infrastructure"
)

type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

func (u *Usecase) Fetch() (string, error) {
	data, err := u.Repository.InvokeScreenOSRepository().Fetch()
	if err != nil {
		return "", err
	}
	return data, nil
}
