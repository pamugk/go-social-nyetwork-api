package rpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pamugk/social-nyetwork-server/internal/app"
	"github.com/pamugk/social-nyetwork-server/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServer struct {
	UnimplementedUserServiceServer
}

func (s *userServer) CreateUser(ctx context.Context, in *CreateUserRequest) (*emptypb.Empty, error) {
	if _, err := app.CreateUser(convertFromUserData(in.User), in.Password); err == nil {
		return &emptypb.Empty{}, nil
	} else {
		return &emptypb.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}
}

func (s *userServer) GetUser(ctx context.Context, in *GetUserRequest) (*GetUserResponse, error) {
	if user, err := app.GetUser(in.Id); err == nil {
		return convertToUser(user), nil
	} else if err.Error() == "Not found" {
		return nil, status.New(codes.NotFound, "User not found").Err()
	} else {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
}

func (s *userServer) SearchUsers(ctx context.Context, in *SearchUsersRequest) (*SearchUsersResponse, error) {
	if users, total, err := app.SearchUsers(in.LoginPart, in.Page, in.Limit); err == nil {
		return convertToSearchUserResponse(in.Page, in.Limit, total, users), nil
	} else {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
}

func (s *userServer) UpdateUser(ctx context.Context, in *UpdateUserRequest) (*emptypb.Empty, error) {
	if err := app.UpdateUser(in.Id, &domain.UserData{}); err == nil {
		return &emptypb.Empty{}, nil
	} else if err.Error() == "Not found" {
		return &emptypb.Empty{}, status.New(codes.NotFound, "User not found").Err()
	} else {
		return &emptypb.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}
}

func (s *userServer) ChangePassword(ctx context.Context, in *ChangePasswordRequest) (*emptypb.Empty, error) {
	if err := app.ChangePassword(in.Id, in.NewPassword); err == nil {
		return &emptypb.Empty{}, nil
	} else if err.Error() == "Not found" {
		return &emptypb.Empty{}, status.New(codes.NotFound, "User not found").Err()
	} else {
		return &emptypb.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}
}

func (s *userServer) DeleteUser(ctx context.Context, in *DeleteUserRequest) (*emptypb.Empty, error) {
	if err := app.DeleteUser(in.Id); err == nil {
		return &emptypb.Empty{}, nil
	} else if err.Error() == "Not found" {
		return &emptypb.Empty{}, status.New(codes.NotFound, "User not found").Err()
	} else {
		return &emptypb.Empty{}, status.New(codes.Internal, err.Error()).Err()
	}
}

func StartServer(port *int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	s := grpc.NewServer()
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	RegisterUserServiceServer(s, &userServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve RPC: %v", err)
	}
}
