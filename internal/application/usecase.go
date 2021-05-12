package application

import (
	aireosUsecase "telee/internal/application/usecases/aireos"
	alliedwareUsecase "telee/internal/application/usecases/alliedware"
	asasoftwareUsecase "telee/internal/application/usecases/asasoftware"
	iosUsecase "telee/internal/application/usecases/ios"
	ironwareUsecase "telee/internal/application/usecases/ironware"
	screenosUsecase "telee/internal/application/usecases/screenos"
	yamahaosUsecase "telee/internal/application/usecases/yamahaos"
	"telee/internal/config"
	"telee/internal/infrastructure"
)

// Usecase struct
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// New returns Usecase struct
func New(c *config.Config, r *infrastructure.Repository) Usecase {
	return Usecase{
		Config:     c,
		Repository: r,
	}
}

// InvokeAireOSUsecase returns AireOSUsecase struct
func (u *Usecase) InvokeAireOSUsecase() *aireosUsecase.Usecase {
	return &aireosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeAlliedWareUsecase returns AlliedWareUsecase struct
func (u *Usecase) InvokeAlliedWareUsecase() *alliedwareUsecase.Usecase {
	return &alliedwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeASASoftwareUsecase returns ASASoftwareUsecase struct
func (u *Usecase) InvokeASASoftwareUsecase() *asasoftwareUsecase.Usecase {
	return &asasoftwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeIOSUsecase returns IOSUsecase struct
func (u *Usecase) InvokeIOSUsecase() *iosUsecase.Usecase {
	return &iosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeIronWareUsecase returns IronWareUsecase
func (u *Usecase) InvokeIronWareUsecase() *ironwareUsecase.Usecase {
	return &ironwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeScreenOSUsecase returns ScreenOSUsecase struct
func (u *Usecase) InvokeScreenOSUsecase() *screenosUsecase.Usecase {
	return &screenosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeYamahaOSUsecase returns YamahaOSUsecase struct
func (u *Usecase) InvokeYamahaOSUsecase() *yamahaosUsecase.Usecase {
	return &yamahaosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}
