package model

type VerifyMagicLinkResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyMagicLinkResponse = Response[VerifyMagicLinkResponseBody]
