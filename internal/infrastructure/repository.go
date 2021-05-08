package infrastructure

import (
	"telee/internal/config"
	repository "telee/internal/infrastructure/repositories"
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

// InvokeServerRepository returns ServerRepository struct
func (r *Repository) InvokeServerRepository() *repository.ServerRepository {
	return &repository.ServerRepository{
		Config: r.Config,
	}
}
