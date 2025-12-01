<!-- src/components/SessionList.vue -->
<template>
  <aside class="border-r border-neutral-800 flex flex-col h-full">
    <header class="p-3 border-b border-neutral-800">
      <div class="flex items-center justify-between">
        <h2 class="text-xl font-semibold">BestDoctors • Chat</h2>

        <!-- Open modal -->
        <button @click="openModal">Relatórios</button>
        <ReportView :isVisible="isModalVisible" @close="closeModal" />
      </div>

      <div class="mt-1 flex gap-2 items-center">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Digite um nome ou telefone"
          class="w-full rounded-xl bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none"
          aria-label="Search by name or phone"
        />

        <!-- Keep FiltersMenu API as you have it; it controls status/date internally via useFilters -->
        <FiltersMenu
          v-model:statusFilter="statusFilterProxy"
          menu-id="session-filters"
          @clear="handleClear"
        />
      </div>
    </header>

    <div class="flex-1 overflow-y-auto overflow-x-hidden">
      <ul>
        <SessionListItem
          v-for="s in filteredItems"
          :key="s.session_id"
          :session="s"
          :selected="selectedId === s.session_id"
          @select="$emit('select', s)"
        />
      </ul>

      <div v-if="!filteredItems.length && !loading" class="p-4 text-sm text-neutral-400">
        Nenhuma sessão encontrada.
      </div>
      <div v-if="loading" class="p-4 text-sm text-neutral-400">Carregando Sessões...</div>
      <div v-if="error" class="p-4 text-sm text-red-400">{{ error }}</div>
    </div>
  </aside>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watch, computed } from 'vue'
import { useSessions } from '@/composables/useSessions'
import { useFilters } from '@/composables/useFilters'
import SessionListItem from './SessionListItem.vue'
import FiltersMenu from './FiltersMenu.vue'
import ReportView from './ReportView.vue'

// ----- Modal visibility (single source of truth)
const isModalVisible = ref(false)
const openModal = () => { isModalVisible.value = true }
const closeModal = () => { isModalVisible.value = false }

// ----- Props & Emits
defineProps({ selectedId: { type: String, default: '' } })
defineEmits(['select'])

// ----- Sessions & polling
const { loading, error, sessions, loadSessions, isPolling, startPolling, stopPolling } = useSessions()

/**
 * Filters (status/date) are centralized in useFilters().
 * We compose final list = filterSessions(sessions) then apply local search (name/phone).
 */
const { filterSessions, statusFilter, setStatusFilter } = useFilters()

// ----- Local search state (name/phone)
const searchQuery = ref('')

// ----- Derived list, already filtered by status+date (AND)
const baseFiltered = computed(() => filterSessions(sessions.value))

// ----- Helpers for normalization
function normalizeDigits(str) { return String(str || '').replace(/\D/g, '') }
function normalizeText(str) { return String(str || '').toLowerCase() }

// ----- Apply local search (lead_name OR phone digits) on top of baseFiltered
const filteredItems = computed(() => {
  const list = baseFiltered.value || []
  const raw = (searchQuery.value || '').trim()
  if (!raw) return list

  const queryDigits = normalizeDigits(raw)
  const queryText = normalizeText(raw)

  return list.filter((s) => {
    const name = normalizeText(s?.lead_name ?? '')
    const phoneDigits = normalizeDigits(s?.phone ?? '')
    const matchName = queryText ? name.includes(queryText) : false
    const matchPhone = queryDigits ? phoneDigits.includes(queryDigits) : false
    return matchName || matchPhone
  })
})

// ----- v-model compatibility for FiltersMenu (statusFilter)
const statusFilterProxy = computed({
  get: () => statusFilter.value,
  set: (val) => setStatusFilter(val),
})

function handleClear() {
  // optional: clearing filters also clears local text search
  searchQuery.value = ''
}

const auto = ref(true)

onMounted(async () => {
  await loadSessions()
  if (auto.value) startPolling()
})

onUnmounted(() => { stopPolling() })

watch(auto, (val) => (val ? startPolling() : stopPolling()))
watch(isPolling, (val) => { if (val !== auto.value) auto.value = val })
</script>
