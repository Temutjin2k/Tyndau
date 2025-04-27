package frontend

import userv1 "github.com/Temutjin2k/TyndauProto/gen/go/user"

type UserV1 struct {
	userv1.UnimplementedAuthServer
	uc UserUseCase
}

func NewUserV1(uc UserUseCase) *UserV1 {
	return &UserV1{uc: uc}
}
