package events_logging

import (
	"context"
	"encoding/json"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/slntopp/nocloud/pkg/nocloud"
	"go.uber.org/zap"
)

type GelfServer struct {
	*gelf.Reader
	rep *SqliteRepository

	log *zap.Logger
}

type ShortLogMessage struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`

	Entity    string `json:"entity,omitempty"`
	Uuid      string `json:"uuid,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Action    string `json:"action,omitempty"`
	Rc        int32  `json:"rc,omitempty"`
	Requestor string `json:"requestor,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`

	Diff string `json:"diff,omitempty"`
}

func NewGelfServer(_log *zap.Logger, host string, rep *SqliteRepository) *GelfServer {
	log := _log.Named("GelfServer")

	log.Debug("Creating Gelf Server")

	reader, err := gelf.NewReader(host)
	if err != nil {
		log.Fatal("Failed to create GelfServer", zap.Error(err))
		return nil
	}
	return &GelfServer{Reader: reader, rep: rep, log: log}
}

func (s *GelfServer) Run() {
	log := s.log.Named("Run")

	log.Info("Start accepting messages")
	nocloudLevelVal := nocloud.NOCLOUD_LOG_LEVEL.String()

	for {
		message, err := s.ReadMessage()
		log.Info("Accept message", zap.String("Short", message.Short))
		if err != nil {
			log.Error("Failed to read message", zap.Error(err))
			continue
		}
		var shortMessage ShortLogMessage

		err = json.Unmarshal([]byte(message.Short), &shortMessage)
		if err != nil {
			log.Warn("Wrong message", zap.String("Short message", message.Short))
			log.Error("Failed to parse short message", zap.Error(err))
			continue
		}

		if shortMessage.Level != nocloudLevelVal {
			continue
		}

		err = s.rep.CreateEvent(context.Background(), &shortMessage)
		if err != nil {
			log.Error("Failed to create event", zap.Error(err))
		}
	}
}
