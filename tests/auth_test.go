package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vvauijij/avito-tech-assignment/tests/token"
)

type authTestCase struct {
	URL        string
	HTTPMethod string
	Token      string

	ExpectedHTTPStatus int
}

var authTestCases = []authTestCase{
	{UserBannerEndpoint(0, 0, false), http.MethodGet, token.Generator.NewInvalidToken(), http.StatusUnauthorized},

	{AdminBannerEndpoint, http.MethodGet, token.Generator.NewInvalidToken(), http.StatusUnauthorized},
	{AdminBannerEndpoint, http.MethodGet, token.Generator.NewUserToken(), http.StatusForbidden},

	{AdminBannerEndpoint, http.MethodPost, token.Generator.NewInvalidToken(), http.StatusUnauthorized},
	{AdminBannerEndpoint, http.MethodPost, token.Generator.NewUserToken(), http.StatusForbidden},

	{AdminBannerEndpoint, http.MethodDelete, token.Generator.NewInvalidToken(), http.StatusUnauthorized},
	{AdminBannerEndpoint, http.MethodDelete, token.Generator.NewUserToken(), http.StatusForbidden},

	{AdminBannerWithIDEndpoint(0), http.MethodPatch, token.Generator.NewInvalidToken(), http.StatusUnauthorized},
	{AdminBannerWithIDEndpoint(0), http.MethodPatch, token.Generator.NewUserToken(), http.StatusForbidden},

	{AdminBannerWithIDEndpoint(0), http.MethodDelete, token.Generator.NewInvalidToken(), http.StatusUnauthorized},
	{AdminBannerWithIDEndpoint(0), http.MethodDelete, token.Generator.NewUserToken(), http.StatusForbidden},
}

func TestAuth(t *testing.T) {
	CleanUp(t)

	client := &http.Client{}
	for _, tc := range authTestCases {
		t.Run(tc.URL, func(t *testing.T) {
			req, err := http.NewRequest(tc.HTTPMethod, tc.URL, nil)
			assert.NoError(t, err)
			req.Header.Add("token", tc.Token)

			resp, err := client.Do(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedHTTPStatus, resp.StatusCode)
		})
	}
}
