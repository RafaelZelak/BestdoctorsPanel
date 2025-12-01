<template>
  <div class="relative" ref="triggerWrap">
    <!-- Trigger Button -->
    <button
      type="button"
      @click="toggle"
      class="rounded-xl border border-neutral-700 bg-neutral-800 px-3 py-2 hover:bg-neutral-700/60 outline-none focus:ring-2 focus:ring-blue-500"
      aria-haspopup="dialog"
      :aria-expanded="open ? 'true' : 'false'"
      :aria-controls="menuIdComputed"
      title="Filters"
    >
      <!-- Filter icon -->
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
        <path d="M3 5h18v2H3V5zm4 6h10v2H7v-2zm4 6h2v2h-2v-2z"/>
      </svg>
      <span class="sr-only">Abrir Filtros</span>
    </button>

    <!-- Popover -->
    <transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        :id="menuIdComputed"
        role="dialog"
        aria-modal="true"
        class="absolute right-0 mt-2 w-80 rounded-xl border border-neutral-700 bg-neutral-900/80 backdrop-blur p-3 shadow-xl z-50"
      >
        <div class="flex items-center justify-between pb-2 border-b border-neutral-800">
          <h3 class="text-sm font-semibold">Filtros</h3>
          <button
            type="button"
            @click="close"
            class="p-1 rounded-lg hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-label="Close filters"
            title="Close"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
              <path d="M18.3 5.71L12 12.01l-6.3-6.3-1.4 1.41 6.29 6.29-6.3 6.3 1.41 1.41 6.3-6.3 6.29 6.3 1.41-1.41-6.3-6.3 6.3-6.29z"/>
            </svg>
          </button>
        </div>

        <!-- Status Filter (TRIPLE TOGGLE) -->
        <div class="mt-3">
          <label class="block text-xs text-neutral-300 mb-2">Status</label>

          <div
            role="group"
            aria-label="Status filter"
            class="inline-flex w-full rounded-lg overflow-hidden border border-neutral-700"
          >
            <button
              type="button"
              :aria-pressed="currentStatus === 'all' ? 'true' : 'false'"
              @click="setStatus('all')"
              class="flex-1 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-500"
              :class="btnClass(currentStatus === 'all')"
            >
              Todos
            </button>
            <button
              type="button"
              :aria-pressed="currentStatus === 'ai_active' ? 'true' : 'false'"
              @click="setStatus('ai_active')"
              class="flex-1 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-500 border-l border-neutral-700"
              :class="btnClass(currentStatus === 'ai_active')"
            >
              IA Ativa
            </button>
            <button
              type="button"
              :aria-pressed="currentStatus === 'ai_off' ? 'true' : 'false'"
              @click="setStatus('ai_off')"
              class="flex-1 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-500 border-l border-neutral-700"
              :class="btnClass(currentStatus === 'ai_off')"
            >
              IA Inativa
            </button>
          </div>
        </div>

        <!-- Date Filter -->
        <div class="mt-4 space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-xs text-neutral-300">Período</span>
            <button
              type="button"
              class="text-xs underline hover:no-underline text-neutral-300 hover:text-neutral-100 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded"
              @click="resetDate"
              :disabled="!hasActiveDateFilter"
            >
              Reset
            </button>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="dateFrom" class="block text-xs text-neutral-400">De</label>
              <input
                id="dateFrom"
                type="date"
                class="w-full rounded-lg bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none focus:ring-2 focus:ring-blue-500"
                :max="end || undefined"
                :value="start || ''"
                @input="onDateChange($event.target.value, end)"
              />
            </div>
            <div>
              <label for="dateTo" class="block text-xs text-neutral-400">Até</label>
              <input
                id="dateTo"
                type="date"
                class="w-full rounded-lg bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none focus:ring-2 focus:ring-blue-500"
                :min="start || undefined"
                :value="end || ''"
                @input="onDateChange(start, $event.target.value)"
              />
            </div>
          </div>

          <p v-if="hasActiveDateFilter" class="text-[11px] text-neutral-400">
            Aplicado de filtro de <span class="text-neutral-200 font-medium">{{ start || new Date().toISOString().slice(0,10)}}</span>
            até <span class="text-neutral-200 font-medium">{{ end || new Date().toISOString().slice(0,10) }}</span>
          </p>
          <p v-else class="text-[11px] text-neutral-500">Sem período selecionado</p>
        </div>

        <!-- Footer -->
        <div class="mt-4 pt-3 border-t border-neutral-800 flex items-center justify-end gap-2">
          <button
            type="button"
            @click="emitClear"
            class="px-3 py-2 text-sm rounded-lg border border-neutral-700 hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Reset filters
          </button>
          <button
            type="button"
            @click="close"
            class="px-3 py-2 text-sm rounded-lg border border-neutral-700 hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Fechar
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useFilters } from '@/composables/useFilters'

/**
 * v-model compatibility for statusFilter + internal centralized store
 */
const props = defineProps({
  statusFilter: { type: String, default: undefined },
  menuId: { type: String, default: '' },
})
const emit = defineEmits(['update:statusFilter', 'clear'])

const menuIdComputed = props.menuId || `filters-popover-${Math.random().toString(36).slice(2, 8)}`
const open = ref(false)
const triggerWrap = ref(null)

// Centralized filters store
const {
  statusFilter: statusFromStore,
  setStatusFilter,
  dateFilter,
  hasActiveDateFilter,
  setDateFilter,
  clearDateFilter,
} = useFilters()

// Effective current status: prefer v-model value if provided, fallback to store
const currentStatus = computed(() => (props.statusFilter ?? statusFromStore.value))

function setStatus(val) {
  emit('update:statusFilter', val)
  setStatusFilter(val) // keep store in sync
}

const start = computed(() => dateFilter.value.start)
const end = computed(() => dateFilter.value.end)

function toggle() { open.value = !open.value }
function close() { open.value = false }

function onClickOutside(e) {
  if (!triggerWrap.value || !open.value) return
  if (!triggerWrap.value.contains(e.target)) open.value = false
}

function onDateChange(startVal, endVal) {
  const s = startVal || null
  const e = endVal || null
  if (s && e && s > e) {
    setDateFilter({ start: e, end: s })
    return
  }
  setDateFilter({ start: s, end: e })
}

function resetDate() {
  clearDateFilter()
}

function emitClear() {
  clearDateFilter()
  setStatus('all')
  emit('clear')
}

function btnClass(isActive) {
  return isActive
    ? 'bg-blue-600 text-white'
    : 'bg-neutral-800 text-neutral-200 hover:bg-neutral-700/60'
}

onMounted(() => document.addEventListener('click', onClickOutside))
onBeforeUnmount(() => document.removeEventListener('click', onClickOutside))
</script>

<style scoped>
/* Styles kept minimal; Tailwind handles visuals */
</style>
