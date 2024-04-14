package server

import (
	"banners/config"
	auth "banners/internal/auth/infrastructure/delivery/http"
	authrepository "banners/internal/auth/infrastructure/repository"
	banners "banners/internal/banner/infrastructure/delivery/http"
	"banners/internal/banner/infrastructure/repository"
	"banners/pkg/cache"
	"banners/pkg/db"
	slogger "banners/pkg/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {

	slogger.Logger = slogger.GetLogger()

	dbConnection, err := db.NewDBConnection()

	if err != nil {
		panic("cannot connect to DB")
	}

	authRepo := authrepository.NewAuthRepository(dbConnection)
	authHandler := auth.AuthHandler{Store: authRepo}

	cache.Cache = *cache.LoadCache()

	bannerRepo := repository.NewBannerRepository(dbConnection, &cache.Cache)
	bannerHadler := banners.NewBannerHandler(bannerRepo)

	mux := http.NewServeMux()

	mux.Handle("/user/", &authHandler)
	mux.Handle("/", bannerHadler)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", config.Cfg.Server.Port), mux); err != nil {
			slogger.Logger.Error("error, server is crashed: ", "err", err)
		}
	}()
	slogger.Logger.Info("Listening to", "HOST", config.Cfg.Server.Host, "PORT", config.Cfg.Server.Port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-sigChan
	slogger.Logger.Info("Shutdown server", "signal", s)
}
