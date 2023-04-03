package server

import (
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/support-pl/nocloud-gelf/pkg/repository"
)

type GelfServer struct {
	*gelf.Reader
	rep *repository.SqliteRepository
}

func NewGelfServer(port string, rep *repository.SqliteRepository) *GelfServer {
	reader, err := gelf.NewReader(port)
	if err != nil {
		return nil
	}
	return &GelfServer{Reader: reader, rep: rep}
}

func (s *GelfServer) Run() {

}
