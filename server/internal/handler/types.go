package handler

import "net/http"

type HTTPHandler interface {
	UserBanner(w http.ResponseWriter, req *http.Request)

	AdminBanner(w http.ResponseWriter, req *http.Request)

	AdminBannerWithID(w http.ResponseWriter, req *http.Request)

	// test environment only, system state clean up
	TestCleanUp(w http.ResponseWriter, req *http.Request)
}
