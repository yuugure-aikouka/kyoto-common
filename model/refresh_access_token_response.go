package model

type RefreshAccessTokenResponseBody struct {
	AccessToken string `json:"access_token"`
}

type RefreshAccessTokenResponse = Response[RefreshAccessTokenResponseBody]
