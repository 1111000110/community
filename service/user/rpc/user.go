package main

import (
	"flag"
	"fmt"

	"community/service/user/rpc/internal/config"
	relationsserviceServer "community/service/user/rpc/internal/server/relationsservice"
	userserviceServer "community/service/user/rpc/internal/server/userservice"
	"community/service/user/rpc/internal/svc"
	"community/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterUserServiceServer(grpcServer, userserviceServer.NewUserServiceServer(ctx))
		__.RegisterRelationsServiceServer(grpcServer, relationsserviceServer.NewRelationsServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
