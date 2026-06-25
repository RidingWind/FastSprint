import request from '../utils/request';
import {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  User,
  Project,
  ProjectMember,
  Issue,
  IssueComment,
  IssueChangelog,
  Sprint,
  Status,
  IssueType,
  Priority,
  BoardData,
  PageResult,
} from '../types';

export const authApi = {
  login: (data: LoginRequest) =>
    request.post<LoginResponse>('/auth/login', data).then((res) => res.data.data!),

  register: (data: RegisterRequest) =>
    request.post('/auth/register', data).then((res) => res.data.data!),

  getCurrentUser: () =>
    request.get<User>('/users/me').then((res) => res.data.data!),
};

export const userApi = {
  getUserById: (id: number) =>
    request.get<User>(`/users/${id}`).then((res) => res.data.data!),

  listUsers: (params?: any) =>
    request.get<PageResult<User>>('/users', { params }).then((res) => res.data.data!),
};

export const projectApi = {
  createProject: (data: any) =>
    request.post<Project>('/projects', data).then((res) => res.data.data!),

  getProject: (id: number) =>
    request.get<Project>(`/projects/${id}`).then((res) => res.data.data!),

  getProjectByKey: (key: string) =>
    request.get<Project>(`/projects/key/${key}`).then((res) => res.data.data!),

  listProjects: (params?: any) =>
    request.get<PageResult<Project>>('/projects', { params }).then((res) => res.data.data!),

  updateProject: (id: number, data: any) =>
    request.put<Project>(`/projects/${id}`, data).then((res) => res.data.data!),

  deleteProject: (id: number) =>
    request.delete(`/projects/${id}`).then((res) => res.data.data!),

  listMembers: (projectId: number) =>
    request.get<ProjectMember[]>(`/projects/${projectId}/members`).then((res) => res.data.data!),

  addMember: (projectId: number, data: { user_id: number; role: string }) =>
    request.post(`/projects/${projectId}/members`, data).then((res) => res.data.data!),

  removeMember: (projectId: number, userId: number) =>
    request.delete(`/projects/${projectId}/members/${userId}`).then((res) => res.data.data!),

  updateMemberRole: (projectId: number, userId: number, role: string) =>
    request.put(`/projects/${projectId}/members/${userId}/role`, { role }).then((res) => res.data.data!),

  listStatuses: (projectId: number) =>
    request.get<Status[]>(`/projects/${projectId}/statuses`).then((res) => res.data.data!),

  listIssueTypes: (projectId: number) =>
    request.get<IssueType[]>(`/projects/${projectId}/issue-types`).then((res) => res.data.data!),

  listPriorities: (projectId: number) =>
    request.get<Priority[]>(`/projects/${projectId}/priorities`).then((res) => res.data.data!),
};

export const issueApi = {
  createIssue: (data: any) =>
    request.post<Issue>('/issues', data).then((res) => res.data.data!),

  getIssue: (id: number) =>
    request.get<Issue>(`/issues/${id}`).then((res) => res.data.data!),

  getIssueByKey: (key: string) =>
    request.get<Issue>(`/issues/key/${key}`).then((res) => res.data.data!),

  listIssues: (params?: any) =>
    request.get<PageResult<Issue>>('/issues', { params }).then((res) => res.data.data!),

  updateIssue: (id: number, data: any) =>
    request.put<Issue>(`/issues/${id}`, data).then((res) => res.data.data!),

  deleteIssue: (id: number) =>
    request.delete(`/issues/${id}`).then((res) => res.data.data!),

  getComments: (issueId: number) =>
    request.get<IssueComment[]>(`/issues/${issueId}/comments`).then((res) => res.data.data!),

  addComment: (issueId: number, content: string) =>
    request.post<IssueComment>(`/issues/${issueId}/comments`, { content }).then((res) => res.data.data!),

  getChangelogs: (issueId: number) =>
    request.get<IssueChangelog[]>(`/issues/${issueId}/changelogs`).then((res) => res.data.data!),
};

export const sprintApi = {
  createSprint: (data: any) =>
    request.post<Sprint>('/sprints', data).then((res) => res.data.data!),

  getSprint: (id: number) =>
    request.get<Sprint>(`/sprints/${id}`).then((res) => res.data.data!),

  listSprints: (params?: any) =>
    request.get<PageResult<Sprint>>('/sprints', { params }).then((res) => res.data.data!),

  updateSprint: (id: number, data: any) =>
    request.put<Sprint>(`/sprints/${id}`, data).then((res) => res.data.data!),

  deleteSprint: (id: number) =>
    request.delete(`/sprints/${id}`).then((res) => res.data.data!),

  startSprint: (id: number) =>
    request.post(`/sprints/${id}/start`).then((res) => res.data.data!),

  completeSprint: (id: number) =>
    request.post(`/sprints/${id}/complete`).then((res) => res.data.data!),
};

export const boardApi = {
  getBoardData: (projectId: number) =>
    request.get<BoardData>('/board', { params: { project_id: projectId } }).then((res) => res.data.data!),
};
