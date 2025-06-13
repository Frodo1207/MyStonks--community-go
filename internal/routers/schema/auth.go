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

type TgInfo struct {
	FirstName  string `json:"first_name"`
	TelegramID int64  `json:"telegram_id" binding:"required"` // 或者改为 telegramId
	Username   string `json:"username"`
	PhotoURL   string `json:"photo_url"`
	AuthDate   int64  `json:"auth_date"`
	Hash       string `json:"hash" binding:"required"`
}

type UserInfo struct {
	Addr string `json:"addr"`
}
type LoginResp struct {
	Data   TokenPair `json:"data"`
	TgInfo TgInfo    `json:"tg_info"`
	User   UserInfo  `json:"user"`
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

type BindTgReq struct {
	Addr       string `json:"addr" binding:"required"`
	FirstName  string `json:"first_name"`
	TelegramID int64  `json:"telegram_id" binding:"required"` // 或者改为 telegramId
	Username   string `json:"username"`
	PhotoURL   string `json:"photo_url"`
	AuthDate   int64  `json:"auth_date"`
	Hash       string `json:"hash" binding:"required"`
}
