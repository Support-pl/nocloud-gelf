package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"

	events "github.com/Support-pl/nocloud-gelf/pkg"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/slntopp/nocloud-proto/events_logging"
	"github.com/slntopp/nocloud/pkg/nocloud"
	"github.com/slntopp/nocloud/pkg/nocloud/auth"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

var (
	port, web_port string
	log            *zap.Logger

	sqliteHost  string
	gelfHost    string
	SIGNING_KEY []byte
	redisHost   string
)

func init() {
	viper.AutomaticEnv()
	log = nocloud.NewLogger()

	viper.SetDefault("PORT", "8000")
	viper.SetDefault("WEB_PORT", "8080")

	viper.SetDefault("GELF_HOST", ":12201")
	viper.SetDefault("SQLITE_HOST", "sqlite.db")
	viper.SetDefault("DB_HOST", "db:8529")
	viper.SetDefault("DB_CRED", "root:openSesame")
	viper.SetDefault("SIGNING_KEY", "seeeecreet")

	port = viper.GetString("PORT")
	web_port = viper.GetString("WEB_PORT")

	sqliteHost = viper.GetString("SQLITE_HOST")
	gelfHost = viper.GetString("GELF_HOST")
	SIGNING_KEY = []byte(viper.GetString("SIGNING_KEY"))

	viper.SetDefault("REDIS_HOST", "redis:6379")
	redisHost = viper.GetString("REDIS_HOST")
}

func main() {
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Setting up Sqlite Connection")
	repository := events.NewSqliteRepository(log, sqliteHost)
	log.Info("Sqlite connection established")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Failed to listen", zap.String("address", port), zap.Error(err))
	}

	log.Info("Connecting redis", zap.String("url", redisHost))
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0, // use default DB
	})

	auth.SetContext(log, rdb, SIGNING_KEY)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(log),
			grpc.UnaryServerInterceptor(auth.JWT_AUTH_INTERCEPTOR),
		)),
	)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gelfServer := events.NewGelfServer(log, gelfHost, repository)

	go gelfServer.Run()

	server := events.NewEventsLoggingServer(log, repository)
	pb.RegisterEventsLoggingServiceServer(s, server)
	err = pb.RegisterEventsLoggingServiceHandlerFromEndpoint(context.Background(), mux, "localhost:"+port, opts)
	if err != nil {
		log.Fatal("Cannot register AnsibleService Gateway handler", zap.Error(err))
	}

	go func() {
		log.Info(fmt.Sprintf("Serving HTTP on 0.0.0.0:%v", web_port), zap.Skip())
		http.ListenAndServe(fmt.Sprintf(":%s", web_port), wsproxy.WebsocketProxy(mux))
	}()

	log.Info(fmt.Sprintf("Serving gRPC on 0.0.0.0:%v", port), zap.Skip())
	log.Fatal("Failed to serve gRPC", zap.Error(s.Serve(lis)))
}
