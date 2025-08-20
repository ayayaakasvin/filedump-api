package httpserver

import (
	"net/http"
	"sync"
	"time"
	"up-down-server/internal/config"
	"up-down-server/internal/http-server/handlers"
	"up-down-server/internal/http-server/middlewares"
	"up-down-server/internal/models"

	"github.com/ayayaakasvin/lightmux"
	"github.com/sirupsen/logrus"
)

type ServerApp struct {
	server *http.Server

	lmux *lightmux.LightMux

	cfg      *config.HTTPServer
	fileRepo models.FileMetaRepository
	userRepo models.UserRepository
	cache    models.Cache
	wg       *sync.WaitGroup

	logger *logrus.Logger
}

func NewServerApp(cfg *config.HTTPServer, file models.FileMetaRepository, user models.UserRepository, cache models.Cache, logger *logrus.Logger, wg *sync.WaitGroup) *ServerApp {
	return &ServerApp{
		cfg:      cfg,
		fileRepo: file,
		userRepo: user,
		cache:    cache,
		logger:   logger,
		wg:       wg,
	}
}

func (s *ServerApp) Run() {
	defer s.wg.Done()

	s.setupServer()

	s.setupLightMux()

	s.startServer()
}

func (s *ServerApp) startServer() {
	s.logger.Infof("Server has been started on port: %s", s.cfg.Address)
	s.logger.Infof("Available handlers:\n")

	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()

	go func() {
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			s.logger.Info("Server is running...")
		}
	}()

	// RunTLS can be run when server is hosted on domain, acts as seperate service of file storing, for my project, id chose to encapsulate servers under one docker-compose and make nginx-gateaway for my api like auth, file, user service
	// if err := s.lmux.RunTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile); err != nil {
	if err := s.lmux.Run(); err != nil {
		s.logger.Fatalf("Server exited with error: %v", err)
	}
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		// s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.cfg.Address
	s.server.IdleTimeout = s.cfg.IdleTimeout
	s.server.ReadTimeout = s.cfg.Timeout
	s.server.WriteTimeout = s.cfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middlewares.NewHTTPMiddlewares(s.logger, s.cache, s.cfg.Cors)
	handlers := handlers.NewHTTPHandlers(s.fileRepo, s.userRepo, s.cache, s.logger)

	// global middlewares usage | recovery from panic, logger for logging(logrus) and cors
	s.lmux.Use(mws.RecoverMiddleware, mws.LoggerMiddleware, mws.CorsMiddleware)

	s.lmux.NewRoute("/api/ping").Handle(http.MethodGet, handlers.PingHandler())

	// /api file GET | POST
	apiGroup := s.lmux.NewGroup("/api", mws.JWTAuthMiddleware)
	apiGroup.NewRoute("/upload", mws.RateLimitMiddleware).Handle(http.MethodPost, handlers.UploadFile())
	apiGroup.NewRoute("/download", mws.RateLimitMiddleware).Handle(http.MethodGet, handlers.DownloadFile())
	apiGroup.NewRoute("/download/shared/", mws.RateLimitMiddleware).Handle(http.MethodGet, handlers.DownloadFileViaSharedLink())
	apiGroup.NewRoute("/sharelink", mws.RateLimitMiddleware).Handle(http.MethodGet, handlers.CreateShareLink())

	// /api/files metadata CRUD
	filesRoute := apiGroup.NewRoute("/files")
	filesRoute.Handle(http.MethodGet, handlers.ListFile())
	filesRoute.Handle(http.MethodDelete, handlers.DeleteFile())
	apiGroup.NewRoute("/files/metadata").Handle(http.MethodGet, handlers.GetFileMetaData())
	apiGroup.NewRoute("/files/rename").Handle(http.MethodPatch, handlers.UpdateFileName())

	// /api auth 
	authGroup := s.lmux.NewGroup("/api")
	authGroup.NewRoute("/login").Handle(http.MethodPost, handlers.LogIn())
	authGroup.NewRoute("/register").Handle(http.MethodPost, handlers.Register())
	authGroup.NewRoute("/logout", mws.JWTAuthMiddleware).Handle(http.MethodDelete, handlers.LogOut())
	authGroup.NewRoute("/refresh").Handle(http.MethodPost, handlers.RefreshTheToken())

	// /swagger/ that idk why not working and / route for redirecting to 404 page, that has to be done
	// s.lmux.Mux().HandleFunc("/swagger/", handlers.SwaggerHandler())
	s.lmux.Mux().HandleFunc("/", handlers.NotFound404())

	s.logger.Info("LightMux has been set up")
}