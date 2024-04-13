package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vvauijij/avito-tech-assignment/tests/token"
)

func CleanUp(t *testing.T) {
	t.Helper()

	req, _ := http.NewRequest(http.MethodDelete, TestCleanUpEndpoint, nil)
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
}

func GetBannerContent(
	t *testing.T,
	featureID int64,
	tagID int64,
	useLastRevision bool,
) Content {

	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, UserBannerEndpoint(featureID, tagID, useLastRevision), nil)
	req.Header.Add("token", token.Generator.NewUserToken())

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	if resp.StatusCode == http.StatusNotFound {
		return Content{}
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, resp.ContentLength)
	read, err := resp.Body.Read(body)
	defer resp.Body.Close()

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, int(resp.ContentLength), read)

	content := Content{}
	err = json.Unmarshal(body, &content)
	assert.NoError(t, err)

	return content
}

func GetBanners(t *testing.T, filter Filter) []*Banner {
	t.Helper()

	req, _ := http.NewRequest(http.MethodGet, AdminBannerEndpoint, nil)
	req.Header.Add("token", token.Generator.NewAdminToken())
	q := req.URL.Query()
	for key, value := range filter {
		q.Set(key, fmt.Sprint(value))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, resp.ContentLength)
	read, err := resp.Body.Read(body)
	defer resp.Body.Close()

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, int(resp.ContentLength), read)

	banners := []*Banner{}
	err = json.Unmarshal(body, &banners)
	assert.NoError(t, err)

	return banners
}

func CreateBanner(
	t *testing.T,
	featureID int64,
	tagIDs []int64,
	content map[string]any,
	isActive bool,
) *Banner {

	t.Helper()

	banner := Banner{FeatureID: &featureID, TagIDs: tagIDs, Content: content, IsActive: &isActive}
	marshalled, err := json.Marshal(banner)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodPost, AdminBannerEndpoint, bytes.NewReader(marshalled))
	req.Header.Add("token", token.Generator.NewAdminToken())
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, resp.ContentLength)
	read, err := resp.Body.Read(body)
	defer resp.Body.Close()

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, int(resp.ContentLength), read)

	err = json.Unmarshal(body, &banner)
	assert.NoError(t, err)

	return &banner
}

func UpdateBanner(
	t *testing.T,
	bannerID int64,
	featureID *int64,
	tagIDs []int64,
	content map[string]any,
	isActive *bool,
) {

	t.Helper()

	banner := Banner{FeatureID: featureID, TagIDs: tagIDs, Content: content, IsActive: isActive}
	marshalled, err := json.Marshal(banner)
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodPatch, AdminBannerWithIDEndpoint(bannerID), bytes.NewReader(marshalled))
	req.Header.Add("token", token.Generator.NewAdminToken())
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func DeleteBanner(t *testing.T, bannerID int64) {
	t.Helper()

	req, _ := http.NewRequest(http.MethodDelete, AdminBannerWithIDEndpoint(bannerID), nil)
	req.Header.Add("token", token.Generator.NewAdminToken())

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func DeleteBanners(t *testing.T, filter Filter) {
	t.Helper()

	req, _ := http.NewRequest(http.MethodDelete, AdminBannerEndpoint, nil)
	req.Header.Add("token", token.Generator.NewAdminToken())
	q := req.URL.Query()
	for key, value := range filter {
		q.Set(key, fmt.Sprint(value))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
