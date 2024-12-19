package events_logging

import (
	"context"
	"encoding/json"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/slntopp/nocloud/pkg/nocloud"
	"go.uber.org/zap"
	"runtime/debug"
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

	Priority int32 `json:"priority,omitempty"`
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

func (s *GelfServer) Run() int {
	log := s.log.Named("Run")

	defer func() {
		if err := recover(); err != nil {
			log.Error("Recovered from panic", zap.Any("error", err), zap.String("stack_trace", string(debug.Stack())))
		}
	}()

	log.Info("Start accepting messages")
	nocloudLevelVal := nocloud.NOCLOUD_LOG_LEVEL.String()

	for {
		message, err := s.ReadMessage()
		if err != nil {
			log.Error("Failed to read message", zap.Error(err))
			continue
		}
		var shortMessage ShortLogMessage

		err = json.Unmarshal([]byte(message.Short), &shortMessage)
		if err != nil {
			continue
		}

		if shortMessage.Level != nocloudLevelVal {
			continue
		}

		log.Info("Attempt to create event", zap.Any("Short message", shortMessage))
		err = s.rep.CreateEvent(context.Background(), &shortMessage)
		if err != nil {
			log.Error("Failed to create event", zap.Error(err))
		}
	}
}
