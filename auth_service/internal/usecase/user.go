package usecase

import (
	"context"

	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	userpb "github.com/Temutjin2k/TyndauProto/gen/go/user"
	"google.golang.org/grpc"
)

// Provider with user gRPC server
type UserProvider struct {
	client userpb.UserClient
	conn   *grpc.ClientConn
}

func NewUserProvider(conn *grpc.ClientConn) *UserProvider {
	client := userpb.NewUserClient(conn)

	return &UserProvider{
		client: client,
		conn:   conn,
	}
}

// Create sends gRPC request to user service to create new user
func (u *UserProvider) Create(ctx context.Context, user model.User) (int64, error) {
	grpcReq := &userpb.CreateRequest{
		Name:       user.Name,
		Email:      user.Email,
		Password:   user.Password,
		AvatarLink: user.AvatarLink,
	}

	resp, err := u.client.Create(ctx, grpcReq)
	if err != nil {
		return 0, err
	}

	return resp.GetUserId(), nil
}
