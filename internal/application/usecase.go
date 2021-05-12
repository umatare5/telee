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

// InvokeAireOSUsecase returns new AireOSUsecase
func (u *Usecase) InvokeAireOSUsecase() *aireosUsecase.Usecase {
	return &aireosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeAlliedWareUsecase returns new AlliedWareUsecase
func (u *Usecase) InvokeAlliedWareUsecase() *alliedwareUsecase.Usecase {
	return &alliedwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeASASoftwareUsecase returns new ASASoftwareUsecase
func (u *Usecase) InvokeASASoftwareUsecase() *asasoftwareUsecase.Usecase {
	return &asasoftwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeIOSUsecase returns new IOSUsecase
func (u *Usecase) InvokeIOSUsecase() *iosUsecase.Usecase {
	return &iosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeIronWareUsecase returns new IronWareUsecase
func (u *Usecase) InvokeIronWareUsecase() *ironwareUsecase.Usecase {
	return &ironwareUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeScreenOSUsecase returns new ScreenOSUsecase
func (u *Usecase) InvokeScreenOSUsecase() *screenosUsecase.Usecase {
	return &screenosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeYamahaOSUsecase returns new YamahaOSUsecase
func (u *Usecase) InvokeYamahaOSUsecase() *yamahaosUsecase.Usecase {
	return &yamahaosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}
