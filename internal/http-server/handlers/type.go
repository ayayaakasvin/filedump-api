// Handlers that serves for main http server, accessed via handlerd.Handler struct that contains necessary dependencies
package handlers

import (
	"up-down-server/internal/models"

	"github.com/sirupsen/logrus"
)

type Handlers struct {
	fileRepo models.FileMetaRepository
	userRepo models.UserRepository
	cache    models.Cache

	logger *logrus.Logger
}

func NewHTTPHandlers(file models.FileMetaRepository, user models.UserRepository, cache models.Cache, logger *logrus.Logger) *Handlers {
	return &Handlers{
		fileRepo: file,
		userRepo: user,
		cache:    cache,

		logger: logger,
	}
}