package schema

import "MyStonks-go/internal/common/response"

type LoginReq struct {
	Address   string `json:"address"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResp struct {
	response.Response
	Data TokenPair `json:"data"`
}

type NonceResp struct {
	response.Response
	Data struct {
		Nonce string `json:"nonce"`
	} `json:"data"`
}

type LogoutResp struct {
	response.Response
	Data struct {
	} `json:"data"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResp struct {
	response.Response
	Data TokenPair `json:"data"`
}
