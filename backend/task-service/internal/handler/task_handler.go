package handler

import (
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/common/response"
	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/task-service/internal/model"
	"github.com/fastsprint/task-service/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateIssue(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	issue, err := h.taskService.CreateIssue(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, issue)
}

func (h *TaskHandler) GetIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	issue, err := h.taskService.GetIssueByID(id)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	response.Success(c, issue)
}

func (h *TaskHandler) GetIssueByKey(c *gin.Context) {
	key := c.Param("key")

	issue, err := h.taskService.GetIssueByKey(key)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	response.Success(c, issue)
}

func (h *TaskHandler) UpdateIssue(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	var req model.UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	issue, err := h.taskService.UpdateIssue(id, userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, issue)
}

func (h *TaskHandler) DeleteIssue(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	err = h.taskService.DeleteIssue(id, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) ListIssues(c *gin.Context) {
	var query model.IssueQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	issues, total, err := h.taskService.ListIssues(query)
	if err != nil {
		response.InternalError(c, "获取任务列表失败")
		return
	}

	response.Page(c, issues, total, query.Page, query.PageSize)
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	comment, err := h.taskService.AddComment(id, userID, req.Content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, comment)
}

func (h *TaskHandler) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	comments, err := h.taskService.GetComments(id)
	if err != nil {
		response.InternalError(c, "获取评论失败")
		return
	}

	response.Success(c, comments)
}

func (h *TaskHandler) DeleteComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的评论ID")
		return
	}

	err = h.taskService.DeleteComment(id, userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) GetChangelogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	changelogs, err := h.taskService.GetChangelogs(id)
	if err != nil {
		response.InternalError(c, "获取变更记录失败")
		return
	}

	response.Success(c, changelogs)
}

func (h *TaskHandler) CreateSprint(c *gin.Context) {
	var req model.CreateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	sprint, err := h.taskService.CreateSprint(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, sprint)
}

func (h *TaskHandler) GetSprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的冲刺ID")
		return
	}

	sprint, err := h.taskService.GetSprintByID(id)
	if err != nil {
		response.NotFound(c, "冲刺不存在")
		return
	}

	response.Success(c, sprint)
}

func (h *TaskHandler) UpdateSprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的冲刺ID")
		return
	}

	var req model.UpdateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	sprint, err := h.taskService.UpdateSprint(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, sprint)
}

func (h *TaskHandler) DeleteSprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的冲刺ID")
		return
	}

	err = h.taskService.DeleteSprint(id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) ListSprints(c *gin.Context) {
	var query model.SprintQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	sprints, total, err := h.taskService.ListSprints(query)
	if err != nil {
		response.InternalError(c, "获取冲刺列表失败")
		return
	}

	response.Page(c, sprints, total, query.Page, query.PageSize)
}

func (h *TaskHandler) StartSprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的冲刺ID")
		return
	}

	err = h.taskService.StartSprint(id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) CompleteSprint(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的冲刺ID")
		return
	}

	err = h.taskService.CompleteSprint(id)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TaskHandler) GetBoardData(c *gin.Context) {
	projectIDStr := c.Query("project_id")
	projectID, err := utils.ParseUint(projectIDStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	statusIDs := []uint{1, 2, 3, 4, 5}

	boardData, err := h.taskService.GetBoardData(projectID, statusIDs)
	if err != nil {
		response.InternalError(c, "获取看板数据失败")
		return
	}

	response.Success(c, boardData)
}
