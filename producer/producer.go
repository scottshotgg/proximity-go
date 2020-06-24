package producer

import "github.com/scottshotgg/proximity/pkg/listener"

type (
	Producer interface {
		Single(msg *listener.Msg) error
		Stream(ch <-chan *listener.Msg)
	}
)
