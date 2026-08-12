package handler

import (
	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/common/response"
	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/user-service/internal/model"
	"github.com/fastsprint/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Login(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, userInfo)
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	user, err := h.userService.UpdateUser(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err := h.userService.ChangePassword(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := utils.ParseUint(idStr)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var query model.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	users, total, err := h.userService.ListUsers(query)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	userInfos := make([]model.UserInfo, 0, len(users))
	for _, u := range users {
		userInfos = append(userInfos, model.UserInfo{
			ID:          u.ID,
			Username:    u.Username,
			Email:       u.Email,
			DisplayName: u.DisplayName,
			AvatarURL:   u.AvatarURL,
			Status:      u.Status,
		})
	}

	response.Page(c, userInfos, total, query.Page, query.PageSize)
}
