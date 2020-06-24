package consumer

import (
	"github.com/scottshotgg/proximity/pkg/listener"
)

type (
	Client interface {
		ID() string

		Recv() ([]byte, error)
		Listen(ch chan<- []byte)

		Send(msg *listener.Msg) error
		Stream(ch <-chan *listener.Msg)
	}
)
