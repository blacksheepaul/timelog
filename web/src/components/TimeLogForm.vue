<template>
  <div class="bg-white shadow rounded-lg p-6">
    <h2 class="text-lg font-medium text-gray-900 mb-6">
      {{ isEditing ? 'Edit Time Log' : 'Create New Time Log' }}
    </h2>
    
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label for="start_time" class="block text-sm font-medium text-gray-700 mb-2">
            Start Time *
          </label>
          <input
            id="start_time"
            v-model="form.start_time"
            type="datetime-local"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
        
        <div>
          <label for="end_time" class="block text-sm font-medium text-gray-700 mb-2">
            End Time
          </label>
          <input
            id="end_time"
            v-model="form.end_time"
            type="datetime-local"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>
      </div>
      
      <div>
        <label for="tags" class="block text-sm font-medium text-gray-700 mb-2">
          Tags
        </label>
        <input
          id="tags"
          v-model="form.tags"
          type="text"
          placeholder="e.g., work, meeting, development"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>
      
      <div>
        <label for="remarks" class="block text-sm font-medium text-gray-700 mb-2">
          Remarks
        </label>
        <textarea
          id="remarks"
          v-model="form.remarks"
          rows="3"
          placeholder="Add any notes or description..."
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
        ></textarea>
      </div>
      
      <div class="flex justify-end space-x-4">
        <button
          v-if="isEditing"
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
          {{ submitting ? 'Saving...' : (isEditing ? 'Update' : 'Create') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { TimeLog, CreateTimeLogRequest, UpdateTimeLogRequest } from '@/types'
import { formatDateTimeLocal } from '@/utils/date'

const props = defineProps<{
  editingLog?: TimeLog
  submitting: boolean
}>()

const emit = defineEmits<{
  submit: [data: CreateTimeLogRequest | UpdateTimeLogRequest]
  cancel: []
}>()

const isEditing = computed(() => !!props.editingLog)

const form = reactive<{
  start_time: string
  end_time: string
  tags: string
  remarks: string
}>({
  start_time: '',
  end_time: '',
  tags: '',
  remarks: '',
})

const resetForm = () => {
  form.start_time = formatDateTimeLocal()
  form.end_time = ''
  form.tags = ''
  form.remarks = ''
}

const loadEditingData = () => {
  if (props.editingLog) {
    form.start_time = new Date(props.editingLog.start_time).toISOString().slice(0, 16)
    form.end_time = props.editingLog.end_time 
      ? new Date(props.editingLog.end_time).toISOString().slice(0, 16) 
      : ''
    form.tags = props.editingLog.tags
    form.remarks = props.editingLog.remarks
  } else {
    resetForm()
  }
}

const handleSubmit = () => {
  const data = {
    start_time: new Date(form.start_time).toISOString(),
    end_time: form.end_time ? new Date(form.end_time).toISOString() : undefined,
    tags: form.tags,
    remarks: form.remarks,
  }
  
  emit('submit', data)
}

const handleCancel = () => {
  emit('cancel')
  resetForm()
}

watch(() => props.editingLog, loadEditingData, { immediate: true })

import { computed } from 'vue'
</script>