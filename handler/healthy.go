package handler

import (
	"context"

	"github.com/micro/go-micro/v2/logger"

	proto "github.com/xtech-cloud/hkg-msp-builder/proto/builder"
)

type Healthy struct{}

func (this *Healthy) Echo(_ctx context.Context, _req *proto.EchoRequest, _rsp *proto.EchoResponse) error {
	logger.Infof("Received Healthy.Echo, msg is %v", _req.Msg)

	_rsp.Msg = _req.Msg

	return nil
}
