package producer

import "github.com/scottshotgg/proximity/pkg/listener"

type (
	Producer interface {
		Send(msg *listener.Msg) error
		Stream(ch <-chan *listener.Msg)
	}
)
