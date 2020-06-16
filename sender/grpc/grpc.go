package grpc

import (
	"github.com/scottshotgg/proximity/pkg/buffs"
)

type (
	grpcSender struct {
		id         string
		route      string
		node       buffs.NodeClient
		sendClient buffs.Node_SendClient
	}
)
