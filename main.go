package main

import (
	"os"
	"path"
	"sync"

	"up-down-server/internal/config"
	httpserver "up-down-server/internal/http-server"
	"up-down-server/internal/logger"
	"up-down-server/internal/models"
	"up-down-server/internal/repository/cache"
	"up-down-server/internal/repository/postgresql"

	_ "up-down-server/docs"

	"github.com/sirupsen/logrus"
)

func main() {
    cfg := config.MustLoadConfig()
	formattedLogger := logger.SetupLogger()

	shutdownChan := models.NewShutdownChannel()
	go func() {
		formattedLogger.Errorf("Error during setup: %s", shutdownChan.Value())
		os.Exit(1)
	}()

	repo := postgresql.NewPostgreSQLConnection(cfg.Database, shutdownChan)
	cache := cache.NewRedisClient(cfg.Redis, shutdownChan)
	mkdirFiles(formattedLogger, shutdownChan)

	wg := new(sync.WaitGroup)	
	wg.Add(1)
	
	app := httpserver.NewServerApp(&cfg.HTTPServer, repo, repo, cache, formattedLogger, wg)

	go app.Run()

	wg.Wait()
}

func mkdirFiles(logger *logrus.Logger, shutdown models.ShutdownChannel) {
	mode := 6400
	pathToFiles := path.Join(".", "files")
	if _, err:= os.Stat(pathToFiles); os.IsNotExist(err) {
		err := os.Mkdir(pathToFiles, os.FileMode(mode))
		if err != nil {
			shutdown.Send(models.ShutdownMessage, "mkdirFiles", "failed to make direction `files`")
			return
		}
	}

	logger.Info("Files direction initialized!")
}