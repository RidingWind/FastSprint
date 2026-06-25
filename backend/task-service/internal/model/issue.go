package model

import (
	"time"

	"github.com/fastsprint/common/models"
)

type Issue struct {
	models.BaseModel
	ProjectID    uint       `gorm:"not null;index" json:"project_id"`
	IssueKey     string     `gorm:"size:30;uniqueIndex;not null" json:"issue_key"`
	Summary      string     `gorm:"size:255;not null" json:"summary"`
	Description  string     `gorm:"type:text" json:"description"`
	IssueTypeID  uint       `gorm:"not null" json:"issue_type_id"`
	StatusID     uint       `gorm:"not null;index" json:"status_id"`
	PriorityID   *uint      `json:"priority_id"`
	ReporterID   uint       `gorm:"not null;index" json:"reporter_id"`
	AssigneeID   *uint      `gorm:"index" json:"assignee_id"`
	ParentID     *uint      `gorm:"index" json:"parent_id"`
	SprintID     *uint      `gorm:"index" json:"sprint_id"`
	StoryPoints  *float64   `gorm:"type:decimal(5,1)" json:"story_points"`
	DueDate      *time.Time `json:"due_date"`
}

func (Issue) TableName() string {
	return "task_service.issues"
}

type IssueComment struct {
	models.BaseModel
	IssueID  uint   `gorm:"not null;index" json:"issue_id"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
	Content  string `gorm:"type:text;not null" json:"content"`
}

func (IssueComment) TableName() string {
	return "task_service.issue_comments"
}

type IssueAttachment struct {
	models.BaseModel
	IssueID    uint   `gorm:"not null;index" json:"issue_id"`
	UploaderID uint   `gorm:"not null" json:"uploader_id"`
	Filename   string `gorm:"size:255;not null" json:"filename"`
	FileSize   int64  `gorm:"not null" json:"file_size"`
	MimeType   string `gorm:"size:100" json:"mime_type"`
	FileURL    string `gorm:"size:500;not null" json:"file_url"`
}

func (IssueAttachment) TableName() string {
	return "task_service.issue_attachments"
}

type IssueChangelog struct {
	models.BaseModel
	IssueID  uint   `gorm:"not null;index" json:"issue_id"`
	AuthorID uint   `gorm:"not null" json:"author_id"`
	Field    string `gorm:"size:50;not null" json:"field"`
	OldValue string `gorm:"type:text" json:"old_value"`
	NewValue string `gorm:"type:text" json:"new_value"`
}

func (IssueChangelog) TableName() string {
	return "task_service.issue_changelogs"
}

type IssueWatcher struct {
	models.BaseModel
	IssueID uint `gorm:"not null;uniqueIndex:idx_issue_watcher" json:"issue_id"`
	UserID  uint `gorm:"not null;uniqueIndex:idx_issue_watcher" json:"user_id"`
}

func (IssueWatcher) TableName() string {
	return "task_service.issue_watchers"
}

type Sprint struct {
	models.BaseModel
	ProjectID   uint       `gorm:"not null;index" json:"project_id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"size:20;default:planning;index" json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

func (Sprint) TableName() string {
	return "task_service.sprints"
}

type CreateIssueRequest struct {
	ProjectID   uint       `json:"project_id" binding:"required"`
	Summary     string     `json:"summary" binding:"required,max=255"`
	Description string     `json:"description"`
	IssueTypeID uint       `json:"issue_type_id" binding:"required"`
	PriorityID  *uint      `json:"priority_id"`
	AssigneeID  *uint      `json:"assignee_id"`
	ParentID    *uint      `json:"parent_id"`
	SprintID    *uint      `json:"sprint_id"`
	StoryPoints *float64   `json:"story_points"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateIssueRequest struct {
	Summary     string     `json:"summary"`
	Description string     `json:"description"`
	IssueTypeID *uint      `json:"issue_type_id"`
	StatusID    *uint      `json:"status_id"`
	PriorityID  *uint      `json:"priority_id"`
	AssigneeID  *uint      `json:"assignee_id"`
	SprintID    *uint      `json:"sprint_id"`
	StoryPoints *float64   `json:"story_points"`
	DueDate     *time.Time `json:"due_date"`
}

type IssueQuery struct {
	models.PaginationQuery
	ProjectID  uint    `form:"project_id"`
	IssueTypeID *uint   `form:"issue_type_id"`
	StatusID   *uint   `form:"status_id"`
	PriorityID *uint   `form:"priority_id"`
	AssigneeID *uint   `form:"assignee_id"`
	ReporterID *uint   `form:"reporter_id"`
	SprintID   *uint   `form:"sprint_id"`
	Keyword    string  `form:"keyword"`
	SortBy     string  `form:"sort_by,default=created_at"`
	Order      string  `form:"order,default=desc"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateSprintRequest struct {
	ProjectID   uint       `json:"project_id" binding:"required"`
	Name        string     `json:"name" binding:"required,max=100"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UpdateSprintRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type SprintQuery struct {
	models.PaginationQuery
	ProjectID uint   `form:"project_id"`
	Status    string `form:"status"`
}

type BoardColumn struct {
	StatusID uint    `json:"status_id"`
	Name     string  `json:"name"`
	Issues   []Issue `json:"issues"`
}

type BoardData struct {
	Columns []BoardColumn `json:"columns"`
}
