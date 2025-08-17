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
          <div class="flex space-x-2">
            <input
              id="start_time"
              v-model="form.start_time"
              type="datetime-local"
              required
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              type="button"
              @click="setCurrentTime('start')"
              class="px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-500 whitespace-nowrap transition-colors"
              title="Set to current time"
            >
              <div class="flex items-center space-x-1">
                <ClockIcon class="h-4 w-4" />
                <span class="hidden sm:inline">Now</span>
              </div>
            </button>
          </div>
        </div>
        
        <div>
          <label for="end_time" class="block text-sm font-medium text-gray-700 mb-2">
            End Time
          </label>
          <div class="flex space-x-2">
            <input
              id="end_time"
              v-model="form.end_time"
              type="datetime-local"
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              type="button"
              @click="setCurrentTime('end')"
              class="px-3 py-2 text-xs font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 focus:outline-none focus:ring-2 focus:ring-blue-500 whitespace-nowrap transition-colors"
              title="Set to current time"
            >
              <div class="flex items-center space-x-1">
                <ClockIcon class="h-4 w-4" />
                <span class="hidden sm:inline">Now</span>
              </div>
            </button>
            <button
              v-if="form.end_time"
              type="button"
              @click="clearEndTime"
              class="px-3 py-2 text-xs font-medium text-gray-600 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors"
              title="Clear end time"
            >
              <XMarkIcon class="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>
      
      <div>
        <label for="tag_id" class="block text-sm font-medium text-gray-700 mb-2">
          Tag *
        </label>
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
import { reactive, watch, computed } from 'vue'
import { ClockIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import type { TimeLog, Tag, CreateTimeLogRequest, UpdateTimeLogRequest } from '@/types'
import { formatDateTimeLocal, formatUTCToLocal, formatLocalToUTC } from '@/utils/date'

const props = defineProps<{
  editingLog?: TimeLog
  submitting: boolean
  availableTags: Tag[]
  lastEndTime?: string | null
}>()

const emit = defineEmits<{
  submit: [data: CreateTimeLogRequest | UpdateTimeLogRequest]
  cancel: []
}>()

const isEditing = computed(() => !!props.editingLog)

const form = reactive<{
  start_time: string
  end_time: string
  tag_id: number | ''
  remarks: string
}>({
  start_time: '',
  end_time: '',
  tag_id: '',
  remarks: '',
})

const resetForm = () => {
  // 如果有上一个 timelog 的结束时间，使用它作为默认的开始时间
  if (props.lastEndTime) {
    form.start_time = formatUTCToLocal(props.lastEndTime)
  } else {
    form.start_time = formatDateTimeLocal()
  }
  form.end_time = ''
  form.tag_id = ''
  form.remarks = ''
}

const setCurrentTime = (field: 'start' | 'end') => {
  const currentTime = formatDateTimeLocal()
  if (field === 'start') {
    form.start_time = currentTime
  } else {
    form.end_time = currentTime
  }
}

const clearEndTime = () => {
  form.end_time = ''
}

const loadEditingData = () => {
  if (props.editingLog) {
    form.start_time = formatUTCToLocal(props.editingLog.start_time)
    form.end_time = props.editingLog.end_time 
      ? formatUTCToLocal(props.editingLog.end_time)
      : ''
    form.tag_id = props.editingLog.tag_id
    form.remarks = props.editingLog.remarks
  } else {
    resetForm()
  }
}

const handleSubmit = () => {
  if (form.tag_id === '') {
    return // Tag is required
  }
  
  const data = {
    start_time: formatLocalToUTC(form.start_time),
    end_time: form.end_time ? formatLocalToUTC(form.end_time) : undefined,
    tag_id: Number(form.tag_id),
    remarks: form.remarks,
  }
  
  emit('submit', data)
}

const handleCancel = () => {
  emit('cancel')
  resetForm()
}

watch(() => props.editingLog, loadEditingData, { immediate: true })
</script>