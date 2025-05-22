package model

type (
	User struct {
		ID         int64
		Name       string
		Email      string
		AvatarLink string
		Password   string

		Version int32
	}
)
