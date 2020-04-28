package endpoints

import (
	"GoKitDemo/cmd/user/service"
	repo "GoKitDemo/repositories"
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	LoginEndpoint endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger /*, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer*/) Set {
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = MakeLoginEndpoint(svc)
		loginEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{Timeout: time.Second}))(loginEndpoint)
		loginEndpoint = LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
	}

	return Set{
		LoginEndpoint: loginEndpoint,
	}
}

func MakeLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginRequest)
		v, err := svc.Login(ctx, req.LoginType, req.Value, req.DeviceType, req.DeviceCode)
		return LoginResponse{V: v, Err: err}, nil
	}
}

type LoginRequest struct {
	LoginType  int32
	Value      string
	Code       string
	DeviceType int32
	DeviceCode string
}

type LoginResponse struct {
	V   *repo.User `json:"v"`
	Err error      `json:"-"`
}

func (l LoginResponse) Failed() error { return l.Err }
