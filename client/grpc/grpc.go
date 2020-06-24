package grpc

import (
	"context"
	"errors"
	"log"

	"github.com/scottshotgg/proximity-go/consumer"
	"github.com/scottshotgg/proximity/pkg/buffs"
	"github.com/scottshotgg/proximity/pkg/listener"
	"google.golang.org/grpc"
)

const (
	maxMsgSize = 100 * 1024 * 1024
)

type (
	grpcClient struct {
		id  string
		sub buffs.Node_SubscribeClient
		pub buffs.Node_PublishClient
	}
)

var (
	defaultOps = grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsgSize),
		grpc.MaxCallSendMsgSize(maxMsgSize),
	)
)

func New(addr string, topics []string) (consumer.Consumer, error) {
	var (
		ctx       = context.Background()
		conn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure(), defaultOps)
	)

	if err != nil {
		return nil, err
	}

	var client = buffs.NewNodeClient(conn)

	// Clients don't listen to topics; maybe we can change this later
	sub, err := client.Subscribe(ctx, &buffs.SubscribeReq{
		Topics: topics,
	})
	if err != nil {
		return nil, err
	}

	md, err := sub.Header()
	if err != nil {
		return nil, err
	}

	var idArray, ok = md["id"]
	if !ok {
		return nil, errors.New("ID header not passed back")
	}

	switch {
	case len(idArray) == 0:
		return nil, errors.New("No ID recieved")

	case len(idArray) < 1:
		return nil, errors.New("More than one ID recieved")
	}

	// Register a publisher
	pub, err := client.Publish(ctx)
	if err != nil {
		return nil, err
	}

	return &grpcClient{
		id:  idArray[0],
		pub: pub,
		sub: sub,
	}, nil
}

func (g *grpcClient) ID() string {
	return g.id
}

func (g *grpcClient) Recv() ([]byte, error) {
	var res, err = g.sub.Recv()

	return res.GetMessage().GetContents(), err
}

// TODO: need some way to close this channel later on
func (g *grpcClient) Listen(ch chan<- []byte) {
	go func() {
		for {
			var res, err = g.sub.Recv()
			if err != nil {
				log.Fatalln("err g.stream.Recv:", err)
			}

			ch <- res.GetMessage().GetContents()
		}
	}()
}

func (g *grpcClient) Send(msg *listener.Msg) error {
	return g.pub.Send(&buffs.PublishReq{
		Routes:   []string{msg.Route},
		Contents: msg.Contents,
	})
}

// TODO: should probably return an error channel here?
func (g *grpcClient) Stream(ch <-chan *listener.Msg) {
	go func() {
		for {
			select {
			case msg := <-ch:
				// TODO: check error here later probably
				g.pub.Send(&buffs.PublishReq{
					Routes:   []string{msg.Route},
					Contents: msg.Contents,
				})
			}
		}
	}()
}
