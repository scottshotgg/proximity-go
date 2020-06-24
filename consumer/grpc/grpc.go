package grpc

import (
	"context"
	"errors"
	"fmt"
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

func New(addr string, topics []string) (consumer.Consumer, error) {
	var conn, err = grpc.Dial(addr, grpc.WithInsecure(), defaultOps)
	if err != nil {
		return nil, err
	}

	sub, err := buffs.NewNodeClient(conn).Subscribe(context.Background(), &buffs.SubscribeReq{
		Topics: topics,
	})

	if err != nil {
		return nil, err
	}

	md, err := sub.Header()
	if err != nil {
		return nil, err
	}

	fmt.Println("header:", md)

	var idArray = md.Get("id")
	if len(idArray) != 1 {
		return nil, errors.New("something was wrong with the ID header")
	}

	return &grpcConsumer{
		id:     idArray[0],
		stream: sub,
	}, nil
}

func (g *grpcConsumer) ID() string {
	return g.id
}

func (g *grpcConsumer) Recv() ([]byte, error) {
	var res, err = g.stream.Recv()

	return res.GetMessage().GetContents(), err
}

func (g *grpcConsumer) Listen(ch chan<- []byte) {
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
