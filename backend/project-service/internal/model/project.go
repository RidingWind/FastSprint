package model

import "github.com/fastsprint/common/models"

type Project struct {
	models.BaseModel
	Key              string `gorm:"size:20;uniqueIndex;not null" json:"key"`
	Name             string `gorm:"size:100;not null" json:"name"`
	Description      string `gorm:"type:text" json:"description"`
	OwnerID          uint   `gorm:"not null" json:"owner_id"`
	Status           string `gorm:"size:20;default:active" json:"status"`
	DefaultAssigneeID uint  `json:"default_assignee_id"`
}

func (Project) TableName() string {
	return "project_service.projects"
}

type ProjectMember struct {
	models.BaseModel
	ProjectID uint   `gorm:"not null;uniqueIndex:idx_project_user" json:"project_id"`
	UserID    uint   `gorm:"not null;uniqueIndex:idx_project_user" json:"user_id"`
	Role      string `gorm:"size:30;default:member" json:"role"`
}

func (ProjectMember) TableName() string {
	return "project_service.project_members"
}

type Status struct {
	models.BaseModel
	ProjectID uint   `gorm:"not null;uniqueIndex:idx_project_status" json:"project_id"`
	Name      string `gorm:"size:50;not null;uniqueIndex:idx_project_status" json:"name"`
	Category  string `gorm:"size:30;not null" json:"category"`
	Color     string `gorm:"size:7;default:#666666" json:"color"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	IsInitial bool   `gorm:"default:false" json:"is_initial"`
}

func (Status) TableName() string {
	return "project_service.statuses"
}

type IssueType struct {
	models.BaseModel
	ProjectID  uint   `gorm:"not null;uniqueIndex:idx_project_issue_type" json:"project_id"`
	Name       string `gorm:"size:50;not null;uniqueIndex:idx_project_issue_type" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Icon       string `gorm:"size:50" json:"icon"`
	Color      string `gorm:"size:7;default:#666666" json:"color"`
	IsSubtask  bool   `gorm:"default:false" json:"is_subtask"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
}

func (IssueType) TableName() string {
	return "project_service.issue_types"
}

type Priority struct {
	models.BaseModel
	ProjectID   uint   `gorm:"not null;uniqueIndex:idx_project_priority" json:"project_id"`
	Name        string `gorm:"size:50;not null;uniqueIndex:idx_project_priority" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"size:50" json:"icon"`
	Color       string `gorm:"size:7;default:#666666" json:"color"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
}

func (Priority) TableName() string {
	return "project_service.priorities"
}

type CreateProjectRequest struct {
	Key              string `json:"key" binding:"required,min=2,max=20,uppercase"`
	Name             string `json:"name" binding:"required,min=1,max=100"`
	Description      string `json:"description"`
	DefaultAssigneeID uint  `json:"default_assignee_id"`
}

type UpdateProjectRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	DefaultAssigneeID uint  `json:"default_assignee_id"`
}

type ProjectQuery struct {
	models.PaginationQuery
	Keyword string `form:"keyword"`
	Status  string `form:"status"`
	UserID  uint   `form:"user_id"`
}

type AddMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member viewer"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member viewer"`
}

type CreateStatusRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Category  string `json:"category" binding:"required,oneof=todo in_progress done"`
	Color     string `json:"color"`
	SortOrder int    `json:"sort_order"`
	IsInitial bool   `json:"is_initial"`
}

type CreateIssueTypeRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsSubtask   bool   `json:"is_subtask"`
	SortOrder   int    `json:"sort_order"`
}

type CreatePriorityRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sort_order"`
}

type ProjectDetail struct {
	Project
	OwnerName        string `json:"owner_name"`
	DefaultAssigneeName string `json:"default_assignee_name,omitempty"`
}

type MemberInfo struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
	JoinedAt    string `json:"joined_at"`
}
