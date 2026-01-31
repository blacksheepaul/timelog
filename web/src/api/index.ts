import axios from 'axios'
import type {
  TimeLog,
  Tag,
  Task,
  CreateTimeLogRequest,
  UpdateTimeLogRequest,
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskStats,
  ApiResponse,
  Constraint,
  CreateConstraintRequest,
  UpdateConstraintRequest,
} from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Store notification handler to be set by App.vue
let notificationHandler: ((type: 'success' | 'error', message: string) => void) | null = null

export const setNotificationHandler = (
  handler: (type: 'success' | 'error', message: string) => void
) => {
  notificationHandler = handler
}

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error)

    // Check if the error is a timeout or cancellation error
    // These can have various code/message combinations depending on the scenario
    const isTimeoutOrCancellation =
      error.code === 'ECONNABORTED' || // Request abort (timeout)
      axios.isCancel?.(error) || // Explicit cancellation via axios cancel token/AbortController
      error.message?.includes('timeout') || // Message contains timeout
      error.message?.includes('Timeout') // Case-insensitive

    if (isTimeoutOrCancellation && notificationHandler) {
      // Show browser notification for timeout
      notificationHandler('error', 'Request timeout - The server took too long to respond')
    }

    return Promise.reject(error)
  }
)

export const timelogAPI = {
  getAll: (): Promise<ApiResponse<TimeLog[]>> => api.get('/timelogs').then(res => res.data),

  getRecent: (limit: number = 5): Promise<ApiResponse<TimeLog[]>> =>
    api
      .get(`/timelogs?limit=${limit}&order=${encodeURIComponent('created_at DESC')}`)
      .then(res => res.data),

  getById: (id: number): Promise<ApiResponse<TimeLog>> =>
    api.get(`/timelogs/${id}`).then(res => res.data),

  create: (data: CreateTimeLogRequest): Promise<ApiResponse<TimeLog>> =>
    api.post('/timelogs', data).then(res => res.data),

  update: (id: number, data: UpdateTimeLogRequest): Promise<ApiResponse<TimeLog>> =>
    api.put(`/timelogs/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/timelogs/${id}`).then(res => res.data),
}

export const tagAPI = {
  getAll: (): Promise<ApiResponse<Tag[]>> => api.get('/tags').then(res => res.data),

  getById: (id: number): Promise<ApiResponse<Tag>> => api.get(`/tags/${id}`).then(res => res.data),

  create: (data: Partial<Tag>): Promise<ApiResponse<Tag>> =>
    api.post('/tags', data).then(res => res.data),

  update: (id: number, data: Partial<Tag>): Promise<ApiResponse<Tag>> =>
    api.put(`/tags/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/tags/${id}`).then(res => res.data),
}

export const taskAPI = {
  getAll: (
    date?: string,
    includeSuspended?: boolean,
    includeCompleted?: boolean
  ): Promise<ApiResponse<Task[]>> => {
    const params = new URLSearchParams()
    if (date) params.append('date', date)
    if (includeSuspended) params.append('include_suspended', 'true')
    if (includeCompleted) params.append('include_completed', 'true')
    const url = `/tasks${params.toString() ? '?' + params.toString() : ''}`
    return api.get(url).then(res => res.data)
  },

  getById: (id: number): Promise<ApiResponse<Task>> =>
    api.get(`/tasks/${id}`).then(res => res.data),

  create: (data: CreateTaskRequest): Promise<ApiResponse<Task>> =>
    api.post('/tasks', data).then(res => res.data),

  update: (id: number, data: UpdateTaskRequest): Promise<ApiResponse<Task>> =>
    api.put(`/tasks/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/tasks/${id}`).then(res => res.data),

  complete: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/complete`).then(res => res.data),

  incomplete: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/incomplete`).then(res => res.data),

  suspend: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/suspend`).then(res => res.data),

  unsuspend: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/tasks/${id}/unsuspend`).then(res => res.data),

  getStats: (date: string): Promise<ApiResponse<TaskStats>> =>
    api.get(`/tasks/stats/${date}`).then(res => res.data),
}

export const constraintAPI = {
  getAll: (active?: boolean): Promise<ApiResponse<Constraint[]>> => {
    const url = active !== undefined ? `/constraints?active=${active}` : '/constraints'
    return api.get(url).then(res => res.data)
  },

  getById: (id: number): Promise<ApiResponse<Constraint>> =>
    api.get(`/constraints/${id}`).then(res => res.data),

  create: (data: CreateConstraintRequest): Promise<ApiResponse<Constraint>> =>
    api.post('/constraints', data).then(res => res.data),

  update: (id: number, data: UpdateConstraintRequest): Promise<ApiResponse<Constraint>> =>
    api.put(`/constraints/${id}`, data).then(res => res.data),

  delete: (id: number): Promise<ApiResponse<null>> =>
    api.delete(`/constraints/${id}`).then(res => res.data),

  complete: (id: number, endReason: string): Promise<ApiResponse<null>> =>
    api.post(`/constraints/${id}/complete`, { end_reason: endReason }).then(res => res.data),

  reactivate: (id: number): Promise<ApiResponse<null>> =>
    api.post(`/constraints/${id}/reactivate`).then(res => res.data),
}

export default api
