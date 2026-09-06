// Package application coordinates business logic across platform-specific use cases.
package application

import (
	aireosUsecase "github.com/umatare5/telee/internal/application/usecases/aireos"
	alliedwareUsecase "github.com/umatare5/telee/internal/application/usecases/alliedware"
	asasoftwareUsecase "github.com/umatare5/telee/internal/application/usecases/asasoftware"
	iosUsecase "github.com/umatare5/telee/internal/application/usecases/ios"
	ironwareUsecase "github.com/umatare5/telee/internal/application/usecases/ironware"
	junosUsecase "github.com/umatare5/telee/internal/application/usecases/junos"
	nxosUsecase "github.com/umatare5/telee/internal/application/usecases/nxos"
	screenosUsecase "github.com/umatare5/telee/internal/application/usecases/screenos"
	yamahaosUsecase "github.com/umatare5/telee/internal/application/usecases/yamahaos"
	"github.com/umatare5/telee/internal/config"
	"github.com/umatare5/telee/internal/infrastructure"
)

type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

func New(c *config.Config, r *infrastructure.Repository) Usecase {
	return Usecase{
		Config:     c,
		Repository: r,
	}
}

func (u *Usecase) InvokeAireOSUsecase() *aireosUsecase.Usecase {
	return &aireosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeAlliedWareUsecase() *alliedwareUsecase.Usecase {
	return &alliedwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeASASoftwareUsecase() *asasoftwareUsecase.Usecase {
	return &asasoftwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeIOSUsecase() *iosUsecase.Usecase {
	return &iosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeIronWareUsecase() *ironwareUsecase.Usecase {
	return &ironwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeJunOSUsecase() *junosUsecase.Usecase {
	return &junosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeNXOSUsecase() *nxosUsecase.Usecase {
	return &nxosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeScreenOSUsecase() *screenosUsecase.Usecase {
	return &screenosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

func (u *Usecase) InvokeYamahaOSUsecase() *yamahaosUsecase.Usecase {
	return &yamahaosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}
