export interface User {
  id: number;
  username: string;
  email: string;
  display_name: string;
  avatar_url: string;
  status: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  display_name?: string;
}

export interface Project {
  id: number;
  key: string;
  name: string;
  description: string;
  owner_id: number;
  status: string;
  default_assignee_id: number;
  created_at: string;
  updated_at: string;
}

export interface ProjectMember {
  id: number;
  project_id: number;
  user_id: number;
  role: string;
  joined_at: string;
}

export interface Status {
  id: number;
  project_id: number;
  name: string;
  category: string;
  color: string;
  sort_order: number;
  is_initial: boolean;
}

export interface IssueType {
  id: number;
  project_id: number;
  name: string;
  description: string;
  icon: string;
  color: string;
  is_subtask: boolean;
  sort_order: number;
}

export interface Priority {
  id: number;
  project_id: number;
  name: string;
  description: string;
  icon: string;
  color: string;
  sort_order: number;
}

export interface Issue {
  id: number;
  project_id: number;
  issue_key: string;
  summary: string;
  description: string;
  issue_type_id: number;
  status_id: number;
  priority_id?: number;
  reporter_id: number;
  assignee_id?: number;
  parent_id?: number;
  sprint_id?: number;
  story_points?: number;
  due_date?: string;
  created_at: string;
  updated_at: string;
}

export interface IssueComment {
  id: number;
  issue_id: number;
  author_id: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface IssueChangelog {
  id: number;
  issue_id: number;
  author_id: number;
  field: string;
  old_value: string;
  new_value: string;
  created_at: string;
}

export interface Sprint {
  id: number;
  project_id: number;
  name: string;
  description: string;
  status: string;
  start_date?: string;
  end_date?: string;
  created_at: string;
  updated_at: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data?: T;
}

export interface BoardColumn {
  status_id: number;
  name: string;
  issues: Issue[];
}

export interface BoardData {
  columns: BoardColumn[];
}
