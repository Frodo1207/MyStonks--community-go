package v1

import (
	"MyStonks-go/internal/common/response"
	"MyStonks-go/internal/routers/schema"
	"MyStonks-go/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

type UserApi struct {
	userSrv *service.UserSrv
}

func NewUserApi(uSrv *service.UserSrv) *UserApi {
	return &UserApi{
		userSrv: uSrv,
	}
}

func (u UserApi) GetUserInfo(c *gin.Context) {
	solAddress := c.Query("sol_address")
	user, err := u.userSrv.GetUserInfo(solAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeUserNotFound, []string{}))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(user))
}

func (u *UserApi) Login(c *gin.Context) {
	var req schema.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	valid, err := u.userSrv.VerifySolanaWalletSignature(&req)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidSignature, []string{}))
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidSignature, []string{}))
		return
	}

	tokenPair, err := u.userSrv.GenerateTokenPair(req.Address)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	if err := u.userSrv.CreateUserIfNotExists(req.Address); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(tokenPair))
}

func (u *UserApi) Logout(c *gin.Context) {
	var req schema.LogoutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if err := u.userSrv.SetTokenBlacklist(token, time.Hour*24); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}
	if err := u.userSrv.SetTokenBlacklist(req.RefreshToken, time.Hour*24*7); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
	}
	c.JSON(http.StatusOK, response.SuccessResponse(gin.H{}))
}
func (u *UserApi) RefreshToken(c *gin.Context) {
	var req schema.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusBadRequest, response.ErrorResponse(response.ErrorCodeInvalidRequest, []string{}))
		return
	}

	tokenPair, err := u.userSrv.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(tokenPair))
}

func (u *UserApi) GetNonce(c *gin.Context) {
	nonce, err := u.userSrv.GenerateNonce()
	if err != nil {
		log.Warn().Err(err).Send()
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeInternalError, []string{}))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(gin.H{"nonce": nonce}))
}
