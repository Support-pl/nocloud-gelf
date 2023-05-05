package events_logging

import (
	"context"
	pb "github.com/slntopp/nocloud-proto/events_logging"
	"github.com/slntopp/nocloud/pkg/nocloud"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type EventsLoggingServer struct {
	pb.UnimplementedEventsLoggingServiceServer
	rep *SqliteRepository

	log *zap.Logger
}

func NewEventsLoggingServer(_log *zap.Logger, rep *SqliteRepository) *EventsLoggingServer {
	log := _log.Named("EventsLoggingServer")
	log.Debug("New EventsLogging Server Creating")

	return &EventsLoggingServer{log: log, rep: rep}
}

func (s *EventsLoggingServer) GetEvents(ctx context.Context, req *pb.GetEventsRequest) (*pb.Events, error) {
	log := s.log.Named("GetEvents")

	requestor := ctx.Value(nocloud.NoCloudAccount).(string)
	log.Debug("Request received", zap.Any("request", req), zap.String("requestor", requestor))

	events, err := s.rep.GetEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.Events{Events: events}, nil
}

func (s *EventsLoggingServer) GetCount(ctx context.Context, req *pb.GetEventsCountRequest) (*pb.GetEventsCountResponse, error) {
	log := s.log.Named("GetEvents")

	requestor := ctx.Value(nocloud.NoCloudAccount).(string)
	log.Debug("Request received", zap.Any("request", req), zap.String("requestor", requestor))

	total, err := s.rep.GetEventsCount(ctx, req)
	if err != nil {
		return nil, err
	}

	unique, err := s.rep.GetUnique(ctx)
	if err != nil {
		return nil, err
	}

	value, err := structpb.NewValue(unique)
	if err != nil {
		return nil, err
	}

	return &pb.GetEventsCountResponse{
		Total:  total,
		Unique: value,
	}, nil
}
