package grpc

import (
	"github.com/scottshotgg/proximity/pkg/buffs"
	"github.com/scottshotgg/proximity/pkg/listener"
)

type (
	handler = func(msg *listener.Msg) error

	grpcListener struct {
		id     string
		route  string
		handle handler
		node   buffs.NodeClient
	}
)
