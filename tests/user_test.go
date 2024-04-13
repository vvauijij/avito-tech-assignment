package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserGetBannerNoBanner(t *testing.T) {
	CleanUp(t)

	assert.Equal(t, Content{}, GetBannerContent(t, 0, 0, false))
}

func TestUserGetBanner(t *testing.T) {
	CleanUp(t)

	banner := CreateBanner(t, 0, []int64{0}, Content{"banner_string": "string"}, true)
	content := GetBannerContent(t, 0, 0, false)
	assert.Equal(t, banner.Content, content)
}

func TestUserGetBannerComplexStructure(t *testing.T) {
	CleanUp(t)

	content := Content{
		"banner_float":  2.1,
		"banner_string": "string",
		"banner_bool":   false,
		"banner_array":  []any{2.1, "string", false},
		"banner_map": map[string]any{
			"float":  2.1,
			"string": "string",
			"bool":   false,
			"array":  []any{2.1, "string", false},
		},
	}
	banner := CreateBanner(t, 0, []int64{0}, content, true)
	content = GetBannerContent(t, 0, 0, false)
	assert.Equal(t, banner.Content, content)
}

func TestUserGetBannerFromCache(t *testing.T) {
	CleanUp(t)

	banner := CreateBanner(t, 0, []int64{0}, Content{"banner_string": "string"}, true)

	// Load content to cache
	GetBannerContent(t, 0, 0, false)

	UpdateBanner(t, *banner.BannerID, nil, nil, Content{"banner_string": "updated_string"}, nil)

	oldContent := GetBannerContent(t, 0, 0, false)
	newContent := GetBannerContent(t, 0, 0, true)
	assert.Equal(t, banner.Content, oldContent)
	assert.NotEqual(t, banner.Content, newContent)
}

func TestUserGetBannerNotActive(t *testing.T) {
	CleanUp(t)

	CreateBanner(t, 0, []int64{0}, Content{"banner_string": "string"}, false)
	assert.Equal(t, Content{}, GetBannerContent(t, 0, 0, false))
}
