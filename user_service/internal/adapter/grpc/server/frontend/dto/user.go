package dto

import (
	"github.com/Temutjin2k/Tyndau/user_service/internal/model"
	userpb "github.com/Temutjin2k/TyndauProto/gen/go/user"
)

func FromCreateRequest(req *userpb.CreateRequest) model.User {
	return model.User{
		Name:       req.GetName(),
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		AvatarLink: req.GetAvatarLink(),
	}
}

func FromUpdateRequest(req *userpb.UpdateRequest) model.User {
	return model.User{
		ID:         req.GetUserId(),
		Name:       req.GetName(),
		AvatarLink: req.GetAvatarLink(),
		Version:    req.GetVersion(),
	}
}

func ToProfileResponce(user model.User) *userpb.ProfileResponse {
	return &userpb.ProfileResponse{
		UserId:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		AvatarLink: user.AvatarLink,
		Version:    user.Version,
	}
}

func ToUpdateResponce(user model.User) *userpb.UpdateResponse {
	return &userpb.UpdateResponse{
		UserId:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		AvatarLink: user.AvatarLink,
		Version:    user.Version,
	}
}
