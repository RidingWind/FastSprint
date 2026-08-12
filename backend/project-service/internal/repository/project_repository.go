package repository

import (
	"errors"

	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/project-service/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *model.Project) error
	GetByID(id uint) (*model.Project, error)
	GetByKey(key string) (*model.Project, error)
	Update(project *model.Project) error
	Delete(id uint) error
	List(query model.ProjectQuery) ([]model.Project, int64, error)
	IsProjectMember(projectID, userID uint) (bool, error)
	GetMemberRole(projectID, userID uint) (string, error)

	AddMember(member *model.ProjectMember) error
	RemoveMember(projectID, userID uint) error
	UpdateMemberRole(projectID, userID uint, role string) error
	ListMembers(projectID uint) ([]model.ProjectMember, error)
	GetMember(projectID, userID uint) (*model.ProjectMember, error)

	CreateStatus(status *model.Status) error
	GetStatusByID(id uint) (*model.Status, error)
	ListStatuses(projectID uint) ([]model.Status, error)
	UpdateStatus(status *model.Status) error
	DeleteStatus(id uint) error

	CreateIssueType(issueType *model.IssueType) error
	GetIssueTypeByID(id uint) (*model.IssueType, error)
	ListIssueTypes(projectID uint) ([]model.IssueType, error)
	UpdateIssueType(issueType *model.IssueType) error
	DeleteIssueType(id uint) error

	CreatePriority(priority *model.Priority) error
	GetPriorityByID(id uint) (*model.Priority, error)
	ListPriorities(projectID uint) ([]model.Priority, error)
	UpdatePriority(priority *model.Priority) error
	DeletePriority(id uint) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetByID(id uint) (*model.Project, error) {
	var project model.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByKey(key string) (*model.Project, error) {
	var project model.Project
	err := r.db.Where("key = ?", key).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uint) error {
	return r.db.Delete(&model.Project{}, id).Error
}

func (r *projectRepository) List(query model.ProjectQuery) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	db := r.db.Model(&model.Project{})

	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR key LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.UserID > 0 {
		db = db.Joins("JOIN project_service.project_members ON project_members.project_id = projects.id").
			Where("project_members.user_id = ?", query.UserID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.Paginate(query.Page, query.PageSize)
	err := db.Offset(offset).Limit(limit).Order("id desc").Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) IsProjectMember(projectID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) GetMemberRole(projectID, userID uint) (string, error) {
	member, err := r.GetMember(projectID, userID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

func (r *projectRepository) AddMember(member *model.ProjectMember) error {
	existing, _ := r.GetMember(member.ProjectID, member.UserID)
	if existing != nil {
		return errors.New("用户已是项目成员")
	}
	return r.db.Create(member).Error
}

func (r *projectRepository) RemoveMember(projectID, userID uint) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&model.ProjectMember{}).Error
}

func (r *projectRepository) UpdateMemberRole(projectID, userID uint, role string) error {
	return r.db.Model(&model.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Update("role", role).Error
}

func (r *projectRepository) ListMembers(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := r.db.Where("project_id = ?", projectID).Order("id asc").Find(&members).Error
	return members, err
}

func (r *projectRepository) GetMember(projectID, userID uint) (*model.ProjectMember, error) {
	var member model.ProjectMember
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *projectRepository) CreateStatus(status *model.Status) error {
	return r.db.Create(status).Error
}

func (r *projectRepository) GetStatusByID(id uint) (*model.Status, error) {
	var status model.Status
	err := r.db.First(&status, id).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *projectRepository) ListStatuses(projectID uint) ([]model.Status, error) {
	var statuses []model.Status
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc, id asc").Find(&statuses).Error
	return statuses, err
}

func (r *projectRepository) UpdateStatus(status *model.Status) error {
	return r.db.Save(status).Error
}

func (r *projectRepository) DeleteStatus(id uint) error {
	return r.db.Delete(&model.Status{}, id).Error
}

func (r *projectRepository) CreateIssueType(issueType *model.IssueType) error {
	return r.db.Create(issueType).Error
}

func (r *projectRepository) GetIssueTypeByID(id uint) (*model.IssueType, error) {
	var issueType model.IssueType
	err := r.db.First(&issueType, id).Error
	if err != nil {
		return nil, err
	}
	return &issueType, nil
}

func (r *projectRepository) ListIssueTypes(projectID uint) ([]model.IssueType, error) {
	var issueTypes []model.IssueType
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc, id asc").Find(&issueTypes).Error
	return issueTypes, err
}

func (r *projectRepository) UpdateIssueType(issueType *model.IssueType) error {
	return r.db.Save(issueType).Error
}

func (r *projectRepository) DeleteIssueType(id uint) error {
	return r.db.Delete(&model.IssueType{}, id).Error
}

func (r *projectRepository) CreatePriority(priority *model.Priority) error {
	return r.db.Create(priority).Error
}

func (r *projectRepository) GetPriorityByID(id uint) (*model.Priority, error) {
	var priority model.Priority
	err := r.db.First(&priority, id).Error
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

func (r *projectRepository) ListPriorities(projectID uint) ([]model.Priority, error) {
	var priorities []model.Priority
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc, id asc").Find(&priorities).Error
	return priorities, err
}

func (r *projectRepository) UpdatePriority(priority *model.Priority) error {
	return r.db.Save(priority).Error
}

func (r *projectRepository) DeletePriority(id uint) error {
	return r.db.Delete(&model.Priority{}, id).Error
}
