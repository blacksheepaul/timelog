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
  remarks: string
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

export interface CreateTimeLogRequest {
  start_time: string
  end_time?: string
  tag_id: number
  remarks: string
}

export interface UpdateTimeLogRequest {
  start_time?: string
  end_time?: string
  tag_id?: number
  remarks?: string
}

export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}