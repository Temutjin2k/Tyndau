package usecase

import "context"

type SongRepository interface {
	Create(ctx context.Context)
}
