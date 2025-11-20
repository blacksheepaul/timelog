<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Task Management</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        New Task
      </button>
    </div>

    <!-- 任务创建/编辑表单 -->
    <div v-if="showForm" class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-6">
        {{ isEditing ? 'Edit Task' : 'Create New Task' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="title" class="block text-sm font-medium text-gray-700 mb-2"> Title * </label>
          <input
            id="title"
            v-model="form.title"
            type="text"
            required
            maxlength="200"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Task title"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Describe what this task involves..."
          ></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label for="tag_id" class="block text-sm font-medium text-gray-700 mb-2"> Tag * </label>
            <select
              id="tag_id"
              v-model="form.tag_id"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="" disabled>Select a tag</option>
              <option
                v-for="tag in availableTags"
                :key="tag.id"
                :value="tag.id"
                :style="{ color: tag.color }"
              >
                {{ tag.name }} - {{ tag.description }}
              </option>
            </select>
          </div>

          <div>
            <label for="due_date" class="block text-sm font-medium text-gray-700 mb-2">
              Due Date *
            </label>
            <input
              id="due_date"
              v-model="form.due_date"
              type="date"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="estimated_minutes" class="block text-sm font-medium text-gray-700 mb-2">
              Estimated Time (minutes) *
            </label>
            <input
              id="estimated_minutes"
              v-model.number="form.estimated_minutes"
              type="number"
              min="1"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="e.g. 60"
            />
          </div>
        </div>

        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="handleCancel"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ submitting ? 'Saving...' : isEditing ? 'Update' : 'Create' }}
          </button>
        </div>
      </form>
    </div>

    <!-- 任务列表 -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <h2 class="text-lg font-medium text-gray-900">Tasks</h2>
          <div class="flex items-center space-x-4">
            <select
              v-model="dateFilter"
              @change="loadTasks"
              class="px-3 py-1 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">All dates</option>
              <option value="today">Today</option>
              <option value="tomorrow">Tomorrow</option>
              <option value="this-week">This week</option>
            </select>
          </div>
        </div>
      </div>

      <div v-if="loading" class="p-6 text-center">
        <div
          class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
        ></div>
        <p class="mt-2 text-gray-600">Loading...</p>
      </div>

      <div v-else-if="error" class="p-6 text-center text-red-600">
        {{ error }}
      </div>

      <div v-else-if="tasks.length === 0" class="p-6 text-center text-gray-500">
        No tasks found. Create your first one!
      </div>

      <div v-else class="divide-y divide-gray-200">
        <div
          v-for="task in tasks"
          :key="task.id"
          class="p-6 hover:bg-gray-50"
          :class="{ 'opacity-60': task.is_completed }"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-2">
                <h3
                  class="text-lg font-medium"
                  :class="task.is_completed ? 'line-through text-gray-500' : 'text-gray-900'"
                >
                  {{ task.title }}
                </h3>
                <span
                  v-if="task.tag"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                  :style="{ backgroundColor: task.tag.color }"
                >
                  {{ task.tag.name }}
                </span>
                <span
                  v-if="task.is_completed"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  ✓ Completed
                </span>
              </div>

              <p v-if="task.description" class="text-gray-600 mb-3">
                {{ task.description }}
              </p>

              <div class="flex items-center space-x-4 text-sm text-gray-500">
                <span>Due: {{ formatDate(task.due_date) }}</span>
                <span>Estimated: {{ task.estimated_minutes }}min</span>
                <span v-if="task.completed_at">
                  Completed: {{ formatDateTime(task.completed_at) }}
                </span>
              </div>
            </div>

            <div class="flex items-center space-x-2 ml-4">
              <button
                v-if="!task.is_completed"
                @click="completeTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                Complete
              </button>
              <button
                v-else
                @click="incompleteTask(task.id)"
                class="px-3 py-1 text-xs font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                Reopen
              </button>
              <button
                @click="handleEdit(task)"
                class="text-blue-600 hover:text-blue-900 text-sm font-medium"
              >
                Edit
              </button>
              <button
                @click="handleDelete(task.id)"
                class="text-red-600 hover:text-red-900 text-sm font-medium"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed, inject } from 'vue'
  import { PlusIcon } from '@heroicons/vue/24/outline'
  import { formatDate, formatDateTime } from '@/utils/date'
  import { taskAPI, tagAPI } from '@/api'
  import type { Task, Tag, CreateTaskRequest, UpdateTaskRequest } from '@/types'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  const tasks = ref<Task[]>([])
  const availableTags = ref<Tag[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const error = ref<string | null>(null)
  const showForm = ref(false)
  const editingTask = ref<Task | undefined>()
  const dateFilter = ref('')

  const isEditing = computed(() => !!editingTask.value)

  const form = reactive({
    title: '',
    description: '',
    tag_id: '',
    due_date: '',
    estimated_minutes: 60,
  })

  const resetForm = () => {
    form.title = ''
    form.description = ''
    form.tag_id = ''
    form.due_date = new Date().toISOString().split('T')[0] // Today's date
    form.estimated_minutes = 60
  }

  const loadEditingData = () => {
    if (editingTask.value) {
      form.title = editingTask.value.title
      form.description = editingTask.value.description
      form.tag_id = editingTask.value.tag_id.toString()
      form.due_date = editingTask.value.due_date.split('T')[0]
      form.estimated_minutes = editingTask.value.estimated_minutes
    } else {
      resetForm()
    }
  }

  // API 调用函数已从 @/api 导入

  const loadTasks = async () => {
    loading.value = true
    error.value = null

    try {
      let dateParam = ''
      const today = new Date().toISOString().split('T')[0]

      if (dateFilter.value === 'today') {
        dateParam = today
      } else if (dateFilter.value === 'tomorrow') {
        const tomorrow = new Date()
        tomorrow.setDate(tomorrow.getDate() + 1)
        dateParam = tomorrow.toISOString().split('T')[0]
      }

      const response = await taskAPI.getAll(dateParam)
      tasks.value = response.data || []
    } catch (err) {
      error.value = 'Failed to load tasks'
      console.error('Error loading tasks:', err)
      showNotification('error', 'Failed to load tasks')
    } finally {
      loading.value = false
    }
  }

  const loadTags = async () => {
    try {
      const response = await tagAPI.getAll()
      availableTags.value = response.data || []
    } catch (err) {
      console.error('Error loading tags:', err)
      showNotification('error', 'Failed to load tags')
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (!showForm.value) {
      editingTask.value = undefined
    }
    loadEditingData()
  }

  const handleSubmit = async () => {
    submitting.value = true

    try {
      const data: CreateTaskRequest | UpdateTaskRequest = {
        title: form.title.trim(),
        description: form.description.trim(),
        tag_id: Number(form.tag_id),
        due_date: new Date(form.due_date + 'T12:00:00Z').toISOString(),
        estimated_minutes: form.estimated_minutes,
      }

      if (editingTask.value) {
        await taskAPI.update(editingTask.value.id, data as UpdateTaskRequest)
        showNotification('success', 'Task updated successfully')
      } else {
        await taskAPI.create(data as CreateTaskRequest)
        showNotification('success', 'Task created successfully')
      }

      await loadTasks()
      showForm.value = false
      editingTask.value = undefined
      resetForm()
    } catch (err) {
      console.error('Error saving task:', err)
      showNotification('error', 'Failed to save task')
    } finally {
      submitting.value = false
    }
  }

  const handleEdit = (task: Task) => {
    editingTask.value = task
    showForm.value = true
    loadEditingData()
  }

  const handleCancel = () => {
    showForm.value = false
    editingTask.value = undefined
    resetForm()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this task? This action cannot be undone.')) {
      return
    }

    try {
      await taskAPI.delete(id)
      showNotification('success', 'Task deleted successfully')
      await loadTasks()
    } catch (err) {
      console.error('Error deleting task:', err)
      showNotification('error', 'Failed to delete task')
    }
  }

  const completeTask = async (id: number) => {
    try {
      await taskAPI.complete(id)
      showNotification('success', 'Task marked as completed')
      await loadTasks()
    } catch (err) {
      console.error('Error completing task:', err)
      showNotification('error', 'Failed to complete task')
    }
  }

  const incompleteTask = async (id: number) => {
    try {
      await taskAPI.incomplete(id)
      showNotification('success', 'Task reopened')
      await loadTasks()
    } catch (err) {
      console.error('Error reopening task:', err)
      showNotification('error', 'Failed to reopen task')
    }
  }

  onMounted(() => {
    loadTasks()
    loadTags()
  })
</script>
