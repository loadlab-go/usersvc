package main

import (
	"context"

	userpb "github.com/loadlab-go/pkg/proto/user"
	"github.com/loadlab-go/usersvc/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userSvc struct {
	userpb.UnimplementedUserServer
	repo userRepo
}

type userRepo interface {
	CreateUser(username, password string) (*model.User, error)
}

func (s *userSvc) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	u, err := s.repo.CreateUser(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "create user failed: %v", err)
	}
	return &userpb.CreateResponse{Id: int64(u.ID)}, nil
}
