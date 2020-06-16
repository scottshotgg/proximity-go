package grpc

import (
	"github.com/scottshotgg/proximity/pkg/buffs"
)

func (g *grpcSender) Send(msg []byte) error {
	return g.sendClient.Send(&buffs.SendReq{
		Msg: &buffs.Message{
			Route:    g.route,
			Contents: string(msg),
		},
	})
}
