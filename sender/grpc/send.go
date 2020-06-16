package grpc

func (g *grpcSender) Send(msg []byte) error {
	g.sendClient.Send(&buffs.SendReq{
		Msg: &buffs.Message{
			Route:    g.route,
			Contents: string(msg),
		},
	})
}