package main

import (
	"ginskeleton/app/global/variable"
	grpc_interceptor "ginskeleton/app/grpc/server/interceptor"
	_ "ginskeleton/bootstrap"
	"ginskeleton/routers"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", variable.ConfigYml.GetString("GrpcServer.Port"))
	if err != nil {
		log.Fatalf("Tcp 监听失败: %v", err)
	}

	grpcServ := grpc.NewServer(grpc.UnaryInterceptor(grpc_interceptor.GrpcRequestLog()))

	routers.InitGrpcService(grpcServ)
	variable.ZapLog.Info("开始启动 grpc 服务, 监听端口: " + variable.ConfigYml.GetString("GrpcServer.Port"))

	if err = grpcServ.Serve(lis); err != nil {
		log.Fatalf("grpc 服务启动失败,错误: %v", err)
	}

}
