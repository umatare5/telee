package framework

import (
	"fmt"
	"telee/internal/application"
	"telee/internal/config"
	"telee/internal/domain"
	"telee/internal/infrastructure"
)

// Exec struct
type Exec struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// New returns Exec struct
func New(c *config.Config, r *infrastructure.Repository, u *application.Usecase) Exec {
	return Exec{
		Config:     c,
		Repository: r,
		Usecase:    u,
	}
}

// Run displays the command result
func (e *Exec) Run() {
	var err error
	var data string

	if e.Config.ExecPlatform == domain.AireOSPlatformName {
		data, err = e.Usecase.InvokeAireOSUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.AlliedWarePlatformName {
		data, err = e.Usecase.InvokeAlliedWareUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.ASASoftwarePlatformName {
		data, err = e.Usecase.InvokeASASoftwareUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.IOSPlatformName {
		data, err = e.Usecase.InvokeIOSUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.IronWarePlatformName {
		data, err = e.Usecase.InvokeIronWareUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.JunOSPlatformName {
		data, err = e.Usecase.InvokeJunOSUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.ScreenOSPlatformName {
		data, err = e.Usecase.InvokeScreenOSUsecase().Fetch()
	}
	if e.Config.ExecPlatform == domain.YamahaOSPlatformName {
		data, err = e.Usecase.InvokeYamahaOSUsecase().Fetch()
	}

	if err != nil {
		fmt.Println(err)
		fmt.Println(domain.HintTelnetFailed)
		return
	}
	fmt.Println(data)
}
