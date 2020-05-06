package endpoints

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"time"
)

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("transport_error", err,
					"took", fmt.Sprintf("%.2f ms", float64(time.Since(begin))/float64(time.Millisecond)))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
