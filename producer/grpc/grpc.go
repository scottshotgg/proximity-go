package grpc

import (
	"context"
	"log"

	"github.com/scottshotgg/proximity-go/producer"
	"github.com/scottshotgg/proximity/pkg/buffs"
	"github.com/scottshotgg/proximity/pkg/listener"
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
	grpcProducer struct {
		stream buffs.Node_PublishClient
	}
)

func New(addr string) (producer.Producer, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), defaultOps)
	if err != nil {
		return nil, err
	}

	var client = buffs.NewNodeClient(conn)

	pub, err := client.Publish(context.Background())
	if err != nil {
		log.Fatalln("err", err)
	}

	return &grpcProducer{
		stream: pub,
	}, nil
}

func (g *grpcProducer) Single(msg *listener.Msg) error {
	return g.stream.Send(&buffs.PublishReq{
		Route:    msg.Route,
		Contents: msg.Contents,
	})
}

func (g *grpcProducer) Stream(ch <-chan *listener.Msg) {
	go func() {
		for {
			select {
			case msg := <-ch:
				// TODO: check error here later probably
				g.stream.Send(&buffs.PublishReq{
					Route:    msg.Route,
					Contents: msg.Contents,
				})
			}
		}
	}()
}
