<template>
  <!--
    Desktop: static right aside (w-80).
    Mobile: persistent bottom sheet — no backdrop, no overlay.
    The sheet collapses to a slim bar when a version is selected so
    the diff content above is fully interactive.
  -->
  <aside
    :class="[
      'bg-neutral-800 border-neutral-700 flex flex-col overflow-hidden',
      'transition-all duration-300 ease-in-out',
      // Desktop: right sidebar
      'md:relative md:w-80 md:border-l md:h-auto md:translate-y-0 md:z-auto md:rounded-none md:max-h-none',
      // Mobile: bottom sheet anchored at the bottom, no backdrop
      'fixed bottom-0 left-0 right-0 z-40 border-t rounded-t-2xl',
      mobileSheetHeight,
    ]"
  >
    <!-- Drag handle (mobile only) -->
    <div class="md:hidden flex justify-center pt-2.5 pb-1 flex-shrink-0">
      <div class="w-10 h-1 rounded-full bg-neutral-600"></div>
    </div>

    <!-- Header -->
    <div class="px-4 py-3 border-b border-neutral-700 flex-shrink-0 flex items-center justify-between">
      <!-- Back to list button — only on mobile when a version is selected -->
      <button
        v-if="selectedVersion"
        @click="$emit('clear-version')"
        class="md:hidden flex items-center gap-1.5 text-sm text-blue-400 hover:text-blue-300 transition active:scale-95"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Versões
      </button>
      <h2 v-else class="font-semibold text-neutral-200 text-sm">Histórico de Versões</h2>

      <button 
        @click="$emit('close')" 
        class="ml-auto text-neutral-400 hover:text-white transition p-1 rounded-lg hover:bg-neutral-700 active:scale-95"
        title="Fechar Histórico"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Version list (hidden on mobile when a version is selected — the diff is shown above) -->
    <div
      :class="[
        'flex-1 overflow-y-auto p-3 space-y-3 relative',
        selectedVersion ? 'hidden md:block' : ''
      ]"
    >
      <div v-if="loading" class="text-center p-4 text-neutral-400 text-sm">
        Carregando histórico...
      </div>
      <div v-else-if="error" class="text-sm text-red-400 bg-red-900/20 p-3 rounded-lg border border-red-500/50">
        {{ error }}
      </div>
      <div v-else-if="historyList.length === 0" class="text-center p-4 text-neutral-500 text-sm">
        Nenhum histórico encontrado.
      </div>
      
      <div v-else class="space-y-3">
        <div 
          v-for="(item, index) in sortedHistory" 
          :key="item.version"
          class="relative pl-4"
        >
          <!-- Timeline vertical line -->
          <div 
            v-if="index !== sortedHistory.length - 1" 
            class="absolute left-1.5 top-6 bottom-[-16px] w-[2px] bg-neutral-700"
          ></div>
          <!-- Timeline dot -->
          <div 
            class="absolute left-0 top-1.5 w-3.5 h-3.5 rounded-full border-2 border-neutral-800"
            :class="getActionColor(item.change_type)"
          ></div>

          <div 
            @click="$emit('select-version', item)"
            :class="[
              'p-3 rounded-lg border cursor-pointer transition active:scale-[0.98]',
              selectedVersion?.version === item.version 
                ? 'bg-blue-900/30 border-blue-500' 
                : 'bg-neutral-900/50 border-neutral-700 hover:border-neutral-500 hover:bg-neutral-800'
            ]"
          >
            <div class="flex justify-between items-start mb-1">
              <span class="font-bold text-sm text-white">v{{ item.version }}</span>
              <span class="text-xs text-neutral-400">{{ formatDate(item.changed_at) }}</span>
            </div>
            
            <div class="flex items-center gap-2 mb-2">
              <span 
                class="px-2 py-0.5 text-[10px] font-bold uppercase rounded"
                :class="getBadgeColor(item.change_type)"
              >
                {{ item.change_type }}
              </span>
              <span v-if="item.edited_by" class="text-[10px] text-neutral-400 font-medium truncate flex items-center gap-1" :title="item.edited_by">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
                {{ item.edited_by }}
              </span>
            </div>

            <p v-if="item.snapshot" class="text-xs text-neutral-400 line-clamp-2">
              {{ item.snapshot }}
            </p>
            <p v-else class="text-xs text-neutral-500 italic">
              Nenhum snapshot
            </p>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  historyList: {
    type: Array,
    required: true
  },
  loading: Boolean,
  error: String,
  selectedVersion: Object,
  isOpen: {
    type: Boolean,
    default: true
  }
})

defineEmits(['close', 'select-version', 'clear-version'])

// Collapsed to a slim bar when a diff is selected (mobile only).
// On desktop the height is always auto.
const mobileSheetHeight = computed(() => {
  if (props.selectedVersion) {
    return 'h-[52px] md:h-auto'
  }
  return 'h-[62vh] md:h-auto'
})

const sortedHistory = computed(() => {
  return [...props.historyList].sort((historyA, historyB) => historyB.version - historyA.version)
})

function formatDate(isoString) {
  if (!isoString) return ''
  const parsedDate = new Date(isoString)
  return parsedDate.toLocaleString('pt-BR', { 
    day: '2-digit', 
    month: '2-digit', 
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getActionColor(type) {
  if (type === 'created') return 'bg-green-500'
  if (type === 'updated') return 'bg-blue-500'
  if (type === 'deleted') return 'bg-red-500'
  return 'bg-neutral-500'
}

function getBadgeColor(type) {
  if (type === 'created') return 'bg-green-500/20 text-green-400 border border-green-500/30'
  if (type === 'updated') return 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
  if (type === 'deleted') return 'bg-red-500/20 text-red-400 border border-red-500/30'
  return 'bg-neutral-700 text-neutral-300'
}
</script>
