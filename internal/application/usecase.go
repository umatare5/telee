package application

import (
	aireosUsecase "telee/internal/application/usecases/aireos"
	alliedwareUsecase "telee/internal/application/usecases/alliedware"
	iosUsecase "telee/internal/application/usecases/ios"
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

// InvokeIOSUsecase returns IOSUsecase struct
func (u *Usecase) InvokeIOSUsecase() *iosUsecase.Usecase {
	return &iosUsecase.Usecase{
		Config:     u.Config,
		Repository: u.Repository,
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
