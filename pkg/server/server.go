package server

import (
	"encoding/json"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/slntopp/nocloud/pkg/nocloud"
	models "github.com/support-pl/nocloud-gelf/pkg/gelf_models"
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
	logLevelString := nocloud.NOCLOUD_LOG_LEVEL.String()

	for {
		logMessage, err := s.ReadMessage()
		if err != nil {
			s.log.Error("Cannot read log message", zap.Error(err))
			return
		}

		var shortMessage models.ShortMessage

		err = json.Unmarshal([]byte(logMessage.Short), &shortMessage)
		if err != nil {
			s.log.Error("Cannot unmarshal short message", zap.Error(err))
			return
		}

		if shortMessage.Level != logLevelString {
			continue
		}

	}
}
