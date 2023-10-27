package main

import (
	"context"
	"github.com/0xhoang/go-kit/common"
	"github.com/0xhoang/go-kit/internal/must"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/status"
	"net/http"
)

type middleware struct {
	TokenSecretKey string
}

func NewMiddleware(tokenSecretKey string) *middleware {
	return &middleware{TokenSecretKey: tokenSecretKey}
}

func (m *middleware) AuthMiddleware(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := must.ParseToken(token, m.TokenSecretKey)
	if err != nil {
		return nil, status.Errorf(http.StatusUnauthorized, "invalid auth token: %v", err)
	}

	return context.WithValue(ctx, common.CustomerKey, tokenInfo), nil
}
