package receiver

import (
	"context"
	"net"

	"github.com/huzhongqing/qelog/pkg/common/kit"
	"google.golang.org/grpc/peer"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/api/receiverpb"
	"github.com/huzhongqing/qelog/pkg/storage"
	"google.golang.org/grpc"
)

type GRPCService struct {
	server   *grpc.Server
	receiver *Service
}

func NewGRPCService(sharding *storage.Sharding) *GRPCService {
	srv := &GRPCService{
		server:   nil,
		receiver: NewService(sharding),
	}

	return srv
}

func (srv *GRPCService) Run(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	srv.server = server

	receiverpb.RegisterReceiverServer(srv.server, srv)

	//reflection.Register(srv.server)

	if err := server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (srv *GRPCService) Close() error {
	srv.receiver.Sync()
	if srv.server != nil {
		srv.server.Stop()
	}
	return nil
}

func (srv *GRPCService) PushPacket(ctx context.Context, in *receiverpb.Packet) (*receiverpb.BaseResp, error) {
	// 获取 clientIP
	if err := srv.receiver.InsertPacket(ctx, srv.clientIP(ctx), in); err != nil {
		e, ok := err.(httputil.Error)
		if ok {
			// 数据库操作错误
			if e.Code == httputil.ErrCodeSystemException {
				return nil, httputil.ErrSystemException
			}
			return &receiverpb.BaseResp{
				Code:    int32(e.Code),
				Message: e.Message,
			}, nil
		}
		return nil, err
	}
	return &receiverpb.BaseResp{
		Code:    httputil.CodeSuccess,
		Message: "success",
	}, nil
}

func (srv *GRPCService) clientIP(ctx context.Context) string {
	ctxPeer, ok := peer.FromContext(ctx)
	if ok && ctxPeer.Addr != nil {
		if ipnet, ok := ctxPeer.Addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
				return ipnet.IP.String()
			}
		}
		// 上述解析不成功，则自行拼接
		return kit.AddrStringToIP(ctxPeer.Addr)
	}
	return ""
}
