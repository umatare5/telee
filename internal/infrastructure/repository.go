package infrastructure

import (
	"telee/internal/config"

	aireosRepository "telee/internal/infrastructure/repositories/aireos"
	alliedwareRepository "telee/internal/infrastructure/repositories/alliedware"
	asasoftwareRepository "telee/internal/infrastructure/repositories/asasoftware"
	iosRepository "telee/internal/infrastructure/repositories/ios"
	ironwareRepository "telee/internal/infrastructure/repositories/ironware"
	screenosRepository "telee/internal/infrastructure/repositories/screenos"
	yamahaosRepository "telee/internal/infrastructure/repositories/yamahaos"
)

// Repository struct
type Repository struct {
	Config *config.Config
}

// New returns Repository struct
func New(c *config.Config) Repository {
	return Repository{
		Config: c,
	}
}

// InvokeAireOSRepository returns AireOSRepository struct
func (r *Repository) InvokeAireOSRepository() *aireosRepository.Repository {
	return &aireosRepository.Repository{
		Config: r.Config,
	}
}

// InvokeAlliedWareRepository returns AlliedWareRepository struct
func (r *Repository) InvokeAlliedWareRepository() *alliedwareRepository.Repository {
	return &alliedwareRepository.Repository{
		Config: r.Config,
	}
}

// InvokeIOSRepository returns IOSRepository struct
func (r *Repository) InvokeIOSRepository() *iosRepository.Repository {
	return &iosRepository.Repository{
		Config: r.Config,
	}
}

// InvokeASASoftwareRepository returns ASASoftwareRepository struct
func (r *Repository) InvokeASASoftwareRepository() *asasoftwareRepository.Repository {
	return &asasoftwareRepository.Repository{
		Config: r.Config,
	}
}

// InvokeIronWareRepository returns IronWareRepository struct
func (r *Repository) InvokeIronWareRepository() *ironwareRepository.Repository {
	return &ironwareRepository.Repository{
		Config: r.Config,
	}
}

// InvokeScreenOSRepository returns ScreenOSRepository struct
func (r *Repository) InvokeScreenOSRepository() *screenosRepository.Repository {
	return &screenosRepository.Repository{
		Config: r.Config,
	}
}

// InvokeYamahaOSRepository returns YamahaOSRepository struct
func (r *Repository) InvokeYamahaOSRepository() *yamahaosRepository.Repository {
	return &yamahaosRepository.Repository{
		Config: r.Config,
	}
}
