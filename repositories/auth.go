package repositories

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type authRepo struct {
	rds redis.UniversalClient
}

func NewAuthRepo() *authRepo {
	return &authRepo{rds: rdsClient}
}

func (a *authRepo) Get(sn string) (string, error) {
	key := fmt.Sprintf("auth:%s", sn)
	return a.rds.Get(key).Result()
}

func (a *authRepo) Set(sn, value string, exp time.Duration) error {
	key := fmt.Sprintf("auth:%s", sn)
	return a.rds.Set(key, value, exp).Err()
}
