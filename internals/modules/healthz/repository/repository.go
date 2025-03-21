package repository

import (
	"context"
	"errors"

	"github.com/ocuris/go-template-backend/internals/modules/healthz"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewHealthzRepository(db *gorm.DB) healthz.Repository {
	return &repository{
		db: db,
	}
}

// PingDatabase checks if the database is reachable
func (r *repository) PingDatabase() (string, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return "unhealthy", errors.New("failed to get database instance")
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		return "unhealthy", errors.New("database is unreachable")
	}

	return "healthy", nil
}

// PingCache checks if the Redis cache is reachable (if applicable)
// func (r *repository) PingCache() (string, error) {
// 	if err := r.redisClient.Ping(context.Background()).Err(); err != nil {
// 		return "unhealthy", errors.New("cache is unreachable")
// 	}
// 	return "healthy", nil
// }
