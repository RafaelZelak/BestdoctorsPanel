<template>
  <!--
    Desktop: static right aside (w-80).
    Mobile: bottom sheet that slides up from the bottom.

    When a version is selected and the diff is being shown, this panel
    collapses on mobile to give full screen to the diff viewer.
  -->

  <!-- Mobile backdrop -->
  <Transition name="fade">
    <div
      v-if="isEffectivelyOpen"
      class="fixed inset-0 bg-black/50 md:hidden z-30"
      @click="$emit('close')"
    />
  </Transition>

  <aside
    :class="[
      'bg-neutral-800 border-neutral-700 flex flex-col overflow-hidden',
      'transition-transform duration-300 ease-in-out',
      // Desktop: right sidebar, always visible when mounted
      'md:relative md:w-80 md:border-l md:translate-y-0 md:z-auto md:h-auto',
      // Mobile: bottom sheet
      'fixed bottom-0 left-0 right-0 z-40 border-t rounded-t-2xl',
      mobileSheetHeight,
      isEffectivelyOpen ? 'translate-y-0' : 'translate-y-full md:translate-y-0'
    ]"
  >
    <!-- Drag handle (mobile only) -->
    <div class="md:hidden flex justify-center pt-3 pb-1 flex-shrink-0 cursor-pointer" @click="$emit('close')">
      <div class="w-10 h-1 rounded-full bg-neutral-600"></div>
    </div>

    <!-- Header -->
    <div class="px-4 py-3 border-b border-neutral-700 flex-shrink-0 flex items-center justify-between">
      <!-- Back button on mobile when diff is selected -->
      <button
        v-if="selectedVersion && hasVersionSelected"
        @click="$emit('clear-version')"
        class="md:hidden flex items-center gap-1.5 text-sm text-neutral-400 hover:text-white transition active:scale-95"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Versões
      </button>
      <h2 v-else class="font-semibold text-neutral-200">Histórico de Versões</h2>
      <button 
        @click="$emit('close')" 
        class="text-neutral-400 hover:text-white transition"
        title="Fechar Histórico"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-3 space-y-3 relative">
      <div v-if="loading" class="text-center p-4 text-neutral-400">
        Carregando histórico...
      </div>
      <div v-else-if="error" class="text-sm text-red-400 bg-red-900/20 p-3 rounded-lg border border-red-500/50">
        {{ error }}
      </div>
      <div v-else-if="historyList.length === 0" class="text-center p-4 text-neutral-500">
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

// On desktop the panel is always "open" (it's mounted via v-if in parent).
// On mobile we need to know if it should slide up.
// We treat it as open when the `isOpen` prop is true OR on non-mobile (md+).
const isEffectivelyOpen = computed(() => props.isOpen)

// When a version is selected, shrink the sheet to a smaller handle so the
// diff in PromptViewer can take the remaining visible area.
// On desktop this class has no effect since we override with md:*.
const hasVersionSelected = computed(() => !!props.selectedVersion)

const mobileSheetHeight = computed(() => {
  if (hasVersionSelected.value) {
    // Collapsed — just show the header with "← Versões" button
    return 'h-[56px] md:h-auto'
  }
  // Full-ish bottom sheet showing the version list
  return 'h-[60vh] md:h-auto'
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

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
