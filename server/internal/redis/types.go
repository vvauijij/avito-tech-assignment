package redis

import (
	"context"

	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

type (
	Client interface {
		GetOne(ctx context.Context, featureID int64, tagID int64) (*types.Banner, error)

		Update(ctx context.Context, banner *types.Banner) error

		TestCleanUp(ctx context.Context) error
	}
)
