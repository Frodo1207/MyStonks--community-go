package v1

import (
	"MyStonks-go/internal/service"
	"MyStonks-go/internal/service/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body schema.UserCreateReq true "用户信息"
// @Success 200 {object} schema.UserResp
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var req schema.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := service.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetUser 获取用户
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} schema.UserResp
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
