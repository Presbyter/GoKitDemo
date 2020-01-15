package transport

import (
	"GoKitDemo/cmd/user/endpoints"
	"GoKitDemo/pb"
	"context"

	"github.com/go-kit/kit/log"
	t "github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	login grpctransport.Handler
}

func NewGRPCServer(endpoints endpoints.Set, logger log.Logger) pb.UserServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(t.NewLogErrorHandler(logger)),
	}
	return &grpcServer{
		login: grpctransport.NewServer(
			endpoints.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
			options...,
		),
	}
}

func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoints.LoginResponse)
	if resp.Err != nil {
		return &pb.LoginResponse{
			Code: 500,
			Msg:  resp.Err.Error(),
			Data: nil,
		}, nil
	} else {
		return &pb.LoginResponse{
			Code: 0,
			Msg:  "SUCCESS",
			Data: &pb.LoginResponse_DataStruct{
				UserID:   int32(resp.V.ID),
				NickName: resp.V.NickName,
				Gender:   int32(resp.V.Gender),
				Avatar:   resp.V.Avatar,
			},
		}, nil
	}
}

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return endpoints.LoginRequest{
		LoginType:  int32(req.GetLoginType()),
		Value:      req.GetValue(),
		Code:       req.GetCode(),
		DeviceType: int32(req.GetDeviceType()),
		DeviceCode: req.GetDeviceCode(),
	}, nil
}

func (g *grpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := g.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LoginResponse), nil
}
