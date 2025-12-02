// src/composables/useResponsive.js
import { onMounted, onUnmounted, ref } from 'vue'

/**
 * Composable for responsive design detection
 * Detects if the current viewport is mobile or desktop
 * Uses Tailwind's 'md' breakpoint (768px) as threshold
 */
export function useResponsive() {
  const isMobile = ref(false)
  let mediaQuery = null

  function updateIsMobile(e) {
    isMobile.value = e.matches
  }

  onMounted(() => {
    // Create media query for mobile detection (max-width: 767px)
    // This matches Tailwind's 'md' breakpoint
    mediaQuery = window.matchMedia('(max-width: 767px)')
    
    // Set initial value
    isMobile.value = mediaQuery.matches

    // Listen for changes
    mediaQuery.addEventListener('change', updateIsMobile)
  })

  onUnmounted(() => {
    // Cleanup listener
    if (mediaQuery) {
      mediaQuery.removeEventListener('change', updateIsMobile)
    }
  })

  return {
    isMobile
  }
}
