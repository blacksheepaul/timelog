export interface Tag {
  id: number
  name: string
  color: string
  description: string
  created_at: string
  updated_at: string
}

export interface TimeLog {
  ID: number
  start_time: string
  end_time?: string | null
  tag_id: number
  tag: Tag
  task_id?: number | null
  task?: Task | null
  remarks: string
  CreatedAt: string
  UpdatedAt: string
  DeletedAt?: string | null
}

export interface CreateTimeLogRequest {
  start_time: string
  end_time?: string
  tag_id: number
  task_id?: number | null
  remarks: string
}

export interface UpdateTimeLogRequest {
  start_time?: string
  end_time?: string
  tag_id?: number
  task_id?: number | null
  remarks?: string
}

export interface Task {
  id: number
  title: string
  description: string
  tag_id: number
  tag?: Tag
  due_date: string
  estimated_minutes: number
  is_completed: boolean
  completed_at?: string | null
  is_suspended: boolean
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  tag_id: number
  due_date: string
  estimated_minutes: number
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  tag_id?: number
  due_date?: string
  estimated_minutes?: number
  is_completed?: boolean
}

export interface TaskStats {
  total_tasks: number
  completed_tasks: number
  completion_rate: number
}

export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}

export interface Constraint {
  id: number
  description: string
  end_reason?: string
  punishment_quote: string
  start_date: string
  end_date?: string | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateConstraintRequest {
  description: string
  punishment_quote: string
  start_date: string
  end_date?: string
}

export interface UpdateConstraintRequest {
  description?: string
  punishment_quote?: string
  start_date?: string
  end_date?: string
}
