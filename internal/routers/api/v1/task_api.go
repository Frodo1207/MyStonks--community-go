package v1

import (
	"MyStonks-go/internal/common/response"
	"MyStonks-go/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskApi struct {
	taskService *service.TaskService
}

func NewTaskApi(taskService *service.TaskService) *TaskApi {
	return &TaskApi{
		taskService: taskService,
	}
}

func (t *TaskApi) GetTasksByCategory(c *gin.Context) {
	category := c.Query("category")
	addr := c.Query("addr")

	if category == "" {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "缺少分类参数",
		})
		return
	}

	// 获取该分类下的所有任务
	tasks, err := t.taskService.GetTasksByCategory(category, addr)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取任务列表失败",
		})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, response.SuccessResponse(tasks))
}

func (t *TaskApi) GetUserInfoTask(c *gin.Context) {
	addr := c.Query("addr")

	userInfoTask, err := t.taskService.GetUserInfoTask(addr)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse(userInfoTask))
}

func (t *TaskApi) FinishTask(c *gin.Context) {
	addr := c.Query("addr")
	taskIDStr := c.Query("task_id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "任务ID必须是数字",
		})
		return
	}
	err = t.taskService.CompleteTask(addr, taskID)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "完成任务失败",
		})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse(gin.H{"message": "完成任务成功"}))
}
func (t *TaskApi) CheckStonksTrade(c *gin.Context) {
	address := c.Query("sol_address")
	isTradeOK, err := t.taskService.CheckStonksTrade(address)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":     500,
			"message":  "获取交易信息失败",
			"is_trade": false,
		})
		return
	}
	if !isTradeOK {
		c.JSON(200, gin.H{
			"code":     200,
			"message":  "no trade",
			"is_trade": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"code":     200,
		"message":  "success trade",
		"is_trade": true,
	})
}

func (t *TaskApi) RefreshDailyTask(c *gin.Context) {
	err := t.taskService.RefreshDailyTask()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "刷新任务失败",
		})
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("刷新任务成功"))
}

func (t *TaskApi) GetRankBoard(c *gin.Context) {
	rankBoard, err := t.taskService.GetRankBoard()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取排行榜失败",
		})
		return
	}
	c.JSON(http.StatusOK, rankBoard)
}

func (t *TaskApi) UpdateTaskProgress(c *gin.Context) {}

func (t *TaskApi) GetUserRank(c *gin.Context) {}

func (t *TaskApi) GetLeaderboard(c *gin.Context) {}

func (t *TaskApi) CheckCompleteTask(c *gin.Context) {}
