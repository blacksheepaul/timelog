import axios from 'axios'
import type { TimeLog, CreateTimeLogRequest, UpdateTimeLogRequest, ApiResponse } from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const timelogAPI = {
  getAll: (): Promise<ApiResponse<TimeLog[]>> => 
    api.get('/timelogs').then(res => res.data),
  
  getById: (id: number): Promise<ApiResponse<TimeLog>> => 
    api.get(`/timelogs/${id}`).then(res => res.data),
  
  create: (data: CreateTimeLogRequest): Promise<ApiResponse<TimeLog>> => 
    api.post('/timelogs', data).then(res => res.data),
  
  update: (id: number, data: UpdateTimeLogRequest): Promise<ApiResponse<TimeLog>> => 
    api.put(`/timelogs/${id}`, data).then(res => res.data),
  
  delete: (id: number): Promise<ApiResponse<null>> => 
    api.delete(`/timelogs/${id}`).then(res => res.data),
}

export default api