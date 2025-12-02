// src/composables/useMobileNavigation.js
import { computed, ref } from 'vue'

/**
 * Composable for managing mobile navigation state
 * Controls which view is visible on mobile devices
 */
export function useMobileNavigation() {
  // Current view: 'list' or 'chat'
  const currentView = ref('list')

  // Actions to change view
  const showList = () => {
    currentView.value = 'list'
  }

  const showChat = () => {
    currentView.value = 'chat'
  }

  // Computed properties for checking current view
  const isListVisible = computed(() => currentView.value === 'list')
  const isChatVisible = computed(() => currentView.value === 'chat')

  return {
    currentView,
    showList,
    showChat,
    isListVisible,
    isChatVisible
  }
}
