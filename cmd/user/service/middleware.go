package service

import (
	repo "GoKitDemo/repositories"
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (l loggingMiddleware) Login(ctx context.Context, loginType int32, value string, deviceType int32, deviceCode string) (*repo.User, error) {
	defer func() {
		l.logger.Log("method", "Login",
			"loginType", loginType,
			"value", value,
			"deviceType", deviceType,
			"deviceCode", deviceCode)
	}()
	return l.next.Login(ctx, loginType, value, deviceType, deviceCode)
}

func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

func (i instrumentingMiddleware) Login(ctx context.Context, loginType int32, value string, deviceType int32, deviceCode string) (*repo.User, error) {
	v, err := i.next.Login(ctx, loginType, value, deviceType, deviceCode)
	//i.ints.Add()
	return v, err
}
