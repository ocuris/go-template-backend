package usecase

import (
	"errors"

	"github.com/ocuris/go-template-backend/internals/modules/healthz"
)

type usecase struct {
	repository healthz.Repository
}

func NewHealthzUsecase(repo healthz.Repository) healthz.Usecase {
	return &usecase{
		repository: repo,
	}
}

func (u *usecase) CheckHealth() (map[string]string, error) {
	// Check database health
	dbStatus, err := u.repository.PingDatabase()
	if err != nil {
		return map[string]string{"database": "unhealthy"}, errors.New("database is unreachable")
	}

	// Check cache health (if applicable)
	// cacheStatus, err := u.repository.PingCache()
	// if err != nil {
	// 	return map[string]string{"cache": "unhealthy"}, errors.New("cache is unreachable")
	// }

	// If everything is fine, return a healthy response
	return map[string]string{
		"database": dbStatus,
		// "cache":    cacheStatus,
		"status": "healthy",
	}, nil
}
