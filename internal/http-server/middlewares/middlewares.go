package middlewares

import (
	"strings"
	"up-down-server/internal/config"
	"up-down-server/internal/models"

	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache  models.Cache
	logger *logrus.Logger

	allowed_origins []string
	allowed_headers string
	allowed_methods string
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache models.Cache, cfg config.CORSConfig) *Middlewares {
	return &Middlewares{
		logger: logger,
		cache:  cache,

		allowed_origins: cfg.AllowedOrigins,
		allowed_headers: strings.Join(cfg.AllowedHeaders, ", "),
		allowed_methods: strings.Join(cfg.AllowedMethods, ", "),
	}
}
