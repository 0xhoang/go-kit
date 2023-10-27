package services

import (
	"context"
	"github.com/0xhoang/go-kit/gen"
)

func (e *GokitService) Profile(ctx context.Context, in *gen.EmptyRequest) (*gen.UserInfoResponse, error) {
	user, err := e.userFromContext(ctx)
	if err != nil {
		return nil, err
	}

	resp := &gen.UserInfoResponse{
		ID:        int64(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return resp, nil
}
