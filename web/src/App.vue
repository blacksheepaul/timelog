<template>
  <div id="app" class="min-h-screen bg-gray-100">
    <!-- 顶部导航栏 -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <h1 class="text-3xl font-bold text-gray-900 mr-8">TimeLog</h1>
            <!-- 导航菜单 -->
            <nav class="hidden md:flex space-x-8">
              <router-link
                to="/"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'Home',
                }"
              >
                Dashboard
              </router-link>
              <router-link
                to="/timelogs"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'TimeLog',
                }"
              >
                Time Logs
              </router-link>
              <router-link
                to="/tasks"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'Tasks',
                }"
              >
                Tasks
              </router-link>
              <router-link
                to="/tags"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'Tags',
                }"
              >
                Tags
              </router-link>
              <router-link
                to="/statistics"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'Statistics',
                }"
              >
                Statistics
              </router-link>
              <router-link
                to="/constraints"
                class="text-gray-500 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                :class="{
                  'text-blue-600 font-semibold': $route.name === 'Constraints',
                }"
              >
                约束
              </router-link>
            </nav>
          </div>

          <!-- 移动端菜单按钮 -->
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
          >
            <Bars3Icon v-if="!mobileMenuOpen" class="h-6 w-6" />
            <XMarkIcon v-else class="h-6 w-6" />
          </button>
        </div>

        <!-- 移动端导航菜单 -->
        <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200 py-4">
          <nav class="space-y-1">
            <router-link
              to="/"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'Home',
              }"
              @click="mobileMenuOpen = false"
            >
              Dashboard
            </router-link>
            <router-link
              to="/timelogs"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'TimeLog',
              }"
              @click="mobileMenuOpen = false"
            >
              Time Logs
            </router-link>
            <router-link
              to="/tasks"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'Tasks',
              }"
              @click="mobileMenuOpen = false"
            >
              Tasks
            </router-link>
            <router-link
              to="/tags"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'Tags',
              }"
              @click="mobileMenuOpen = false"
            >
              Tags
            </router-link>
            <router-link
              to="/statistics"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'Statistics',
              }"
              @click="mobileMenuOpen = false"
            >
              Statistics
            </router-link>
            <router-link
              to="/constraints"
              class="block px-3 py-2 text-base font-medium text-gray-500 hover:text-gray-900 hover:bg-gray-50 transition-colors"
              :class="{
                'text-blue-600 bg-blue-50': $route.name === 'Constraints',
              }"
              @click="mobileMenuOpen = false"
            >
              约束
            </router-link>
          </nav>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <router-view />
    </main>

    <!-- 全局通知组件 -->
    <div
      v-if="notification.show"
      class="fixed bottom-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg p-4 max-w-sm z-50"
      :class="{
        'border-green-200 bg-green-50': notification.type === 'success',
        'border-red-200 bg-red-50': notification.type === 'error',
      }"
    >
      <div class="flex items-center">
        <CheckCircleIcon
          v-if="notification.type === 'success'"
          class="h-5 w-5 text-green-600 mr-2"
        />
        <XCircleIcon v-if="notification.type === 'error'" class="h-5 w-5 text-red-600 mr-2" />
        <p
          class="text-sm font-medium"
          :class="{
            'text-green-800': notification.type === 'success',
            'text-red-800': notification.type === 'error',
          }"
        >
          {{ notification.message }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, provide, onMounted, onUnmounted } from 'vue'
  import { CheckCircleIcon, XCircleIcon, Bars3Icon, XMarkIcon } from '@heroicons/vue/24/outline'
  import { useSettings } from '@/composables/useSettings'
  import { setNotificationHandler } from '@/api'

  // 移动端菜单状态
  const mobileMenuOpen = ref(false)

  // 全局通知系统
  const notification = reactive({
    show: false,
    type: 'success' as 'success' | 'error',
    message: '',
  })

  const showNotification = (type: 'success' | 'error', message: string) => {
    notification.type = type
    notification.message = message
    notification.show = true

    setTimeout(() => {
      notification.show = false
    }, 3000)
  }

  // 通过provide向子组件提供全局通知功能
  provide('showNotification', showNotification)

  const REMINDER_INTERVAL_MS = 25 * 60 * 1000
  let reminderTimer: number | undefined

  const sendReminderNotification = () => {
    try {
      new Notification('该记录TimeLog了', {
        body: '又过去了25分钟，别忘了记录你的时间日志。',
        tag: 'timelog-reminder',
      })
    } catch (error) {
      console.error('Failed to send reminder notification', error)
    }
  }

  const startReminderTimer = () => {
    if (reminderTimer) {
      window.clearInterval(reminderTimer)
    }
    reminderTimer = window.setInterval(sendReminderNotification, REMINDER_INTERVAL_MS)
  }

  const initSystemNotifications = () => {
    if (typeof window === 'undefined' || !('Notification' in window)) {
      showNotification('error', '当前浏览器不支持系统通知提醒')
      return
    }

    if (Notification.permission === 'granted') {
      startReminderTimer()
      return
    }

    if (Notification.permission === 'denied') {
      showNotification('error', '请在浏览器设置中允许通知以启用提醒')
      return
    }

    Notification.requestPermission().then(permission => {
      if (permission === 'granted') {
        startReminderTimer()
        showNotification('success', '已启用每25分钟一次的记录提醒')
      } else {
        showNotification('error', '未授予通知权限，无法启用提醒')
      }
    })
  }

  // Initialize settings on app mount
  onMounted(() => {
    const { loadSettings } = useSettings()
    loadSettings()
    initSystemNotifications()

    // Register notification handler for API timeout errors
    setNotificationHandler(showNotification)
  })

  onUnmounted(() => {
    if (reminderTimer) {
      window.clearInterval(reminderTimer)
    }
  })
</script>
