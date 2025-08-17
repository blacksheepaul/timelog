<template>
  <div class="space-y-6">
    <!-- 欢迎信息 -->
    <div class="bg-white shadow rounded-lg p-6">
      <h2 class="text-2xl font-bold text-gray-900 mb-4">Welcome to TimeLog Dashboard</h2>
      <p class="text-gray-600 mb-6">
        Track your time efficiently with our intuitive time logging system. 
        Get insights into your productivity and manage your daily activities.
      </p>
      
      <!-- 快速操作按钮 -->
      <div class="flex flex-wrap gap-4">
        <router-link
          to="/timelogs"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <ClockIcon class="h-5 w-5 mr-2" />
          View Time Logs
        </router-link>
        <router-link
          to="/tags"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <TagIcon class="h-5 w-5 mr-2" />
          Manage Tags
        </router-link>
        <router-link
          to="/statistics"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <ChartBarIcon class="h-5 w-5 mr-2" />
          View Statistics
        </router-link>
      </div>
    </div>

    <!-- 今日概览 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <ClockIcon class="h-8 w-8 text-blue-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Today's Logs</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.count }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <PlayIcon class="h-8 w-8 text-green-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Active Sessions</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.activeSessions }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <StopIcon class="h-8 w-8 text-red-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Total Time</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.totalTime }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg p-6">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <TagIcon class="h-8 w-8 text-purple-600" />
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-500">Tags Used</p>
            <p class="text-2xl font-semibold text-gray-900">{{ todayStats.tagsUsed }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 最近的时间记录 -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Recent Time Logs</h3>
      </div>
      <div v-if="loading" class="p-6 text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading...</p>
      </div>
      <div v-else-if="recentLogs.length === 0" class="p-6 text-center text-gray-500">
        No recent time logs found.
      </div>
      <div v-else class="divide-y divide-gray-200">
        <div v-for="log in recentLogs" :key="log.id" class="p-6 hover:bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <span
                v-if="log.tag"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                :style="{ backgroundColor: log.tag.color }"
              >
                {{ log.tag.name }}
              </span>
              <div>
                <p class="text-sm font-medium text-gray-900">{{ log.remarks || 'No remarks' }}</p>
                <p class="text-xs text-gray-500">
                  {{ formatDateTime(log.start_time) }} - 
                  {{ log.end_time ? formatDateTime(log.end_time) : 'Ongoing' }}
                </p>
              </div>
            </div>
            <span class="text-sm text-gray-500">
              {{ calculateDuration(log.start_time, log.end_time) }}
            </span>
          </div>
        </div>
      </div>
      <div class="px-6 py-3 bg-gray-50 text-center">
        <router-link
          to="/timelogs"
          class="text-sm font-medium text-blue-600 hover:text-blue-500"
        >
          View all time logs →
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { 
  ClockIcon, 
  PlayIcon, 
  StopIcon, 
  TagIcon, 
  ChartBarIcon 
} from '@heroicons/vue/24/outline'
import { timelogAPI } from '@/api'
import { formatDateTime, calculateDuration } from '@/utils/date'
import type { TimeLog } from '@/types'

const loading = ref(false)
const recentLogs = ref<TimeLog[]>([])

// 今日统计数据
const todayStats = computed(() => {
  const today = new Date().toDateString()
  const todayLogs = recentLogs.value.filter(log => 
    new Date(log.start_time).toDateString() === today
  )
  
  const activeSessions = todayLogs.filter(log => !log.end_time).length
  const tagsUsed = new Set(todayLogs.map(log => log.tag_id)).size
  
  // 计算总时间（小时）
  const totalMinutes = todayLogs
    .filter(log => log.end_time)
    .reduce((total, log) => {
      const start = new Date(log.start_time)
      const end = new Date(log.end_time!)
      return total + (end.getTime() - start.getTime()) / (1000 * 60)
    }, 0)
  
  const totalTime = totalMinutes > 60 
    ? `${Math.floor(totalMinutes / 60)}h ${Math.round(totalMinutes % 60)}m`
    : `${Math.round(totalMinutes)}m`

  return {
    count: todayLogs.length,
    activeSessions,
    totalTime,
    tagsUsed
  }
})

const loadRecentLogs = async () => {
  loading.value = true
  try {
    const response = await timelogAPI.getAll()
    // 只显示最近10条记录
    recentLogs.value = (response.data || []).slice(0, 10)
  } catch (err) {
    console.error('Error loading recent logs:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecentLogs()
})
</script>