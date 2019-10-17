package config

import (
	"services/apps/testing/models/settings"
	"services/package/shared_models"
)

// Service represent the caches's functions
type Service interface {
	Parse() error
	GetApplicationSettings() settings.Application
	GetCacheSettings() settings.Cache
	GetDatabaseSettings() shared_models.Database
	GetUserServiceSettings() shared_models.ServiceConnection
}
