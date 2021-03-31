package main

import (
	"context"

	"github.com/loadlab-go/usersvc/idl/proto/userpb"
)

type userSvc struct {
	userpb.UnimplementedUserServer
}

func (s *userSvc) Create(_ context.Context, _ *userpb.CreateRequest) (*userpb.CreateResponse, error) {
	panic("not implemented") // TODO: Implement
}
