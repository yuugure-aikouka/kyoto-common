package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuugure-aikouka/kyoto-common/model"
	helper "github.com/yuugure-aikouka/kyoto-common/test/helper"
)

func TestSendMagicLink(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		setupTest()

		testEmail := "anondev71@gmail.com"
		sendMagicLink(t, testEmail)
	})

	t.Run("Missing email field", func(t *testing.T) {
		setupTest()

		url := fmt.Sprintf("/v1/auth/magic")
		req := httptest.NewRequest(http.MethodPost, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Invalid email", func(t *testing.T) {
		setupTest()

		testEmail := "invalidemail"
		body := model.SendMagicLinkRequest{Email: testEmail}

		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err)

		url := fmt.Sprintf("/v1/auth/magic")
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestVerifyMagicLink(t *testing.T) {
	testEmail := "anondev71@gmail.com"
	t.Run("Success", func(t *testing.T) {
		setupTest()

		sendMagicLink(t, testEmail)
		sentEmail := mailer.SentEmails[0]
		token := strings.Split(sentEmail.Link, "token=")[1]

		url := fmt.Sprintf("/v1/auth/magic?token=%s", token)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		resBody, err := helper.UnmarshalResponseBody[model.VerifyMagicLinkResponseBody](res.Body)
		require.Nil(t, err)

		// validate
		require.Equal(t, http.StatusOK, res.Code)
		require.NotNil(t, resBody.Data.AccessToken)
		require.NotNil(t, resBody.Data.RefreshToken)

		keys, err := helper.GetMagicLinkKeysByEmail(rdbConn, testEmail)
		require.NoError(t, err)
		require.Empty(t, keys)
	})

	t.Run("Missing token", func(t *testing.T) {
		setupTest()

		url := "/v1/auth/magic"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		setupTest()

		url := "/v1/auth/magic?token=invalid-token"
		req := httptest.NewRequest(http.MethodGet, url, nil)
		res := httptest.NewRecorder()
		server.Route().ServeHTTP(res, req)

		require.Equal(t, http.StatusUnauthorized, res.Code)
	})
}

func sendMagicLink(t *testing.T, email string) {
	body := model.SendMagicLinkRequest{Email: email}

	jsonBytes, err := json.Marshal(body)
	require.NoError(t, err)

	url := "/v1/auth/magic"
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	server.Route().ServeHTTP(res, req)

	// validate
	require.Equal(t, http.StatusOK, res.Code)
	require.Len(t, mailer.SentEmails, 1)
	require.Equal(t, email, mailer.SentEmails[0].To)

	keys, err := helper.GetMagicLinkKeysByEmail(rdbConn, email)
	require.NoError(t, err)
	require.NotEmpty(t, keys)
}
