import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export interface SettingsState {
  // Task view settings
  taskShowCompleted: boolean
  taskShowSuspended: boolean
  taskDateFilter: string

  // Time log view settings
  timeLogShowOnlyActive: boolean

  // UI preferences
  theme: 'light' | 'dark' | 'auto'

  // Display preferences
  itemsPerPage: number

  // Other settings can be added here
}

const SETTINGS_STORAGE_KEY = 'timelog-settings'

export const useSettingsStore = defineStore('settings', () => {
  // Default settings
  const settings = ref<SettingsState>({
    taskShowCompleted: false,
    taskShowSuspended: false,
    taskDateFilter: '',
    timeLogShowOnlyActive: true,
    theme: 'auto',
    itemsPerPage: 20
  })

  // Load settings from localStorage
  const loadSettings = () => {
    try {
      const stored = localStorage.getItem(SETTINGS_STORAGE_KEY)
      if (stored) {
        const parsed = JSON.parse(stored)
        // Merge with defaults to handle missing keys
        settings.value = { ...settings.value, ...parsed }
      }
    } catch (error) {
      console.error('Failed to load settings from localStorage:', error)
    }
  }

  // Save settings to localStorage
  const saveSettings = () => {
    try {
      localStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(settings.value))
    } catch (error) {
      console.error('Failed to save settings to localStorage:', error)
    }
  }

  // Update a specific setting
  const updateSetting = <K extends keyof SettingsState>(
    key: K,
    value: SettingsState[K]
  ) => {
    settings.value[key] = value
  }

  // Watch for changes and auto-save
  watch(
    settings,
    () => {
      saveSettings()
    },
    { deep: true }
  )

  // Initialize on store creation
  loadSettings()

  return {
    settings,
    updateSetting,
    loadSettings,
    saveSettings
  }
})