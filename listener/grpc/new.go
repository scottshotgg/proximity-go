package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/scottshotgg/proximity/pkg/buffs"
	"github.com/scottshotgg/proximity/pkg/listener"
	"google.golang.org/grpc"
)

func New(id, route, addr string, handle handler) (listener.Listener, error) {
	if handle == nil {
		return nil, errors.New("handle cannot be nil")
	}

	var conn, err = grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	var client = buffs.NewNodeClient(conn)

	lis, err := client.Attach(context.Background(), &buffs.AttachReq{
		Id:    id,
		Route: route,
	})

	if err != nil {
		return nil, err
	}

	go func() {
		for {
			var res, err = lis.Recv()
			if err != nil {
				log.Fatalln("err listening:", err)
			}

			var msg = res.GetMessage()

			handle(&listener.Msg{
				Route:    msg.GetRoute(),
				Contents: msg.GetContents(),
			})
		}
	}()

	return &grpcListener{
		id:     id,
		route:  route,
		handle: handle,
		node:   client,
	}, nil
}
