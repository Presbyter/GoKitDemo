package service

import (
	"GoKitDemo/util"
	"context"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-kit/kit/log"
	"time"
)

type Service interface {
	GetJwt(ctx context.Context, userID int, validityTime time.Duration, aud []string, sub string) (token, refreshToken string, err error)
	Refresh(ctx context.Context, refToken string) (token, refreshToken string, err error)
	Validate(ctx context.Context, token string) (payload *CustomPayload, err error)
}

func New(logger log.Logger) Service {
	var svc Service
	{
		svc = newAuthService()
		svc = loggingMiddleWare(logger)(svc)
	}
	return svc
}

var (
	hs = jwt.NewHS256([]byte("81f92fb7-b489-4815-b06a-77c2cec57caf_presbyter_secret_key"))
)

type CustomPayload struct {
	jwt.Payload
	UserID int `json:"user_id"`
}

type authService struct{}

func newAuthService() Service {
	return &authService{}
}

func (a authService) GetJwt(ctx context.Context, userID int, validityTime time.Duration, aud []string, sub string) (string, string, error) {
	now := time.Now()
	pl := &CustomPayload{
		Payload: jwt.Payload{
			Issuer:         "auth_service",
			Subject:        sub,
			Audience:       aud,
			ExpirationTime: jwt.NumericDate(now.Add(validityTime)),
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          fmt.Sprintf("%d%s", now.Unix(), util.RandString(6)),
		},
		UserID: userID,
	}
	token, err := jwt.Sign(pl, hs)
	if err != nil {
		return "", "", err
	}

	pl.Payload.ExpirationTime = jwt.NumericDate(now.Add(validityTime + 30*24*time.Hour))
	pl.Payload.Subject = "refresh_token"
	pl.Payload.JWTID = fmt.Sprintf("%d%s", now.Unix(), util.RandString(6))
	refToken, err := jwt.Sign(pl, hs)
	if err != nil {
		return "", "", err
	}

	return string(token), string(refToken), nil
}

func (a authService) Refresh(ctx context.Context, refToken string) (string, string, error) {
	now := time.Now()
	subValidator := jwt.SubjectValidator("refresh_token")
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	nbfValidator := jwt.NotBeforeValidator(now)

	var pl CustomPayload
	validatePayload := jwt.ValidatePayload(&pl.Payload, subValidator, iatValidator, expValidator, nbfValidator)

	_, err := jwt.Verify([]byte(refToken), hs, &pl, validatePayload)
	if err != nil {
		return "", "", err
	}

	return a.GetJwt(ctx, pl.UserID, pl.ExpirationTime.Sub(pl.IssuedAt.Time), pl.Audience, pl.Subject)
}

func (a authService) Validate(ctx context.Context, token string) (*CustomPayload, error) {
	now := time.Now()
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	nbfValidator := jwt.NotBeforeValidator(now)

	var pl CustomPayload
	validatePayload := jwt.ValidatePayload(&pl.Payload, iatValidator, expValidator, nbfValidator)

	_, err := jwt.Verify([]byte(token), hs, &pl, validatePayload)
	if err != nil {
		return nil, err
	}
	return &pl, nil
}
