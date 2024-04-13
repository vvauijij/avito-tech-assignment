package mongo

import (
	"context"

	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

type (
	Client interface {
		GetOne(ctx context.Context, featureID int64, tagID int64) (*types.Banner, error)

		GetMany(ctx context.Context, filter *types.BannerFilter, limit int64, offset int64) ([]*types.Banner, error)

		Create(ctx context.Context, banner *types.Banner) (int64, error)

		Update(ctx context.Context, banner *types.Banner) error

		DeleteOne(ctx context.Context, id int64) error

		DeleteMany(ctx context.Context, filter *types.BannerFilter)

		TestCleanUp(ctx context.Context) error
	}
)
