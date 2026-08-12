package service

import (
	"errors"
	"fmt"

	"github.com/fastsprint/task-service/internal/model"
	"github.com/fastsprint/task-service/internal/repository"
)

type TaskService interface {
	CreateIssue(reporterID uint, req model.CreateIssueRequest) (*model.Issue, error)
	GetIssueByID(id uint) (*model.Issue, error)
	GetIssueByKey(issueKey string) (*model.Issue, error)
	UpdateIssue(id, userID uint, req model.UpdateIssueRequest) (*model.Issue, error)
	DeleteIssue(id, userID uint) error
	ListIssues(query model.IssueQuery) ([]model.Issue, int64, error)
	GetBoardData(projectID uint, statuses []uint) (*model.BoardData, error)

	AddComment(issueID, authorID uint, content string) (*model.IssueComment, error)
	GetComments(issueID uint) ([]model.IssueComment, error)
	DeleteComment(id, authorID uint) error

	GetChangelogs(issueID uint) ([]model.IssueChangelog, error)

	CreateSprint(req model.CreateSprintRequest) (*model.Sprint, error)
	GetSprintByID(id uint) (*model.Sprint, error)
	UpdateSprint(id uint, req model.UpdateSprintRequest) (*model.Sprint, error)
	DeleteSprint(id uint) error
	ListSprints(query model.SprintQuery) ([]model.Sprint, int64, error)
	GetActiveSprint(projectID uint) (*model.Sprint, error)
	StartSprint(id uint) error
	CompleteSprint(id uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateIssue(reporterID uint, req model.CreateIssueRequest) (*model.Issue, error) {
	nextNum, err := s.taskRepo.GetNextIssueNumber(req.ProjectID)
	if err != nil {
		return nil, errors.New("获取任务编号失败")
	}

	issue := &model.Issue{
		ProjectID:   req.ProjectID,
		IssueKey:    fmt.Sprintf("PROJ-%d", nextNum),
		Summary:     req.Summary,
		Description: req.Description,
		IssueTypeID: req.IssueTypeID,
		StatusID:    1,
		PriorityID:  req.PriorityID,
		ReporterID:  reporterID,
		AssigneeID:  req.AssigneeID,
		ParentID:    req.ParentID,
		SprintID:    req.SprintID,
		StoryPoints: req.StoryPoints,
		DueDate:     req.DueDate,
	}

	err = s.taskRepo.CreateIssue(issue)
	if err != nil {
		return nil, err
	}

	s.addChangelog(issue.ID, reporterID, "创建", "", issue.IssueKey)

	return issue, nil
}

func (s *taskService) GetIssueByID(id uint) (*model.Issue, error) {
	return s.taskRepo.GetIssueByID(id)
}

func (s *taskService) GetIssueByKey(issueKey string) (*model.Issue, error) {
	return s.taskRepo.GetIssueByKey(issueKey)
}

func (s *taskService) UpdateIssue(id, userID uint, req model.UpdateIssueRequest) (*model.Issue, error) {
	issue, err := s.taskRepo.GetIssueByID(id)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	oldStatusID := issue.StatusID

	if req.Summary != "" {
		if issue.Summary != req.Summary {
			s.addChangelog(issue.ID, userID, "标题", issue.Summary, req.Summary)
		}
		issue.Summary = req.Summary
	}
	if req.Description != "" {
		issue.Description = req.Description
	}
	if req.IssueTypeID != nil {
		issue.IssueTypeID = *req.IssueTypeID
	}
	if req.StatusID != nil {
		issue.StatusID = *req.StatusID
	}
	if req.PriorityID != nil {
		issue.PriorityID = req.PriorityID
	}
	if req.AssigneeID != nil {
		issue.AssigneeID = req.AssigneeID
	}
	if req.SprintID != nil {
		issue.SprintID = req.SprintID
	}
	if req.StoryPoints != nil {
		issue.StoryPoints = req.StoryPoints
	}
	if req.DueDate != nil {
		issue.DueDate = req.DueDate
	}

	if oldStatusID != issue.StatusID {
		s.addChangelog(issue.ID, userID, "状态",
			fmt.Sprintf("%d", oldStatusID), fmt.Sprintf("%d", issue.StatusID))
	}

	err = s.taskRepo.UpdateIssue(issue)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

func (s *taskService) DeleteIssue(id, userID uint) error {
	_, err := s.taskRepo.GetIssueByID(id)
	if err != nil {
		return errors.New("任务不存在")
	}
	return s.taskRepo.DeleteIssue(id)
}

func (s *taskService) ListIssues(query model.IssueQuery) ([]model.Issue, int64, error) {
	return s.taskRepo.ListIssues(query)
}

func (s *taskService) GetBoardData(projectID uint, statusIDs []uint) (*model.BoardData, error) {
	issues, err := s.taskRepo.GetIssuesByStatus(projectID, statusIDs)
	if err != nil {
		return nil, err
	}

	columnMap := make(map[uint][]model.Issue)
	for _, issue := range issues {
		columnMap[issue.StatusID] = append(columnMap[issue.StatusID], issue)
	}

	columns := make([]model.BoardColumn, 0)
	for _, statusID := range statusIDs {
		columns = append(columns, model.BoardColumn{
			StatusID: statusID,
			Issues:   columnMap[statusID],
		})
	}

	return &model.BoardData{Columns: columns}, nil
}

func (s *taskService) AddComment(issueID, authorID uint, content string) (*model.IssueComment, error) {
	comment := &model.IssueComment{
		IssueID:  issueID,
		AuthorID: authorID,
		Content:  content,
	}

	err := s.taskRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *taskService) GetComments(issueID uint) ([]model.IssueComment, error) {
	return s.taskRepo.GetCommentsByIssue(issueID)
}

func (s *taskService) DeleteComment(id, authorID uint) error {
	result := s.taskRepo.DeleteComment(id, authorID)
	return result
}

func (s *taskService) addChangelog(issueID, authorID uint, field, oldValue, newValue string) {
	changelog := &model.IssueChangelog{
		IssueID:  issueID,
		AuthorID: authorID,
		Field:    field,
		OldValue: oldValue,
		NewValue: newValue,
	}
	s.taskRepo.CreateChangelog(changelog)
}

func (s *taskService) GetChangelogs(issueID uint) ([]model.IssueChangelog, error) {
	return s.taskRepo.GetChangelogsByIssue(issueID)
}

func (s *taskService) CreateSprint(req model.CreateSprintRequest) (*model.Sprint, error) {
	sprint := &model.Sprint{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Status:      "planning",
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := s.taskRepo.CreateSprint(sprint)
	if err != nil {
		return nil, err
	}

	return sprint, nil
}

func (s *taskService) GetSprintByID(id uint) (*model.Sprint, error) {
	return s.taskRepo.GetSprintByID(id)
}

func (s *taskService) UpdateSprint(id uint, req model.UpdateSprintRequest) (*model.Sprint, error) {
	sprint, err := s.taskRepo.GetSprintByID(id)
	if err != nil {
		return nil, errors.New("冲刺不存在")
	}

	if req.Name != "" {
		sprint.Name = req.Name
	}
	if req.Description != "" {
		sprint.Description = req.Description
	}
	if req.Status != "" {
		sprint.Status = req.Status
	}
	if req.StartDate != nil {
		sprint.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		sprint.EndDate = req.EndDate
	}

	err = s.taskRepo.UpdateSprint(sprint)
	if err != nil {
		return nil, err
	}

	return sprint, nil
}

func (s *taskService) DeleteSprint(id uint) error {
	_, err := s.taskRepo.GetSprintByID(id)
	if err != nil {
		return errors.New("冲刺不存在")
	}
	return s.taskRepo.DeleteSprint(id)
}

func (s *taskService) ListSprints(query model.SprintQuery) ([]model.Sprint, int64, error) {
	return s.taskRepo.ListSprints(query)
}

func (s *taskService) GetActiveSprint(projectID uint) (*model.Sprint, error) {
	return s.taskRepo.GetActiveSprint(projectID)
}

func (s *taskService) StartSprint(id uint) error {
	sprint, err := s.taskRepo.GetSprintByID(id)
	if err != nil {
		return errors.New("冲刺不存在")
	}

	if sprint.Status != "planning" {
		return errors.New("只能启动规划中的冲刺")
	}

	sprint.Status = "active"
	return s.taskRepo.UpdateSprint(sprint)
}

func (s *taskService) CompleteSprint(id uint) error {
	sprint, err := s.taskRepo.GetSprintByID(id)
	if err != nil {
		return errors.New("冲刺不存在")
	}

	if sprint.Status != "active" {
		return errors.New("只能完成进行中的冲刺")
	}

	sprint.Status = "completed"
	return s.taskRepo.UpdateSprint(sprint)
}
