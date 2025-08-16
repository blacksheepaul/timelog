export interface TimeLog {
  id: number
  start_time: string
  end_time?: string
  tags: string
  remarks: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface CreateTimeLogRequest {
  start_time: string
  end_time?: string
  tags: string
  remarks: string
}

export interface UpdateTimeLogRequest {
  start_time?: string
  end_time?: string
  tags?: string
  remarks?: string
}

export interface ApiResponse<T> {
  data: T
  message?: string
  status: number
}