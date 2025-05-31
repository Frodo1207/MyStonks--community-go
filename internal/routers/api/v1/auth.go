package v1

import (
	"MyStonks-go/internal/service"
	"MyStonks-go/internal/service/schema"
	"net/http"
	"strings"
	"time"

	"MyStonks-go/internal/common/response"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// GetNonce 获取登录随机数
// @Summary 获取登录随机数
// @Description 获取用于钱包签名的随机数
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} schema.NonceResp
// @Router /api/v1/auth/nonce [get]
func GetNonce(c *gin.Context) {
	nonce, err := service.GenerateNonce()
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(gin.H{"nonce": nonce}))
}

// Login 钱包登录
// @Summary 钱包登录
// @Description 使用钱包签名登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body schema.LoginReq true "登录请求参数"
// @Success 200 {object} schema.LoginResp
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req schema.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	valid, err := service.VerifySolanaWalletSignature(&req)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidSignature, []string{}))
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidSignature, []string{}))
		return
	}

	tokenPair, err := service.GenerateTokenPair(req.Address)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	if err := service.CreateUserIfNotExists(req.Address); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(tokenPair))
}

// Logout 退出登录
// @Summary 退出登录
// @Description 退出系统登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body schema.LogoutReq true "退出登录请求参数"
// @Success 200 {object} schema.LogoutResp
// @Router /api/v1/auth/logout [post]
func Logout(c *gin.Context) {
	var req schema.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if err := service.SetTokenBlacklist(token, time.Hour*24); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}
	if err := service.SetTokenBlacklist(req.RefreshToken, time.Hour*24*7); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
	}
	c.JSON(http.StatusOK, response.SuccessResponse(gin.H{}))
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用刷新令牌刷新访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body schema.RefreshTokenReq true "刷新令牌请求参数"
// @Success 200 {object} schema.RefreshTokenResp
// @Router /api/v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req schema.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	tokenPair, err := service.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(tokenPair))
}
