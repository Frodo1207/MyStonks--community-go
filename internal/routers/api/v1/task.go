package v1

import (
	"MyStonks-go/internal/common/response"
	"MyStonks-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary 获取任务列表
// @Description 根据任务分类获取任务列表
// @Tags 任务
// @Accept  json
// @Produce  json
// @Param category query string false "任务分类(daily/newbie/other)"
// @Success 200 {object} app.Response{data=[]models.Task}
// @Failure 500 {object} app.Response
// @Router /api/v1/tasks [get]
func GetTasksByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "缺少分类参数",
		})
		return
	}
	tasks, err := service.GetTasksByCategory(category)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取任务列失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(tasks))
}

// CheckCompleteTask CompleteTask @Summary 完成任务
// @Description 标记指定任务为已完成，并增加用户积分
// @Tags 任务
// @Accept  json
// @Produce  json
// @Param user_id query int true "用户ID"
// @Param task_id query int true "任务ID"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/task/complete [post]
func CheckCompleteTask(c *gin.Context) {
	userId := c.Query("user_id")
	taskid := c.Query("task_id")
	taskComplete, err := service.IsTaskComplete(userId, taskid)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "检查任务失败",
		})
		return
	}

	if taskComplete {
		c.JSON(200, map[string]interface{}{
			"code":    200,
			"message": "任务已完成",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "任务未完成，正在进行中...",
	})

}

// GetUserInfoTask @Summary 获取用户信息和任务状态
// @Description 获取指定钱包地址的用户信息和任务完成状态
// @Tags 用户
// @Accept  json
// @Produce  json
// @Param address query string true "钱包地址"
// @Success 200 {object} app.Response{data=models.UserInfo}
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user [get]
func GetUserInfoTask(c *gin.Context) {
	userid := c.Query("user_id")
	userInfoTask, err := service.GetUserInfoTask(userid)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(userInfoTask))
}

// @Summary 更新任务进度
// @Description 更新任务进度，若进度达到最大值则自动完成任务
// @Tags 任务
// @Accept  json
// @Produce  json
// @Param user_id query int true "用户ID"
// @Param task_id query int true "任务ID"
// @Param progress query int true "任务进度（百分比）"
// @Success 200 {object} app.Response
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/task/progress [post]
func UpdateTaskProgress(c *gin.Context) {
}

// @Summary 获取排行榜
// @Description 获取积分排行榜
// @Tags 排行榜
// @Accept  json
// @Produce  json
// @Param limit query int false "排行榜返回数量，默认10"
// @Success 200 {object} app.Response{data=[]models.LeaderboardEntry}
// @Failure 500 {object} app.Response
// @Router /api/v1/leaderboard [get]
func GetLeaderboard(c *gin.Context) {
}

// @Summary 获取用户排名
// @Description 获取指定用户的排行榜信息
// @Tags 排行榜
// @Accept  json
// @Produce  json
// @Param user_id query int true "用户ID"
// @Success 200 {object} app.Response{data=models.UserRank}
// @Failure 400 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/user/rank [get]
func GetUserRank(c *gin.Context) {
}
