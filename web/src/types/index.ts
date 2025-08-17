export interface Tag {
  id: number
  name: string
  color: string
  description: string
  created_at: string
  updated_at: string
}

export interface TimeLog {
  id: number
  start_time: string
  end_time?: string | null
  tag_id: number
  tag: Tag
  task_id?: number | null
  task?: Task | null
  remarks: string
  created_at: string
  updated_at: string
  deleted_at?: string | null
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