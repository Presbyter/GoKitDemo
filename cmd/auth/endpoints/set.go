package endpoints

import (
	"GoKitDemo/cmd/auth/service"
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"
)

type Set struct {
	GetJwtEndpoint   endpoint.Endpoint
	RefreshEndpoint  endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger) Set {
	var getJwtEndpoint endpoint.Endpoint
	{
		getJwtEndpoint = MakeGetJwtEndpoint(svc)
		getJwtEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100000))(getJwtEndpoint)
		getJwtEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{Timeout: time.Second}))(getJwtEndpoint)
		getJwtEndpoint = loggingMiddleware(log.With(logger, "method", "GetJwt"))(getJwtEndpoint)
	}

	var refreshEndpoint endpoint.Endpoint
	{
		refreshEndpoint = MakeRefreshEndpoint(svc)
		refreshEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100000))(refreshEndpoint)
		refreshEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{Timeout: time.Second}))(refreshEndpoint)
		refreshEndpoint = loggingMiddleware(log.With(logger, "method", "Refresh"))(refreshEndpoint)
	}

	var validateEndpoint endpoint.Endpoint
	{
		validateEndpoint = MakeValidateEndpoint(svc)
		validateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100000))(validateEndpoint)
		validateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{Timeout: time.Second}))(validateEndpoint)
		validateEndpoint = loggingMiddleware(log.With(logger, "method", "Validate"))(validateEndpoint)
	}

	return Set{
		GetJwtEndpoint:   getJwtEndpoint,
		RefreshEndpoint:  refreshEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}

func MakeGetJwtEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetJwtRequest)
		t, rt, err := svc.GetJwt(ctx, req.UserID, time.Duration(req.ValidityTime), req.Aud, req.Sub)
		resp := &GetJwtResponse{
			Token:        t,
			RefreshToken: rt,
			Err:          err,
		}
		return resp, nil
	}
}

type GetJwtRequest struct {
	UserID       int
	ValidityTime int64
	Aud          []string
	Sub          string
}

type GetJwtResponse struct {
	Token        string
	RefreshToken string
	Err          error
}

func MakeRefreshEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RefreshRequest)
		t, rt, err := svc.Refresh(ctx, req.RefreshToken)
		resp := &RefreshResponse{
			Token:        t,
			RefreshToken: rt,
			Err:          err,
		}
		return resp, nil
	}
}

type RefreshRequest struct {
	RefreshToken string
}

type RefreshResponse struct {
	Token        string
	RefreshToken string
	Err          error
}

func MakeValidateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidateRequest)
		pl, err := svc.Validate(ctx, req.Token)
		resp := &ValidateResponse{
			CustomPayload: pl,
			err:           err,
		}
		return resp, nil
	}
}

type ValidateRequest struct {
	Token string
}

type ValidateResponse struct {
	CustomPayload *service.CustomPayload
	err           error
}
