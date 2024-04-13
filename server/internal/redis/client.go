package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

type ClientImpl struct {
	RedisClient     *redis.Client
	CacheExpireTime time.Duration
}

func NewClient(redisURI string, cacheExpireTime time.Duration) Client {
	clientOpts := &redis.Options{Addr: redisURI}
	client := redis.NewClient(clientOpts)
	return &ClientImpl{RedisClient: client, CacheExpireTime: cacheExpireTime}
}

func (c *ClientImpl) GetOne(ctx context.Context, featureID int64, tagID int64) (*types.Banner, error) {
	BannerID, err := c.RedisClient.Get(
		ctx,
		getBannerKey(featureID, tagID),
	).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrBannerNotFound
	}
	if err != nil {
		return nil, err
	}

	bannerValue, err := c.RedisClient.Get(
		ctx,
		BannerID,
	).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrBannerNotFound
	}
	if err != nil {
		return nil, err
	}

	banner := &types.Banner{}
	err = json.Unmarshal([]byte(bannerValue), banner)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (c *ClientImpl) Update(ctx context.Context, banner *types.Banner) error {
	BannerID := strconv.FormatInt(banner.BannerID, 10)
	bannerValue, err := json.Marshal(banner)
	if err != nil {
		return err
	}
	c.RedisClient.Set(
		ctx,
		BannerID,
		bannerValue,
		c.CacheExpireTime,
	)

	for _, tagID := range banner.TagIDs {
		c.RedisClient.Set(
			ctx,
			getBannerKey(*banner.FeatureID, tagID),
			BannerID,
			c.CacheExpireTime,
		)
	}
	return nil
}

func (c *ClientImpl) TestCleanUp(ctx context.Context) error {
	return c.RedisClient.FlushDB(ctx).Err()
}
