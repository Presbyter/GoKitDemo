package service

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"time"
)

type Middleware func(Service) Service

func loggingMiddleWare(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (l loggingMiddleware) GetJwt(ctx context.Context, userID int, validityTime time.Duration, aud []string, sub string) (token, refreshToken string, err error) {
	start := time.Now()
	defer func() {
		l.logger.Log("Method", "GetJwt",
			"userID", userID,
			"validityTime", validityTime,
			"aud", aud,
			"sub", sub,
			"useTime", fmt.Sprintf("%.2f ms", float64(time.Since(start))/float64(time.Millisecond)))
	}()
	return l.next.GetJwt(ctx, userID, validityTime, aud, sub)
}

func (l loggingMiddleware) Refresh(ctx context.Context, refToken string) (token, refreshToken string, err error) {
	start := time.Now()
	defer func() {
		l.logger.Log("Method", "Refresh",
			"refToken", refToken,
			"useTime", fmt.Sprintf("%.2f ms", float64(time.Since(start))/float64(time.Millisecond)))
	}()
	return l.next.Refresh(ctx, refToken)
}

func (l loggingMiddleware) Validate(ctx context.Context, token string) (payload *CustomPayload, err error) {
	start := time.Now()
	defer func() {
		l.logger.Log("Method", "Validate",
			"token", token,
			"useTime", fmt.Sprintf("%.2f ms", float64(time.Since(start))/float64(time.Millisecond)))
	}()
	return l.next.Validate(ctx, token)
}
