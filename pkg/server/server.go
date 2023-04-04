package server

import (
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/support-pl/nocloud-gelf/pkg/repository"
	"go.uber.org/zap"
)

type GelfServer struct {
	*gelf.Reader
	rep *repository.SqliteRepository

	log *zap.Logger
}

func NewGelfServer(logger *zap.Logger, host string, rep *repository.SqliteRepository) *GelfServer {
	reader, err := gelf.NewReader(host)
	if err != nil {
		return nil
	}
	return &GelfServer{Reader: reader, rep: rep, log: logger}
}

func (s *GelfServer) Run() {
	for {
		_, err := s.ReadMessage()
		if err != nil {
			s.log.Error("Cannot read log message", zap.Error(err))
			return
		}
	}
}
