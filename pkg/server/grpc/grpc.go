package grpc

import (
	"context"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/bkit/util/inet"
	"github.com/bbdshow/qelog/api/receiverpb"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/receiver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
)

var (
	receiverSvc *receiver.Service
)

type ReceiverGrpc struct {
	*runner.GrpcServer
}

func NewReceiverGRpc(cfg *conf.Config, svc *receiver.Service) runner.Server {
	receiverSvc = svc
	rpc := &ReceiverGrpc{
		GrpcServer: runner.NewGrpcServer(),
	}
	rpc.RunAfter(func(s *grpc.Server) error {
		receiverpb.RegisterReceiverServer(s, rpc)
		return nil
	})
	return rpc
}

func (rpc *ReceiverGrpc) PushPacket(ctx context.Context, in *receiverpb.Packet) (*receiverpb.BaseResp, error) {
	// 获取 clientIP
	if err := receiverSvc.PacketToLogging(ctx, clientIP(ctx), in); err != nil {
		e, ok := err.(errc.Error)
		if ok {
			// 数据库操作错误
			if e.Code == errc.InternalErr {
				logs.Qezap.Error("PushPacket", zap.Error(e))
				return nil, errc.ErrInternalErr
			}
			return &receiverpb.BaseResp{
				Code:    int32(e.Code),
				Message: e.Message,
			}, nil
		}
		return nil, err
	}
	return &receiverpb.BaseResp{
		Code:    errc.Success,
		Message: "success",
	}, nil
}

func clientIP(ctx context.Context) string {
	ctxPeer, ok := peer.FromContext(ctx)
	if ok && ctxPeer.Addr != nil {
		if ipNet, ok := ctxPeer.Addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil || ipNet.IP.To16() != nil {
				return ipNet.IP.String()
			}
		}
		// 上述解析不成功，则自行拼接
		return inet.AddrStringToIP(ctxPeer.Addr)
	}
	return ""
}
