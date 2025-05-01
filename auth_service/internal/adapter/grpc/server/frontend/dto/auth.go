package dto

import (
	"github.com/Temutjin2k/Tyndau/auth_service/internal/model"
	authpb "github.com/Temutjin2k/TyndauProto/gen/go/auth"
)

func FromRegisterRequest(req *authpb.RegisterRequest) model.User {
	user := model.User{
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		AvatarLink: req.GetAvatarLink(),
	}

	return user
}
