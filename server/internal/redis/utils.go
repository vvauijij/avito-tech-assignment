package redis

import "fmt"

func getBannerKey(featureID int64, tagID int64) string {
	return fmt.Sprintf("%d:%d", featureID, tagID)
}
