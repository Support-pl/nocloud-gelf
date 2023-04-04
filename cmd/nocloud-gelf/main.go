package main

import (
	"github.com/slntopp/nocloud/pkg/nocloud"
	"github.com/spf13/viper"
	"github.com/support-pl/nocloud-gelf/pkg/repository"
	"github.com/support-pl/nocloud-gelf/pkg/server"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	db          string
	host        string
	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()

	logger = nocloud.NewLogger()

	viper.SetDefault("DB", "store.db")
	viper.SetDefault("HOST", ":12201")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	db = viper.GetString("DB")
	host = viper.GetString("host")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
}

func main() {
	defer func() {
		_ = logger.Sync()
	}()

	rep := repository.NewSqliteRepository(db)
	if rep == nil {
		logger.Fatal("Failed to access sqlite db")
	}

	gelfServer := server.NewGelfServer(logger, host, rep)
	if gelfServer == nil {
		logger.Fatal("Failed to start server")
	}

	gelfServer.Run()
}
