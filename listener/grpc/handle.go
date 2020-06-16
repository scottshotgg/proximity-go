package grpc

import (
	"github.com/scottshotgg/proximity/pkg/listener"
)

func (g *grpcListener) Handle(msg *listener.Msg) error {
	return g.handle(msg)
}
