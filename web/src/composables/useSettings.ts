import { storeToRefs } from 'pinia'
import { computed, readonly } from 'vue'
import { useSettingsStore } from '@/stores/settings'

export function useSettings() {
  const settingsStore = useSettingsStore()
  const { settings } = storeToRefs(settingsStore)

  // Task-related settings
  const taskShowCompleted = computed({
    get: () => settings.value.taskShowCompleted,
    set: value => settingsStore.updateSetting('taskShowCompleted', value),
  })

  const taskShowSuspended = computed({
    get: () => settings.value.taskShowSuspended,
    set: value => settingsStore.updateSetting('taskShowSuspended', value),
  })

  const taskDateFilter = computed({
    get: () => settings.value.taskDateFilter,
    set: value => settingsStore.updateSetting('taskDateFilter', value),
  })

  // Time log settings
  const timeLogShowOnlyActive = computed({
    get: () => settings.value.timeLogShowOnlyActive,
    set: value => settingsStore.updateSetting('timeLogShowOnlyActive', value),
  })

  // UI settings
  const theme = computed({
    get: () => settings.value.theme,
    set: value => settingsStore.updateSetting('theme', value),
  })

  const itemsPerPage = computed({
    get: () => settings.value.itemsPerPage,
    set: value => settingsStore.updateSetting('itemsPerPage', value),
  })

  return {
    // Settings state
    settings: readonly(settings),

    // Task settings
    taskShowCompleted,
    taskShowSuspended,
    taskDateFilter,

    // Time log settings
    timeLogShowOnlyActive,

    // UI settings
    theme,
    itemsPerPage,

    // Methods
    updateSetting: settingsStore.updateSetting,
    loadSettings: settingsStore.loadSettings,
    saveSettings: settingsStore.saveSettings,
  }
}
