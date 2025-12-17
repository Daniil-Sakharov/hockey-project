package di

import (
	"github.com/Daniil-Sakharov/HockeyProject/internal/config"
)

// BaseContainer содержит общие зависимости для всех приложений
type BaseContainer struct {
	infra   *Infrastructure
	repo    *Repository
	service *Service
}

func NewBaseContainer(cfg *config.Config) *BaseContainer {
	infra := NewInfrastructure(cfg)
	repo := NewRepository(infra)
	service := NewService(cfg, infra, repo)

	return &BaseContainer{
		infra:   infra,
		repo:    repo,
		service: service,
	}
}

func (bc *BaseContainer) Infrastructure() *Infrastructure {
	return bc.infra
}

func (bc *BaseContainer) Repository() *Repository {
	return bc.repo
}

func (bc *BaseContainer) Service() *Service {
	return bc.service
}
