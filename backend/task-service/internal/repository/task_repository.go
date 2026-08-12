package repository

import (
	"fmt"

	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/task-service/internal/model"
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateIssue(issue *model.Issue) error
	GetIssueByID(id uint) (*model.Issue, error)
	GetIssueByKey(issueKey string) (*model.Issue, error)
	UpdateIssue(issue *model.Issue) error
	DeleteIssue(id uint) error
	ListIssues(query model.IssueQuery) ([]model.Issue, int64, error)
	GetNextIssueNumber(projectID uint) (int64, error)
	GetIssuesBySprint(sprintID uint) ([]model.Issue, error)
	GetIssuesByStatus(projectID uint, statusIDs []uint) ([]model.Issue, error)

	CreateComment(comment *model.IssueComment) error
	GetCommentsByIssue(issueID uint) ([]model.IssueComment, error)
	DeleteComment(id, authorID uint) error

	CreateChangelog(changelog *model.IssueChangelog) error
	GetChangelogsByIssue(issueID uint) ([]model.IssueChangelog, error)

	CreateSprint(sprint *model.Sprint) error
	GetSprintByID(id uint) (*model.Sprint, error)
	UpdateSprint(sprint *model.Sprint) error
	DeleteSprint(id uint) error
	ListSprints(query model.SprintQuery) ([]model.Sprint, int64, error)
	GetActiveSprint(projectID uint) (*model.Sprint, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateIssue(issue *model.Issue) error {
	return r.db.Create(issue).Error
}

func (r *taskRepository) GetIssueByID(id uint) (*model.Issue, error) {
	var issue model.Issue
	err := r.db.First(&issue, id).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *taskRepository) GetIssueByKey(issueKey string) (*model.Issue, error) {
	var issue model.Issue
	err := r.db.Where("issue_key = ?", issueKey).First(&issue).Error
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (r *taskRepository) UpdateIssue(issue *model.Issue) error {
	return r.db.Save(issue).Error
}

func (r *taskRepository) DeleteIssue(id uint) error {
	return r.db.Delete(&model.Issue{}, id).Error
}

func (r *taskRepository) ListIssues(query model.IssueQuery) ([]model.Issue, int64, error) {
	var issues []model.Issue
	var total int64

	db := r.db.Model(&model.Issue{})

	if query.ProjectID > 0 {
		db = db.Where("project_id = ?", query.ProjectID)
	}
	if query.IssueTypeID != nil {
		db = db.Where("issue_type_id = ?", *query.IssueTypeID)
	}
	if query.StatusID != nil {
		db = db.Where("status_id = ?", *query.StatusID)
	}
	if query.PriorityID != nil {
		db = db.Where("priority_id = ?", *query.PriorityID)
	}
	if query.AssigneeID != nil {
		db = db.Where("assignee_id = ?", *query.AssigneeID)
	}
	if query.ReporterID != nil {
		db = db.Where("reporter_id = ?", *query.ReporterID)
	}
	if query.SprintID != nil {
		db = db.Where("sprint_id = ?", *query.SprintID)
	}
	if query.Keyword != "" {
		db = db.Where("summary LIKE ? OR description LIKE ? OR issue_key LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.Paginate(query.Page, query.PageSize)
	orderStr := fmt.Sprintf("%s %s", query.SortBy, query.Order)
	err := db.Offset(offset).Limit(limit).Order(orderStr).Find(&issues).Error

	return issues, total, err
}

func (r *taskRepository) GetNextIssueNumber(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Issue{}).
		Where("project_id = ?", projectID).
		Count(&count).Error
	return count + 1, err
}

func (r *taskRepository) GetIssuesBySprint(sprintID uint) ([]model.Issue, error) {
	var issues []model.Issue
	err := r.db.Where("sprint_id = ?", sprintID).Order("id desc").Find(&issues).Error
	return issues, err
}

func (r *taskRepository) GetIssuesByStatus(projectID uint, statusIDs []uint) ([]model.Issue, error) {
	var issues []model.Issue
	db := r.db.Where("project_id = ?", projectID)
	if len(statusIDs) > 0 {
		db = db.Where("status_id IN ?", statusIDs)
	}
	err := db.Order("id desc").Find(&issues).Error
	return issues, err
}

func (r *taskRepository) CreateComment(comment *model.IssueComment) error {
	return r.db.Create(comment).Error
}

func (r *taskRepository) GetCommentsByIssue(issueID uint) ([]model.IssueComment, error) {
	var comments []model.IssueComment
	err := r.db.Where("issue_id = ?", issueID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (r *taskRepository) DeleteComment(id, authorID uint) error {
	return r.db.Where("id = ? AND author_id = ?", id, authorID).
		Delete(&model.IssueComment{}).Error
}

func (r *taskRepository) CreateChangelog(changelog *model.IssueChangelog) error {
	return r.db.Create(changelog).Error
}

func (r *taskRepository) GetChangelogsByIssue(issueID uint) ([]model.IssueChangelog, error) {
	var changelogs []model.IssueChangelog
	err := r.db.Where("issue_id = ?", issueID).Order("created_at desc").Find(&changelogs).Error
	return changelogs, err
}

func (r *taskRepository) CreateSprint(sprint *model.Sprint) error {
	return r.db.Create(sprint).Error
}

func (r *taskRepository) GetSprintByID(id uint) (*model.Sprint, error) {
	var sprint model.Sprint
	err := r.db.First(&sprint, id).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (r *taskRepository) UpdateSprint(sprint *model.Sprint) error {
	return r.db.Save(sprint).Error
}

func (r *taskRepository) DeleteSprint(id uint) error {
	return r.db.Delete(&model.Sprint{}, id).Error
}

func (r *taskRepository) ListSprints(query model.SprintQuery) ([]model.Sprint, int64, error) {
	var sprints []model.Sprint
	var total int64

	db := r.db.Model(&model.Sprint{})

	if query.ProjectID > 0 {
		db = db.Where("project_id = ?", query.ProjectID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.Paginate(query.Page, query.PageSize)
	err := db.Offset(offset).Limit(limit).Order("id desc").Find(&sprints).Error

	return sprints, total, err
}

func (r *taskRepository) GetActiveSprint(projectID uint) (*model.Sprint, error) {
	var sprint model.Sprint
	err := r.db.Where("project_id = ? AND status = ?", projectID, "active").
		Order("id desc").First(&sprint).Error
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}
