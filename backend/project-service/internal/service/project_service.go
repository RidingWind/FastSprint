package service

import (
	"errors"

	"github.com/fastsprint/project-service/internal/model"
	"github.com/fastsprint/project-service/internal/repository"
)

type ProjectService interface {
	CreateProject(ownerID uint, req model.CreateProjectRequest) (*model.Project, error)
	GetProjectByID(id uint) (*model.Project, error)
	GetProjectByKey(key string) (*model.Project, error)
	UpdateProject(id uint, req model.UpdateProjectRequest) (*model.Project, error)
	DeleteProject(id uint) error
	ListProjects(query model.ProjectQuery) ([]model.Project, int64, error)
	IsProjectMember(projectID, userID uint) (bool, error)
	HasProjectPermission(projectID, userID uint, requiredRole string) (bool, error)

	AddMember(projectID uint, req model.AddMemberRequest) error
	RemoveMember(projectID, userID uint) error
	UpdateMemberRole(projectID, userID uint, role string) error
	ListMembers(projectID uint) ([]model.ProjectMember, error)

	CreateStatus(projectID uint, req model.CreateStatusRequest) (*model.Status, error)
	GetStatusByID(id uint) (*model.Status, error)
	ListStatuses(projectID uint) ([]model.Status, error)
	UpdateStatus(id uint, req model.CreateStatusRequest) (*model.Status, error)
	DeleteStatus(id uint) error

	CreateIssueType(projectID uint, req model.CreateIssueTypeRequest) (*model.IssueType, error)
	GetIssueTypeByID(id uint) (*model.IssueType, error)
	ListIssueTypes(projectID uint) ([]model.IssueType, error)
	UpdateIssueType(id uint, req model.CreateIssueTypeRequest) (*model.IssueType, error)
	DeleteIssueType(id uint) error

	CreatePriority(projectID uint, req model.CreatePriorityRequest) (*model.Priority, error)
	GetPriorityByID(id uint) (*model.Priority, error)
	ListPriorities(projectID uint) ([]model.Priority, error)
	UpdatePriority(id uint, req model.CreatePriorityRequest) (*model.Priority, error)
	DeletePriority(id uint) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{projectRepo: projectRepo}
}

var roleHierarchy = map[string]int{
	"viewer": 1,
	"member": 2,
	"admin":  3,
	"owner":  4,
}

func (s *projectService) CreateProject(ownerID uint, req model.CreateProjectRequest) (*model.Project, error) {
	existing, _ := s.projectRepo.GetByKey(req.Key)
	if existing != nil {
		return nil, errors.New("项目键已存在")
	}

	project := &model.Project{
		Key:               req.Key,
		Name:              req.Name,
		Description:       req.Description,
		OwnerID:           ownerID,
		Status:            "active",
		DefaultAssigneeID: req.DefaultAssigneeID,
	}

	err := s.projectRepo.Create(project)
	if err != nil {
		return nil, err
	}

	s.projectRepo.AddMember(&model.ProjectMember{
		ProjectID: project.ID,
		UserID:    ownerID,
		Role:      "admin",
	})

	s.initDefaultStatuses(project.ID)
	s.initDefaultIssueTypes(project.ID)
	s.initDefaultPriorities(project.ID)

	return project, nil
}

func (s *projectService) initDefaultStatuses(projectID uint) {
	statuses := []model.Status{
		{ProjectID: projectID, Name: "待办", Category: "todo", Color: "#94a3b8", SortOrder: 1, IsInitial: true},
		{ProjectID: projectID, Name: "进行中", Category: "in_progress", Color: "#3b82f6", SortOrder: 2},
		{ProjectID: projectID, Name: "评审中", Category: "in_progress", Color: "#8b5cf6", SortOrder: 3},
		{ProjectID: projectID, Name: "已完成", Category: "done", Color: "#22c55e", SortOrder: 4},
		{ProjectID: projectID, Name: "已关闭", Category: "done", Color: "#64748b", SortOrder: 5},
	}
	for _, status := range statuses {
		s.projectRepo.CreateStatus(&status)
	}
}

func (s *projectService) initDefaultIssueTypes(projectID uint) {
	issueTypes := []model.IssueType{
		{ProjectID: projectID, Name: "任务", Description: "普通任务", Icon: "task", Color: "#4ade80", IsSubtask: false, SortOrder: 1},
		{ProjectID: projectID, Name: "故事", Description: "用户故事", Icon: "story", Color: "#60a5fa", IsSubtask: false, SortOrder: 2},
		{ProjectID: projectID, Name: "缺陷", Description: "Bug缺陷", Icon: "bug", Color: "#f87171", IsSubtask: false, SortOrder: 3},
		{ProjectID: projectID, Name: "子任务", Description: "子任务", Icon: "subtask", Color: "#94a3b8", IsSubtask: true, SortOrder: 4},
	}
	for _, it := range issueTypes {
		s.projectRepo.CreateIssueType(&it)
	}
}

func (s *projectService) initDefaultPriorities(projectID uint) {
	priorities := []model.Priority{
		{ProjectID: projectID, Name: "最高", Description: "阻塞性问题", Icon: "highest", Color: "#ef4444", SortOrder: 1},
		{ProjectID: projectID, Name: "高", Description: "重要问题", Icon: "high", Color: "#f97316", SortOrder: 2},
		{ProjectID: projectID, Name: "中", Description: "普通问题", Icon: "medium", Color: "#eab308", SortOrder: 3},
		{ProjectID: projectID, Name: "低", Description: "次要问题", Icon: "low", Color: "#22c55e", SortOrder: 4},
		{ProjectID: projectID, Name: "最低", Description: "可以延后", Icon: "lowest", Color: "#64748b", SortOrder: 5},
	}
	for _, p := range priorities {
		s.projectRepo.CreatePriority(&p)
	}
}

func (s *projectService) GetProjectByID(id uint) (*model.Project, error) {
	return s.projectRepo.GetByID(id)
}

func (s *projectService) GetProjectByKey(key string) (*model.Project, error) {
	return s.projectRepo.GetByKey(key)
}

func (s *projectService) UpdateProject(id uint, req model.UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}
	if req.DefaultAssigneeID > 0 {
		project.DefaultAssigneeID = req.DefaultAssigneeID
	}

	err = s.projectRepo.Update(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) DeleteProject(id uint) error {
	_, err := s.projectRepo.GetByID(id)
	if err != nil {
		return errors.New("项目不存在")
	}
	return s.projectRepo.Delete(id)
}

func (s *projectService) ListProjects(query model.ProjectQuery) ([]model.Project, int64, error) {
	return s.projectRepo.List(query)
}

func (s *projectService) IsProjectMember(projectID, userID uint) (bool, error) {
	return s.projectRepo.IsProjectMember(projectID, userID)
}

func (s *projectService) HasProjectPermission(projectID, userID uint, requiredRole string) (bool, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return false, errors.New("项目不存在")
	}

	if project.OwnerID == userID {
		return true, nil
	}

	role, err := s.projectRepo.GetMemberRole(projectID, userID)
	if err != nil {
		return false, nil
	}

	userLevel := roleHierarchy[role]
	requiredLevel := roleHierarchy[requiredRole]

	return userLevel >= requiredLevel, nil
}

func (s *projectService) AddMember(projectID uint, req model.AddMemberRequest) error {
	member := &model.ProjectMember{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}
	return s.projectRepo.AddMember(member)
}

func (s *projectService) RemoveMember(projectID, userID uint) error {
	return s.projectRepo.RemoveMember(projectID, userID)
}

func (s *projectService) UpdateMemberRole(projectID, userID uint, role string) error {
	return s.projectRepo.UpdateMemberRole(projectID, userID, role)
}

func (s *projectService) ListMembers(projectID uint) ([]model.ProjectMember, error) {
	return s.projectRepo.ListMembers(projectID)
}

func (s *projectService) CreateStatus(projectID uint, req model.CreateStatusRequest) (*model.Status, error) {
	status := &model.Status{
		ProjectID: projectID,
		Name:      req.Name,
		Category:  req.Category,
		Color:     req.Color,
		SortOrder: req.SortOrder,
		IsInitial: req.IsInitial,
	}
	if status.Color == "" {
		status.Color = "#666666"
	}
	err := s.projectRepo.CreateStatus(status)
	return status, err
}

func (s *projectService) GetStatusByID(id uint) (*model.Status, error) {
	return s.projectRepo.GetStatusByID(id)
}

func (s *projectService) ListStatuses(projectID uint) ([]model.Status, error) {
	return s.projectRepo.ListStatuses(projectID)
}

func (s *projectService) UpdateStatus(id uint, req model.CreateStatusRequest) (*model.Status, error) {
	status, err := s.projectRepo.GetStatusByID(id)
	if err != nil {
		return nil, errors.New("状态不存在")
	}

	if req.Name != "" {
		status.Name = req.Name
	}
	if req.Category != "" {
		status.Category = req.Category
	}
	if req.Color != "" {
		status.Color = req.Color
	}
	status.SortOrder = req.SortOrder
	status.IsInitial = req.IsInitial

	err = s.projectRepo.UpdateStatus(status)
	return status, err
}

func (s *projectService) DeleteStatus(id uint) error {
	_, err := s.projectRepo.GetStatusByID(id)
	if err != nil {
		return errors.New("状态不存在")
	}
	return s.projectRepo.DeleteStatus(id)
}

func (s *projectService) CreateIssueType(projectID uint, req model.CreateIssueTypeRequest) (*model.IssueType, error) {
	issueType := &model.IssueType{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		IsSubtask:   req.IsSubtask,
		SortOrder:   req.SortOrder,
	}
	if issueType.Color == "" {
		issueType.Color = "#666666"
	}
	err := s.projectRepo.CreateIssueType(issueType)
	return issueType, err
}

func (s *projectService) GetIssueTypeByID(id uint) (*model.IssueType, error) {
	return s.projectRepo.GetIssueTypeByID(id)
}

func (s *projectService) ListIssueTypes(projectID uint) ([]model.IssueType, error) {
	return s.projectRepo.ListIssueTypes(projectID)
}

func (s *projectService) UpdateIssueType(id uint, req model.CreateIssueTypeRequest) (*model.IssueType, error) {
	issueType, err := s.projectRepo.GetIssueTypeByID(id)
	if err != nil {
		return nil, errors.New("任务类型不存在")
	}

	if req.Name != "" {
		issueType.Name = req.Name
	}
	if req.Description != "" {
		issueType.Description = req.Description
	}
	if req.Icon != "" {
		issueType.Icon = req.Icon
	}
	if req.Color != "" {
		issueType.Color = req.Color
	}
	issueType.IsSubtask = req.IsSubtask
	issueType.SortOrder = req.SortOrder

	err = s.projectRepo.UpdateIssueType(issueType)
	return issueType, err
}

func (s *projectService) DeleteIssueType(id uint) error {
	_, err := s.projectRepo.GetIssueTypeByID(id)
	if err != nil {
		return errors.New("任务类型不存在")
	}
	return s.projectRepo.DeleteIssueType(id)
}

func (s *projectService) CreatePriority(projectID uint, req model.CreatePriorityRequest) (*model.Priority, error) {
	priority := &model.Priority{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		SortOrder:   req.SortOrder,
	}
	if priority.Color == "" {
		priority.Color = "#666666"
	}
	err := s.projectRepo.CreatePriority(priority)
	return priority, err
}

func (s *projectService) GetPriorityByID(id uint) (*model.Priority, error) {
	return s.projectRepo.GetPriorityByID(id)
}

func (s *projectService) ListPriorities(projectID uint) ([]model.Priority, error) {
	return s.projectRepo.ListPriorities(projectID)
}

func (s *projectService) UpdatePriority(id uint, req model.CreatePriorityRequest) (*model.Priority, error) {
	priority, err := s.projectRepo.GetPriorityByID(id)
	if err != nil {
		return nil, errors.New("优先级不存在")
	}

	if req.Name != "" {
		priority.Name = req.Name
	}
	if req.Description != "" {
		priority.Description = req.Description
	}
	if req.Icon != "" {
		priority.Icon = req.Icon
	}
	if req.Color != "" {
		priority.Color = req.Color
	}
	priority.SortOrder = req.SortOrder

	err = s.projectRepo.UpdatePriority(priority)
	return priority, err
}

func (s *projectService) DeletePriority(id uint) error {
	_, err := s.projectRepo.GetPriorityByID(id)
	if err != nil {
		return errors.New("优先级不存在")
	}
	return s.projectRepo.DeletePriority(id)
}
