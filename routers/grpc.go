package routers

import (
	"ginskeleton/app/grpc/server/proto/stu_demo_pb"
	"ginskeleton/app/grpc/server/proto/user_demo_pb"
	"ginskeleton/app/grpc/server/service_implement"

	"google.golang.org/grpc"
)

func InitGrpcService(grpcServ *grpc.Server) {

	stu_demo_pb.RegisterStudentServiceServer(grpcServ, &service_implement.StuService{})

	user_demo_pb.RegisterUserServiceServer(grpcServ, &service_implement.UserService{})

}
