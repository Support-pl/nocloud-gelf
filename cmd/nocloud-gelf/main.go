package main

import (
	"github.com/spf13/viper"
	"github.com/support-pl/nocloud-gelf/pkg/repository"
	"github.com/support-pl/nocloud-gelf/pkg/server"
	"go.uber.org/zap"
)

var (
	log *zap.Logger

	db          string
	host        string
	SIGNING_KEY []byte
)

func init() {
	viper.AutomaticEnv()

	viper.SetDefault("DB", "store.db")
	viper.SetDefault("HOST", ":12201")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	db = viper.GetString("DB")
	host = viper.GetString("host")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	rep := repository.NewSqliteRepository(db)
	if rep == nil {
		log.Fatal("Failed to access sqlite db")
	}

	gelfServer := server.NewGelfServer(host, rep)
	if gelfServer == nil {
		log.Fatal("Failed to start server")
	}

	gelfServer.Run()
}
