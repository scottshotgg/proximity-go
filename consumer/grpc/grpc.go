package grpc

import (
	"context"
	"log"

	"github.com/scottshotgg/proximity-go/consumer"
	"github.com/scottshotgg/proximity/pkg/buffs"
	"google.golang.org/grpc"
)

const (
	maxMsgSize = 100 * 1024 * 1024
)

var defaultOps = grpc.WithDefaultCallOptions(
	grpc.MaxCallRecvMsgSize(maxMsgSize),
	grpc.MaxCallSendMsgSize(maxMsgSize),
)

type (
	grpcConsumer struct {
		id     string
		route  string
		stream buffs.Node_SubscribeClient
	}
)

func New(addr, id, route string) (consumer.Consumer, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), defaultOps)
	if err != nil {
		return nil, err
	}

	var client = buffs.NewNodeClient(conn)

	sub, err := client.Subscribe(context.Background(), &buffs.SubscribeReq{
		Id:    id,
		Route: route,
	})
	if err != nil {
		return nil, err
	}

	return &grpcConsumer{
		stream: sub,
	}, nil
}

func (g *grpcConsumer) Single() ([]byte, error) {
	var res, err = g.stream.Recv()

	return res.GetMessage().GetContents(), err
}

func (g *grpcConsumer) Stream(ch chan<- []byte) {
	go func() {
		for {
			var res, err = g.stream.Recv()
			if err != nil {
				log.Fatalln("err g.stream.Recv:", err)
			}

			ch <- res.GetMessage().GetContents()
		}
	}()
}

// func (g *grpcConsumer) Close()
