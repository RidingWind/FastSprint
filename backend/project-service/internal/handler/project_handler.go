package handler

import (
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/common/response"
	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/project-service/internal/model"
	"github.com/fastsprint/project-service/internal/service"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	project, err := h.projectService.CreateProject(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, project)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	project, err := h.projectService.GetProjectByID(id)
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(id, userID)
	if !isMember && project.OwnerID != userID {
		response.Forbidden(c, "无权限访问该项目")
		return
	}

	response.Success(c, project)
}

func (h *ProjectHandler) GetProjectByKey(c *gin.Context) {
	userID := middleware.GetUserID(c)
	key := c.Param("key")

	project, err := h.projectService.GetProjectByKey(key)
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(project.ID, userID)
	if !isMember && project.OwnerID != userID {
		response.Forbidden(c, "无权限访问该项目")
		return
	}

	response.Success(c, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限修改项目")
		return
	}

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	project, err := h.projectService.UpdateProject(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	project, err := h.projectService.GetProjectByID(id)
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	if project.OwnerID != userID {
		response.Forbidden(c, "只有项目所有者可以删除项目")
		return
	}

	err = h.projectService.DeleteProject(id)
	if err != nil {
		response.InternalError(c, "删除项目失败")
		return
	}

	response.Success(c, nil)
}

func (h *ProjectHandler) ListProjects(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var query model.ProjectQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	query.UserID = userID

	projects, total, err := h.projectService.ListProjects(query)
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	response.Page(c, projects, total, query.Page, query.PageSize)
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限添加成员")
		return
	}

	var req model.AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err = h.projectService.AddMember(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	memberIDStr := c.Param("userId")
	memberID, err := utils.ParseUint(memberIDStr)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限移除成员")
		return
	}

	err = h.projectService.RemoveMember(id, memberID)
	if err != nil {
		response.InternalError(c, "移除成员失败")
		return
	}

	response.Success(c, nil)
}

func (h *ProjectHandler) UpdateMemberRole(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	memberIDStr := c.Param("userId")
	memberID, err := utils.ParseUint(memberIDStr)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限修改成员角色")
		return
	}

	var req model.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err = h.projectService.UpdateMemberRole(id, memberID, req.Role)
	if err != nil {
		response.InternalError(c, "修改角色失败")
		return
	}

	response.Success(c, nil)
}

func (h *ProjectHandler) ListMembers(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(id, userID)
	if !isMember {
		project, _ := h.projectService.GetProjectByID(id)
		if project == nil || project.OwnerID != userID {
			response.Forbidden(c, "无权限查看成员列表")
			return
		}
	}

	members, err := h.projectService.ListMembers(id)
	if err != nil {
		response.InternalError(c, "获取成员列表失败")
		return
	}

	response.Success(c, members)
}

func (h *ProjectHandler) ListStatuses(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(id, userID)
	if !isMember {
		response.Forbidden(c, "无权限访问")
		return
	}

	statuses, err := h.projectService.ListStatuses(id)
	if err != nil {
		response.InternalError(c, "获取状态列表失败")
		return
	}

	response.Success(c, statuses)
}

func (h *ProjectHandler) CreateStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限添加状态")
		return
	}

	var req model.CreateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	status, err := h.projectService.CreateStatus(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, status)
}

func (h *ProjectHandler) ListIssueTypes(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(id, userID)
	if !isMember {
		response.Forbidden(c, "无权限访问")
		return
	}

	issueTypes, err := h.projectService.ListIssueTypes(id)
	if err != nil {
		response.InternalError(c, "获取任务类型列表失败")
		return
	}

	response.Success(c, issueTypes)
}

func (h *ProjectHandler) CreateIssueType(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限添加任务类型")
		return
	}

	var req model.CreateIssueTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	issueType, err := h.projectService.CreateIssueType(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, issueType)
}

func (h *ProjectHandler) ListPriorities(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	isMember, _ := h.projectService.IsProjectMember(id, userID)
	if !isMember {
		response.Forbidden(c, "无权限访问")
		return
	}

	priorities, err := h.projectService.ListPriorities(id)
	if err != nil {
		response.InternalError(c, "获取优先级列表失败")
		return
	}

	response.Success(c, priorities)
}

func (h *ProjectHandler) CreatePriority(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	hasPerm, _ := h.projectService.HasProjectPermission(id, userID, "admin")
	if !hasPerm {
		response.Forbidden(c, "无权限添加优先级")
		return
	}

	var req model.CreatePriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	priority, err := h.projectService.CreatePriority(id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, priority)
}
