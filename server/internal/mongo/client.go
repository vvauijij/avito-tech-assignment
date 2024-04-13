package mongo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

type ClientImpl struct {
	BannerColl *mongo.Collection
}

func NewClient(mongoURI string) Client {
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return &ClientImpl{BannerColl: client.Database("BannerDB").Collection("BannerColl")}
}

func (c *ClientImpl) GetOne(
	ctx context.Context,
	featureID int64,
	tagID int64,
) (banner *types.Banner, err error) {

	filter := bson.D{
		{Key: "feature_id", Value: featureID},
		{Key: "tag_ids", Value: tagID},
	}
	banner = &types.Banner{}
	err = c.BannerColl.FindOne(
		ctx,
		filter,
	).Decode(banner)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &types.Banner{}, ErrBannerNotFound
	}
	return banner, err
}

func (c *ClientImpl) GetMany(
	ctx context.Context,
	filter *types.BannerFilter,
	limit int64,
	offset int64,
) (banners []*types.Banner, err error) {

	options := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	}
	res, err := c.BannerColl.Find(
		ctx,
		filter,
		&options,
	)
	if err != nil {
		return []*types.Banner{}, err
	}
	banners = make([]*types.Banner, res.RemainingBatchLength())
	err = res.All(ctx, &banners)
	return banners, err
}

func (c *ClientImpl) Create(ctx context.Context, banner *types.Banner) (int64, error) {
	banner.BannerID = rand.Int63()
	createdAt := time.Now().Format(time.RFC3339)
	banner.CreatedAt = &createdAt
	banner.UpdatedAt = &createdAt

	_, err := c.BannerColl.InsertOne(ctx, banner)
	return banner.BannerID, err
}

func (c *ClientImpl) Update(ctx context.Context, banner *types.Banner) error {
	updatedAt := time.Now().Format(time.RFC3339)
	banner.UpdatedAt = &updatedAt

	filter := bson.D{
		{Key: "banner_id", Value: banner.BannerID},
	}
	update := bson.D{
		{Key: "$set", Value: banner},
	}
	err := c.BannerColl.FindOneAndUpdate(
		ctx,
		filter,
		update,
	).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrBannerNotFound
	}
	return err
}

func (c *ClientImpl) DeleteOne(ctx context.Context, id int64) error {
	filter := bson.D{
		{Key: "banner_id", Value: id},
	}
	res, err := c.BannerColl.DeleteOne(
		ctx,
		filter,
	)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrBannerNotFound
	}
	return nil
}

func (c *ClientImpl) DeleteMany(
	ctx context.Context,
	filter *types.BannerFilter,
) {

	_, _ = c.BannerColl.DeleteMany(
		ctx,
		filter,
	)
}

func (c *ClientImpl) TestCleanUp(ctx context.Context) error {
	filter := bson.D{}
	_, err := c.BannerColl.DeleteMany(ctx, filter)
	return err
}
