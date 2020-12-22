package receiver

import (
	"context"
	"net"

	"google.golang.org/grpc/reflection"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/libs/mongo"
	"github.com/huzhongqing/qelog/pkg/common/push"
	"github.com/huzhongqing/qelog/pkg/storage"
	"google.golang.org/grpc"
)

type GRPCService struct {
	database *mongo.Database
	server   *grpc.Server
	receiver *Service
}

func NewGRPCService(database *mongo.Database) *GRPCService {
	srv := &GRPCService{
		database: database,
		server:   nil,
		receiver: NewService(storage.New(database)),
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

	push.RegisterPushServer(srv.server, srv)

	reflection.Register(srv.server)

	if err := server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (srv *GRPCService) Close() error {
	if srv.server != nil {
		srv.server.Stop()
	}

	return nil
}

func (srv *GRPCService) PushPacket(ctx context.Context, in *push.Packet) (*push.BaseResp, error) {
	// 获取 clientIP
	if err := srv.receiver.InsertPacket(ctx, "127.0.0.2", in); err != nil {
		return nil, httputil.ErrSystemException
	}

	return &push.BaseResp{
		Code:    httputil.CodeSuccess,
		Message: "success",
	}, nil
}
