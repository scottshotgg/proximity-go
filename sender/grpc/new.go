package grpc

import (
	"context"

	"github.com/scottshotgg/proximity/pkg/buffs"
	"github.com/scottshotgg/proximity/pkg/sender"
	"google.golang.org/grpc"
)

func New(id, route, addr string) (sender.Sender, error) {
	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	var client = buffs.NewNodeClient(conn)

	sendClient, err := client.Send(context.Background())
	if err != nil {
		return nil, err
	}

	return &grpcSender{
		id:         id,
		route:      route,
		node:       client,
		sendClient: sendClient,
	}, nil
}
