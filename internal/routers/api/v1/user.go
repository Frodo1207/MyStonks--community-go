package v1

import (
	"MyStonks-go/internal/common/response"
	"MyStonks-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary get user info
// @Description get user info
// @Tags user
// @Accept json
// @Produce json
// @Param sol_address query string true "sol address"
// @Success 200 {object} schema.UserResp
// @Router /api/v1/user [get]
func GetUserInfo(c *gin.Context) {
	solAddress := c.Query("sol_address")
	user, err := service.GetUserInfo(solAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ErrorCodeUserNotFound, []string{}))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(user))
}
