import { defineStore } from 'pinia'
import { ref } from 'vue'
import { timelogAPI } from '@/api'
import type { TimeLog } from '@/types'

interface TimeLogCache {
  data: TimeLog[]
  timestamp: number
}

const TIMELOG_CACHE_KEY = 'timelog-cache'
const CACHE_DURATION = 25 * 60 * 1000 // 25 minutes in milliseconds

export const useTimeLogStore = defineStore('timelog', () => {
  const timeLogs = ref<TimeLog[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Load timelogs from cache if valid, otherwise fetch from API
  const loadTimeLogs = async (forceRefresh = false) => {
    if (!forceRefresh && isCacheValid()) {
      const cached = loadFromCache()
      if (cached) {
        timeLogs.value = cached
        return
      }
    }

    loading.value = true
    error.value = null

    try {
      const response = await timelogAPI.getAll()
      const logs = response.data || []
      // Sort by start_time in descending order (most recent first)
      timeLogs.value = logs.sort(
        (a, b) => new Date(b.start_time).getTime() - new Date(a.start_time).getTime()
      )
      // Save to cache
      saveToCache(timeLogs.value)
    } catch (err) {
      error.value = 'Failed to load time logs'
      console.error('Error loading time logs:', err)
    } finally {
      loading.value = false
    }
  }

  // Check if cache is valid
  const isCacheValid = (): boolean => {
    try {
      const cached = localStorage.getItem(TIMELOG_CACHE_KEY)
      if (!cached) return false

      const parsed: TimeLogCache = JSON.parse(cached)
      const now = Date.now()
      return now - parsed.timestamp < CACHE_DURATION
    } catch (error) {
      console.error('Error checking cache validity:', error)
      return false
    }
  }

  // Load from cache
  const loadFromCache = (): TimeLog[] | null => {
    try {
      const cached = localStorage.getItem(TIMELOG_CACHE_KEY)
      if (!cached) return null

      const parsed: TimeLogCache = JSON.parse(cached)
      return parsed.data
    } catch (error) {
      console.error('Error loading from cache:', error)
      return null
    }
  }

  // Save to cache
  const saveToCache = (data: TimeLog[]): void => {
    try {
      const cache: TimeLogCache = {
        data,
        timestamp: Date.now(),
      }
      localStorage.setItem(TIMELOG_CACHE_KEY, JSON.stringify(cache))
    } catch (error) {
      console.error('Error saving to cache:', error)
    }
  }

  // Refresh timelogs (force reload)
  const refreshTimeLogs = async () => {
    await loadTimeLogs(true)
  }

  return {
    timeLogs,
    loading,
    error,
    loadTimeLogs,
    refreshTimeLogs,
    isCacheValid,
  }
})
