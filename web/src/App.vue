<template>
  <div id="app" class="min-h-screen bg-gray-100">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <h1 class="text-3xl font-bold text-gray-900">TimeLog</h1>
          <button
            @click="toggleForm"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <PlusIcon class="h-5 w-5 mr-2" />
            New Log
          </button>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <div class="space-y-6">
        <TimeLogForm
          v-if="showForm"
          :editing-log="editingLog"
          :submitting="submitting"
          :available-tags="tags"
          @submit="handleSubmit"
          @cancel="handleCancel"
        />
        
        <TimeLogList
          :time-logs="timeLogs"
          :loading="loading"
          :error="error"
          @edit="handleEdit"
          @delete="handleDelete"
        />
      </div>
    </main>

    <div
      v-if="notification.show"
      class="fixed bottom-4 right-4 bg-white border border-gray-200 rounded-lg shadow-lg p-4 max-w-sm"
      :class="{
        'border-green-200 bg-green-50': notification.type === 'success',
        'border-red-200 bg-red-50': notification.type === 'error'
      }"
    >
      <div class="flex items-center">
        <CheckCircleIcon
          v-if="notification.type === 'success'"
          class="h-5 w-5 text-green-600 mr-2"
        />
        <XCircleIcon
          v-if="notification.type === 'error'"
          class="h-5 w-5 text-red-600 mr-2"
        />
        <p class="text-sm font-medium" :class="{
          'text-green-800': notification.type === 'success',
          'text-red-800': notification.type === 'error'
        }">
          {{ notification.message }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusIcon, CheckCircleIcon, XCircleIcon } from '@heroicons/vue/24/outline'
import TimeLogList from '@/components/TimeLogList.vue'
import TimeLogForm from '@/components/TimeLogForm.vue'
import { timelogAPI, tagAPI } from '@/api'
import type { TimeLog, Tag, CreateTimeLogRequest, UpdateTimeLogRequest } from '@/types'

const timeLogs = ref<TimeLog[]>([])
const tags = ref<Tag[]>([])
const loading = ref(false)
const submitting = ref(false)
const error = ref<string | null>(null)
const showForm = ref(false)
const editingLog = ref<TimeLog | undefined>()

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

const loadTimeLogs = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await timelogAPI.getAll()
    timeLogs.value = response.data || []
  } catch (err) {
    error.value = 'Failed to load time logs'
    console.error('Error loading time logs:', err)
    showNotification('error', 'Failed to load time logs')
  } finally {
    loading.value = false
  }
}

const loadTags = async () => {
  try {
    const response = await tagAPI.getAll()
    tags.value = response.data || []
  } catch (err) {
    console.error('Error loading tags:', err)
    showNotification('error', 'Failed to load tags')
  }
}

const toggleForm = () => {
  showForm.value = !showForm.value
  if (!showForm.value) {
    editingLog.value = undefined
  }
}

const handleSubmit = async (data: CreateTimeLogRequest | UpdateTimeLogRequest) => {
  submitting.value = true
  
  try {
    if (editingLog.value) {
      await timelogAPI.update(editingLog.value.id, data as UpdateTimeLogRequest)
      showNotification('success', 'Time log updated successfully')
    } else {
      await timelogAPI.create(data as CreateTimeLogRequest)
      showNotification('success', 'Time log created successfully')
    }
    
    await loadTimeLogs()
    showForm.value = false
    editingLog.value = undefined
  } catch (err) {
    console.error('Error saving time log:', err)
    showNotification('error', 'Failed to save time log')
  } finally {
    submitting.value = false
  }
}

const handleEdit = (log: TimeLog) => {
  editingLog.value = log
  showForm.value = true
}

const handleCancel = () => {
  showForm.value = false
  editingLog.value = undefined
}

const handleDelete = async (id: number) => {
  if (!confirm('Are you sure you want to delete this time log?')) {
    return
  }
  
  try {
    await timelogAPI.delete(id)
    showNotification('success', 'Time log deleted successfully')
    await loadTimeLogs()
  } catch (err) {
    console.error('Error deleting time log:', err)
    showNotification('error', 'Failed to delete time log')
  }
}

onMounted(() => {
  loadTimeLogs()
  loadTags()
})
</script>