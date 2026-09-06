// Package infrastructure provides data access layer for network device communication.
package infrastructure

import (
	"github.com/umatare5/telee/internal/config"
	aireosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/aireos"
	alliedwareRepository "github.com/umatare5/telee/internal/infrastructure/repositories/alliedware"
	asasoftwareRepository "github.com/umatare5/telee/internal/infrastructure/repositories/asasoftware"
	iosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/ios"
	ironwareRepository "github.com/umatare5/telee/internal/infrastructure/repositories/ironware"
	junosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/junos"
	nxosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/nxos"
	screenosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/screenos"
	yamahaosRepository "github.com/umatare5/telee/internal/infrastructure/repositories/yamahaos"
)

type Repository struct {
	Config *config.Config
}

func New(c *config.Config) Repository {
	return Repository{
		Config: c,
	}
}

func (r *Repository) InvokeAireOSRepository() *aireosRepository.Repository {
	return &aireosRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeAlliedWareRepository() *alliedwareRepository.Repository {
	return &alliedwareRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeASASoftwareRepository() *asasoftwareRepository.Repository {
	return &asasoftwareRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeIOSRepository() *iosRepository.Repository {
	return &iosRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeIronWareRepository() *ironwareRepository.Repository {
	return &ironwareRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeJunOSRepository() *junosRepository.Repository {
	return &junosRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeNXOSRepository() *nxosRepository.Repository {
	return &nxosRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeScreenOSRepository() *screenosRepository.Repository {
	return &screenosRepository.Repository{
		Config: r.Config,
	}
}

func (r *Repository) InvokeYamahaOSRepository() *yamahaosRepository.Repository {
	return &yamahaosRepository.Repository{
		Config: r.Config,
	}
}
