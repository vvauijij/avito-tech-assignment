package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/vvauijij/avito-tech-assignment/server/internal/mongo"
	"github.com/vvauijij/avito-tech-assignment/server/internal/redis"
	"github.com/vvauijij/avito-tech-assignment/server/internal/token"
	"github.com/vvauijij/avito-tech-assignment/server/internal/types"
)

type HTTPHandlerImpl struct {
	TokenClient token.Client
	RedisClient redis.Client
	MongoClient mongo.Client
}

func NewHTTPHandler(
	tokenClient token.Client,
	redisClient redis.Client,
	mongoClient mongo.Client,
) HTTPHandler {
	return &HTTPHandlerImpl{
		TokenClient: tokenClient,
		RedisClient: redisClient,
		MongoClient: mongoClient,
	}
}

func (h *HTTPHandlerImpl) UserBanner(w http.ResponseWriter, req *http.Request) {
	t := req.Header.Get("token")
	switch h.TokenClient.Verify(t) {
	case token.None:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		h.getBanner(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandlerImpl) AdminBanner(w http.ResponseWriter, req *http.Request) {
	t := req.Header.Get("token")
	switch h.TokenClient.Verify(t) {
	case token.User:
		w.WriteHeader(http.StatusForbidden)
		return
	case token.None:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		h.getBanners(w, req)
	case http.MethodPost:
		h.createBanner(w, req)
	case http.MethodDelete:
		h.deleteBanners(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandlerImpl) AdminBannerWithID(w http.ResponseWriter, req *http.Request) {
	t := req.Header.Get("token")
	switch h.TokenClient.Verify(t) {
	case token.User:
		w.WriteHeader(http.StatusForbidden)
		return
	case token.None:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodPatch:
		h.updateBanner(w, req)
	case http.MethodDelete:
		h.deleteBanner(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandlerImpl) TestCleanUp(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		if err := h.RedisClient.TestCleanUp(context.Background()); err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
		if err := h.MongoClient.TestCleanUp(context.Background()); err != nil {
			JSONError(w, err, http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandlerImpl) getBanner(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	featureID, err := strconv.ParseInt(q.Get("feature_id"), 10, 64)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	tagID, err := strconv.ParseInt(q.Get("tag_id"), 10, 64)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	useLastRevision := false
	if q.Has("use_last_revision") {
		useLastRevision, err = strconv.ParseBool(q.Get("use_last_revision"))
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
	}

	var banner *types.Banner
	if !useLastRevision {
		banner, err = h.RedisClient.GetOne(context.Background(), featureID, tagID)
		switch {
		case errors.Is(err, redis.ErrBannerNotFound):
			break
		case err != nil:
			JSONError(w, err, http.StatusInternalServerError)
			return
		default:
			JSONBanner(w, banner)
			return
		}
	}
	banner, err = h.MongoClient.GetOne(context.Background(), featureID, tagID)
	if errors.Is(err, mongo.ErrBannerNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	err = h.RedisClient.Update(context.Background(), banner)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	JSONBanner(w, banner)
}

func (h *HTTPHandlerImpl) getBanners(w http.ResponseWriter, req *http.Request) {
	var filter types.BannerFilter
	var featureID, tagID, limit, offset int64
	var err error

	q := req.URL.Query()
	if q.Has("feature_id") {
		featureID, err = strconv.ParseInt(q.Get("feature_id"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		filter.FeatureID = &featureID
	}
	if q.Has("tag_id") {
		tagID, err = strconv.ParseInt(q.Get("tag_id"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		filter.TagID = &tagID
	}
	if q.Has("limit") {
		limit, err = strconv.ParseInt(q.Get("limit"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
	}
	if q.Has("offset") {
		offset, err = strconv.ParseInt(q.Get("offset"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
	}

	banners, err := h.MongoClient.GetMany(context.Background(), &filter, limit, offset)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	JSONBanners(w, banners)
}

func (h *HTTPHandlerImpl) createBanner(w http.ResponseWriter, req *http.Request) {
	body := make([]byte, req.ContentLength)
	read, err := req.Body.Read(body)
	defer req.Body.Close()

	if err != io.EOF {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	if read != int(req.ContentLength) {
		JSONError(w, ErrInvalidBodyLength, http.StatusBadRequest)
		return
	}

	banner := &types.Banner{}
	err = json.Unmarshal(body, banner)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	bannerID, err := h.MongoClient.Create(context.Background(), banner)
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	JSONBannerID(w, bannerID)
}

func (h *HTTPHandlerImpl) deleteBanners(w http.ResponseWriter, req *http.Request) {
	var filter types.BannerFilter
	var featureID, tagID int64
	var err error

	q := req.URL.Query()
	if q.Has("feature_id") {
		featureID, err = strconv.ParseInt(q.Get("feature_id"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		filter.FeatureID = &featureID
	}
	if q.Has("tag_id") {
		tagID, err = strconv.ParseInt(q.Get("tag_id"), 10, 64)
		if err != nil {
			JSONError(w, err, http.StatusBadRequest)
			return
		}
		filter.TagID = &tagID
	}

	go h.MongoClient.DeleteMany(context.Background(), &filter)
}

func (h *HTTPHandlerImpl) updateBanner(w http.ResponseWriter, req *http.Request) {
	BannerID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}

	body := make([]byte, req.ContentLength)
	read, err := req.Body.Read(body)
	defer req.Body.Close()

	if err != io.EOF {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	if read != int(req.ContentLength) {
		JSONError(w, ErrInvalidBodyLength, http.StatusBadRequest)
		return
	}

	banner := &types.Banner{}
	err = json.Unmarshal(body, banner)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	banner.BannerID = BannerID

	err = h.MongoClient.Update(context.Background(), banner)
	if errors.Is(err, mongo.ErrBannerNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *HTTPHandlerImpl) deleteBanner(w http.ResponseWriter, req *http.Request) {
	BannerID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		JSONError(w, err, http.StatusBadRequest)
		return
	}
	err = h.MongoClient.DeleteOne(context.Background(), BannerID)
	if errors.Is(err, mongo.ErrBannerNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		JSONError(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
