package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

func JSONError(w http.ResponseWriter, err error, HTTPStatus int) {
	w.WriteHeader(HTTPStatus)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]error{"error": err})
}

func JSONBanner(w http.ResponseWriter, banner *types.Banner) {
	if banner.IsActive != nil && *banner.IsActive {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(banner.Content)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func JSONBanners(w http.ResponseWriter, banners []*types.Banner) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(banners)
}

func JSONBannerID(w http.ResponseWriter, bannerID int64) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int64{"banner_id": bannerID})
}
