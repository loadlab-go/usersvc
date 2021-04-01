package main

import (
	"context"

	userpb "github.com/loadlab-go/pkg/proto/user"
	"github.com/loadlab-go/usersvc/model"
	"github.com/loadlab-go/usersvc/passwd"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userSvc struct {
	userpb.UnimplementedUserServer
	repo userRepo
}

type userRepo interface {
	CreateUser(username, password string) (*model.User, error)
	GetUser(id uint) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
}

func (s *userSvc) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	hashed, err := passwd.Generate(req.Password)
	if err != nil {
		logger.Warn("generate password failed", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, "generate password failed: %v", err)
	}
	u, err := s.repo.CreateUser(req.Username, hashed)
	if err != nil {
		logger.Warn("create user failed", zap.Error(err))
		return nil, status.Errorf(codes.Aborted, "create user failed: %v", err)
	}
	return &userpb.CreateResponse{Id: int64(u.ID)}, nil
}

func (s *userSvc) Get(_ context.Context, req *userpb.GetRequest) (*userpb.GetResponse, error) {
	u, err := s.repo.GetUser(uint(req.Id))
	if err != nil {
		logger.Warn("get user failed", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "get user failed: %v", err)
	}
	return &userpb.GetResponse{
		Id:       int64(u.ID),
		Username: u.Name,
	}, nil
}

func (s *userSvc) ValidatePassword(_ context.Context, req *userpb.ValidatePasswordRequest) (*userpb.ValidatePasswordResponse, error) {
	u, err := s.repo.GetUserByName(req.Username)
	if err != nil {
		logger.Warn("get user failed", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "get user failed: %v", err)
	}
	err = passwd.Validate(u.Password, req.Password)
	if err != nil {
		logger.Warn("validate password failed", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "wrong password: %v", err)
	}
	return &userpb.ValidatePasswordResponse{Id: int64(u.ID)}, nil
}
