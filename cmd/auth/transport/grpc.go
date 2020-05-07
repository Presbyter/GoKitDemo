package transport

import (
	"GoKitDemo/cmd/auth/endpoints"
	"GoKitDemo/pb"
	"context"
	"github.com/go-kit/kit/log"
	t "github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	getJwt   grpctransport.Handler
	refresh  grpctransport.Handler
	validate grpctransport.Handler
}

func NewGRPCServer(endpoints endpoints.Set, logger log.Logger) pb.AuthServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(t.NewLogErrorHandler(logger)),
	}
	return &grpcServer{
		getJwt: grpctransport.NewServer(
			endpoints.GetJwtEndpoint,
			decodeGRPCGetJwtRequest,
			encodedGRPCGetJwtResponse,
			options...),
		refresh: grpctransport.NewServer(
			endpoints.RefreshEndpoint,
			decodeGRPCRefreshRequest,
			encodedGRPCRefreshResponse,
			options...),
		validate: grpctransport.NewServer(
			endpoints.ValidateEndpoint,
			decodeGRPCValidateRequest,
			encodedGRPCValidateResponse,
			options...),
	}
}

func encodedGRPCValidateResponse(_ context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(*endpoints.ValidateResponse)
	if resp.Err == nil {
		audJson, err := resp.CustomPayload.Audience.MarshalJSON()
		if err != nil {
			return nil, err
		}
		response = &pb.ValidateResponse{
			Code:   0,
			ErrMsg: "OK",
			Data: &pb.ValidateResponse_DataStruct{
				UserId: int32(resp.CustomPayload.UserID),
				Iss:    resp.CustomPayload.Issuer,
				Exp:    int32(resp.CustomPayload.ExpirationTime.Unix()),
				Sub:    resp.CustomPayload.Subject,
				Aud:    string(audJson),
				Nbf:    int32(resp.CustomPayload.NotBefore.Unix()),
				Iat:    int32(resp.CustomPayload.IssuedAt.Unix()),
				Jti:    resp.CustomPayload.JWTID,
			},
		}
		return
	}
	return &pb.ValidateResponse{
		Code:   500,
		ErrMsg: resp.Err.Error(),
		Data:   nil,
	}, nil
}

func decodeGRPCValidateRequest(_ context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.ValidateRequest)
	request = endpoints.ValidateRequest{
		Token: req.GetToken(),
	}
	return
}

func encodedGRPCRefreshResponse(_ context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(*endpoints.RefreshResponse)
	if resp.Err == nil {
		response = &pb.RefreshResponse{
			Code:   0,
			ErrMsg: "OK",
			Data: &pb.RefreshResponse_DataStruct{
				Token:        resp.Token,
				RefreshToken: resp.RefreshToken,
			},
		}
		return
	}
	return &pb.RefreshResponse{
		Code:   500,
		ErrMsg: resp.Err.Error(),
		Data:   nil,
	}, nil
}

func decodeGRPCRefreshRequest(_ context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.RefreshRequest)
	request = endpoints.RefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	}
	return
}

func encodedGRPCGetJwtResponse(_ context.Context, i interface{}) (response interface{}, err error) {
	resp := i.(*endpoints.GetJwtResponse)
	if resp.Err == nil {
		response = &pb.GetJwtResponse{
			Code:   0,
			ErrMsg: "OK",
			Data: &pb.GetJwtResponse_DataStruct{
				Token:        resp.Token,
				RefreshToken: resp.RefreshToken,
			},
		}
		return
	}
	return &pb.GetJwtResponse{
		Code:   500,
		ErrMsg: resp.Err.Error(),
		Data:   nil,
	}, nil
}

func decodeGRPCGetJwtRequest(_ context.Context, i interface{}) (request interface{}, err error) {
	req := i.(*pb.GetJwtRequest)
	request = endpoints.GetJwtRequest{
		UserID:       int(req.GetUserId()),
		ValidityTime: req.GetValidityTime(),
		Aud:          req.GetAud(),
		Sub:          req.GetSub(),
	}
	return
}

func (g grpcServer) GetJwt(ctx context.Context, request *pb.GetJwtRequest) (*pb.GetJwtResponse, error) {
	_, resp, err := g.getJwt.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetJwtResponse), nil
}

func (g grpcServer) Refresh(ctx context.Context, request *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	_, resp, err := g.refresh.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.RefreshResponse), nil
}

func (g grpcServer) Validate(ctx context.Context, request *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, resp, err := g.validate.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ValidateResponse), nil
}
