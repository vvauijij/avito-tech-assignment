package tests

import (
	"fmt"
	"os"
)

var (
	serverURL           = os.Getenv("SERVER")
	userBannerEndpoint  = fmt.Sprintf("%s/user_banner", serverURL)
	AdminBannerEndpoint = fmt.Sprintf("%s/banner", serverURL)
	TestCleanUpEndpoint = fmt.Sprintf("%s/test_clean_up", serverURL)
)

func AdminBannerWithIDEndpoint(bannerID int64) string {
	return fmt.Sprintf("%s/%d", AdminBannerEndpoint, bannerID)
}

func UserBannerEndpoint(featureID int64, tagID int64, useLastRevision bool) string {
	return fmt.Sprintf(
		"%s?tag_id=%d&feature_id=%d&use_last_revision=%t",
		userBannerEndpoint,
		featureID,
		tagID,
		useLastRevision,
	)
}
