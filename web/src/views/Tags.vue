<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h1 class="text-2xl font-bold text-gray-900">Tag Management</h1>
      <button
        @click="toggleForm"
        class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <PlusIcon class="h-5 w-5 mr-2" />
        New Tag
      </button>
    </div>

    <!-- Tag创建/编辑表单 -->
    <div v-if="showForm" class="bg-white shadow rounded-lg p-6">
      <h2 class="text-lg font-medium text-gray-900 mb-6">
        {{ isEditing ? 'Edit Tag' : 'Create New Tag' }}
      </h2>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 mb-2"> Name * </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              maxlength="50"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Tag name"
            />
          </div>

          <div>
            <label for="color" class="block text-sm font-medium text-gray-700 mb-2">
              Color *
            </label>
            <div class="flex items-center space-x-2">
              <input
                id="color"
                v-model="form.color"
                type="color"
                required
                class="h-10 w-16 border border-gray-300 rounded-md"
              />
              <input
                v-model="form.color"
                type="text"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="#000000"
              />
            </div>
          </div>
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            maxlength="200"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="Describe what this tag is used for..."
          ></textarea>
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

    <!-- Tags列表 -->
    <div class="bg-white shadow rounded-lg">
      <div class="px-6 py-4 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">All Tags</h2>
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

      <div v-else-if="tags.length === 0" class="p-6 text-center text-gray-500">
        No tags found. Create your first one!
      </div>

      <div v-else class="divide-y divide-gray-200">
        <div v-for="tag in tags" :key="tag.id" class="p-6 hover:bg-gray-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-4">
              <span
                class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium text-white"
                :style="{ backgroundColor: tag.color }"
              >
                {{ tag.name }}
              </span>
              <div>
                <p class="text-sm text-gray-900">{{ tag.description || 'No description' }}</p>
                <p class="text-xs text-gray-500">Created: {{ formatDateTime(tag.created_at) }}</p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <button
                @click="handleEdit(tag)"
                class="text-blue-600 hover:text-blue-900 text-sm font-medium"
              >
                Edit
              </button>
              <button
                @click="handleDelete(tag.id)"
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
  import { tagAPI } from '@/api'
  import { formatDateTime } from '@/utils/date'
  import type { Tag } from '@/types'

  // 注入全局通知功能
  const showNotification = inject('showNotification') as (
    type: 'success' | 'error',
    message: string
  ) => void

  const tags = ref<Tag[]>([])
  const loading = ref(false)
  const submitting = ref(false)
  const error = ref<string | null>(null)
  const showForm = ref(false)
  const editingTag = ref<Tag | undefined>()

  const isEditing = computed(() => !!editingTag.value)

  const form = reactive({
    name: '',
    color: '#3B82F6',
    description: '',
  })

  const resetForm = () => {
    form.name = ''
    form.color = '#3B82F6'
    form.description = ''
  }

  const loadEditingData = () => {
    if (editingTag.value) {
      form.name = editingTag.value.name
      form.color = editingTag.value.color
      form.description = editingTag.value.description
    } else {
      resetForm()
    }
  }

  const loadTags = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await tagAPI.getAll()
      tags.value = response.data || []
    } catch (err) {
      error.value = 'Failed to load tags'
      console.error('Error loading tags:', err)
      showNotification('error', 'Failed to load tags')
    } finally {
      loading.value = false
    }
  }

  const toggleForm = () => {
    showForm.value = !showForm.value
    if (!showForm.value) {
      editingTag.value = undefined
    }
    loadEditingData()
  }

  const handleSubmit = async () => {
    submitting.value = true

    try {
      const data = {
        name: form.name.trim(),
        color: form.color,
        description: form.description.trim(),
      }

      if (editingTag.value) {
        await tagAPI.update(editingTag.value.id, data)
        showNotification('success', 'Tag updated successfully')
      } else {
        await tagAPI.create(data)
        showNotification('success', 'Tag created successfully')
      }

      await loadTags()
      showForm.value = false
      editingTag.value = undefined
      resetForm()
    } catch (err) {
      console.error('Error saving tag:', err)
      showNotification('error', 'Failed to save tag')
    } finally {
      submitting.value = false
    }
  }

  const handleEdit = (tag: Tag) => {
    editingTag.value = tag
    showForm.value = true
    loadEditingData()
  }

  const handleCancel = () => {
    showForm.value = false
    editingTag.value = undefined
    resetForm()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this tag? This action cannot be undone.')) {
      return
    }

    try {
      await tagAPI.delete(id)
      showNotification('success', 'Tag deleted successfully')
      await loadTags()
    } catch (err) {
      console.error('Error deleting tag:', err)
      showNotification('error', 'Failed to delete tag')
    }
  }

  onMounted(() => {
    loadTags()
  })
</script>
