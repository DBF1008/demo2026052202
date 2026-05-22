package service_implement

import (
	"context"

	"ginskeleton/app/grpc/server/proto/user_demo_pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/anypb"
)

type UserService struct {

}

func (u *UserService) GetUserInfo(ctx context.Context, req *user_demo_pb.UserRequest) (resp *user_demo_pb.UserResponse, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "GetUserInfo - 接口请求的参数校验不通过: %v", err)
	}
	resp = &user_demo_pb.UserResponse{
		Code: 200,
		Msg:  "success",
		Data: []*user_demo_pb.Data2{
			&user_demo_pb.Data2{
				OneItem: map[string]*anypb.Any{
					"user_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("zhangsan001"),
					},
					"age": {
						TypeUrl: "type.googleapis.com/google.protobuf.Int32Value",
						Value:   []byte("18"),
					},
					"real_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("张三2026"),
					},
				},
			},
			&user_demo_pb.Data2{
				OneItem: map[string]*anypb.Any{
					"user_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("lisi002"),
					},
					"age": {
						TypeUrl: "type.googleapis.com/google.protobuf.Int32Value",
						Value:   []byte("19"),
					},
					"real_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("李四2026"),
					},
				},
			},
		},
	}
	return resp, nil
}

func (u *UserService) GetItem(ctx context.Context, req *user_demo_pb.UserIdRequest) (resp *user_demo_pb.UserResponse, err error) {
	if err = req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "GetItem - 接口请求的参数校验不通过: %v", err)
	}

	resp = &user_demo_pb.UserResponse{
		Code: 200,
		Msg:  "success",
		Data: []*user_demo_pb.Data2{
			&user_demo_pb.Data2{
				OneItem: map[string]*anypb.Any{
					"user_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("zhangsan"),
					},
					"age": {
						TypeUrl: "type.googleapis.com/google.protobuf.Int32Value",
						Value:   []byte("18"),
					},
					"real_name": {
						TypeUrl: "type.googleapis.com/google.protobuf.StringValue",
						Value:   []byte("张三2026"),
					},
				},
			},
		},
	}
	return resp, nil
}
