package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminGetBanners(t *testing.T) {
	CleanUp(t)

	banners := []*Banner{CreateBanner(t, int64(0), []int64{int64(0)}, Content{"banner_string": "string"}, true)}
	assert.Equal(t, banners, GetBanners(t, Filter{}))
}

func TestAdminGetBannersWithFilter(t *testing.T) {
	CleanUp(t)

	nBanners := 10
	banners := make([]*Banner, 0, nBanners)
	for i := 0; i < nBanners; i++ {
		banners = append(banners, CreateBanner(t, int64(i), []int64{int64(i)}, Content{"banner_string": "string"}, true))
	}

	for i := 0; i < nBanners; i++ {
		assert.Equal(t, banners[i], GetBanners(t, Filter{"feature_id": i, "tag_id": i})[0])
	}
}

func TestAdminGetBannersWithLimitAndOffset(t *testing.T) {
	CleanUp(t)

	CreateBanner(t, int64(0), []int64{int64(0)}, Content{"banner_string": "string"}, true)
	CreateBanner(t, int64(1), []int64{int64(1)}, Content{"banner_string": "string"}, true)

	bannersFirstHalf := GetBanners(t, Filter{"limit": 1, "offset": 0})
	bannersSecondHalf := GetBanners(t, Filter{"limit": 1, "offset": 1})

	assert.Equal(t, 1, len(bannersFirstHalf))
	assert.Equal(t, 1, len(bannersSecondHalf))
	assert.NotEqual(t, bannersFirstHalf, bannersSecondHalf)
}

func TestAdminUpdateBanner(t *testing.T) {
	CleanUp(t)

	banner := CreateBanner(t, int64(0), []int64{int64(0)}, Content{"banner_string": "string"}, true)

	updatedFeatureID := int64(1)
	updatedTagIDs := []int64{int64(1)}
	updatedContent := Content{"banner_string": "updated_string"}
	updatedIsActive := false
	UpdateBanner(t, *banner.BannerID, &updatedFeatureID, updatedTagIDs, updatedContent, &updatedIsActive)

	updatedBanner := GetBanners(t, Filter{})[0]
	assert.NotEqual(t, banner, updatedBanner)

	banner.FeatureID = &updatedFeatureID
	banner.TagIDs = updatedTagIDs
	banner.Content = updatedContent
	banner.IsActive = &updatedIsActive
	assert.Equal(t, banner, updatedBanner)
}

func TestAdminDeleteBanner(t *testing.T) {
	CleanUp(t)

	banner := CreateBanner(t, int64(0), []int64{int64(0)}, Content{"banner_string": "string"}, true)
	DeleteBanner(t, *banner.BannerID)
	assert.Equal(t, []*Banner{}, GetBanners(t, Filter{}))
}

func TestAdminDeleteBannersWithFilter(t *testing.T) {
	CleanUp(t)

	nBannersWithZeroTag := 10
	for i := 0; i < nBannersWithZeroTag; i++ {
		CreateBanner(t, int64(i), []int64{int64(0), int64(1), int64(2)}, Content{"banner_string": "string"}, true)
	}
	nBannersWithoutZeroTag := 10
	for i := 0; i < nBannersWithoutZeroTag; i++ {
		CreateBanner(t, int64(i), []int64{int64(1), int64(2), int64(3)}, Content{"banner_string": "string"}, true)
	}

	DeleteBanners(t, Filter{"tag_id": 0})
	assert.Equal(t, 0, len(GetBanners(t, Filter{"tag_id": 0})))
	assert.Equal(t, nBannersWithoutZeroTag, len(GetBanners(t, Filter{})))
}
