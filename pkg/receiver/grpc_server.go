package receiver

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/huzhongqing/qelog/pkg/common/model"

	"google.golang.org/grpc/peer"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/libs/mongo"
	"github.com/huzhongqing/qelog/pb"
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

	if os.Getenv("ENV") == "release" {

	} else {
		if err := srv.database.UpsertCollectionIndexMany(
			model.ModuleMetricsIndexMany()); err != nil {
			return err
		}
	}

	server := grpc.NewServer()
	srv.server = server

	pb.RegisterPushServer(srv.server, srv)

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

func (srv *GRPCService) PushPacket(ctx context.Context, in *pb.Packet) (*pb.BaseResp, error) {
	fmt.Println(in)
	// 获取 clientIP
	if err := srv.receiver.InsertPacket(ctx, srv.clientIP(ctx), in); err != nil {
		e, ok := err.(httputil.Error)
		if ok {
			// 数据库操作错误
			if e.Code == httputil.ErrCodeSystemException {
				return nil, httputil.ErrSystemException
			}
			return &pb.BaseResp{
				Code:    int32(e.Code),
				Message: e.Message,
			}, nil
		}
		return nil, err
	}
	return &pb.BaseResp{
		Code:    httputil.CodeSuccess,
		Message: "success",
	}, nil
}

func (srv *GRPCService) clientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if ok && p.Addr != nil {
		addr := p.Addr.String()
		if strings.HasPrefix(addr, "[") {
			return strings.Split(strings.TrimPrefix(addr, "["), "]:")[0]
		}
		return strings.Split(addr, ":")[0]
	}
	return ""
}
