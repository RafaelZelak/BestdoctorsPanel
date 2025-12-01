// src/composables/useFilters.js
import { reactive, computed } from 'vue'

/**
 * Singleton reactive filters store to be shared across components.
 * Combines statusFilter and dateFilter, but keeps each logic isolated.
 */

const state = reactive({
  statusFilter: 'all', // 'all' | 'ai_active' | 'ai_off' | other statuses if any
  dateFilter: {
    start: null, // string 'YYYY-MM-DD' or null
    end: null,   // string 'YYYY-MM-DD' or null
  },
})

function setStatusFilter(next) {
  state.statusFilter = next || 'all'
}

function setDateFilter({ start, end }) {
  state.dateFilter.start = start || null
  state.dateFilter.end = end || null
}

function clearDateFilter() {
  state.dateFilter.start = null
  state.dateFilter.end = null
}

/**
 * Normalize a date string 'YYYY-MM-DD' to Date at start or end of day.
 */
function toDayEdge(dateStr, edge = 'start') {
  if (!dateStr) return null
  const [y, m, d] = dateStr.split('-').map(Number)
  if (edge === 'start') {
    return new Date(y, (m - 1), d, 0, 0, 0, 0)
  }
  return new Date(y, (m - 1), d, 23, 59, 59, 999)
}

/**
 * Check if session matches dateFilter by last_message_at.
 * Accepts ISO/string/number dates. Inclusive range.
 */
function matchDate(session, df) {
  if (!df?.start && !df?.end) return true
  const raw = session?.last_message_at || session?.sessionphone?.last_message_at
  if (!raw) return false

  const msgDate = new Date(raw)
  if (Number.isNaN(msgDate.getTime())) return false

  const start = toDayEdge(df.start, 'start')
  const end = toDayEdge(df.end, 'end')

  if (start && msgDate < start) return false
  if (end && msgDate > end) return false
  return true
}

/**
 * Check if session matches statusFilter (isolated logic).
 * Default assumed statuses:
 *  - 'all': no restriction
 *  - 'ai_active': session.ai_active === true
 *  - 'ai_off': session.ai_active === false
 * Extend if you have other status definitions.
 */
function matchStatus(session, status) {
  if (!status || status === 'all') return true
  if (status === 'ai_active') return !!session?.ai_active
  if (status === 'ai_off') return !session?.ai_active
  // Add custom statuses here if needed
  return true
}

/**
 * Unified filter function (AND). Keeps each filter's logic isolated.
 */
function filterSessions(sessions) {
  if (!Array.isArray(sessions)) return []
  const { statusFilter, dateFilter } = state
  return sessions.filter((s) => {
    return matchStatus(s, statusFilter) && matchDate(s, dateFilter)
  })
}

const hasActiveDateFilter = computed(() => !!(state.dateFilter.start || state.dateFilter.end))

export function useFilters() {
  return {
    // state
    statusFilter: computed(() => state.statusFilter),
    dateFilter: computed(() => state.dateFilter),
    hasActiveDateFilter,

    // setters
    setStatusFilter,
    setDateFilter,
    clearDateFilter,

    // core
    filterSessions,
  }
}
