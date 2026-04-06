<template>
  <!-- Mobile backdrop overlay -->
  <Transition name="fade">
    <div
      v-if="isMobileOpen"
      class="fixed inset-0 bg-black/60 z-40 md:hidden"
      @click="$emit('close')"
    />
  </Transition>

  <!--
    Desktop: static sidebar (w-72, relative).
    Mobile: fixed full-height drawer sliding from the left (translate-x).
    The 'translate-x-[-100%]' / 'translate-x-0' toggle is driven by isMobileOpen.
  -->
  <aside
    :class="[
      'bg-neutral-800 border-r border-neutral-700 flex flex-col overflow-hidden',
      'transition-transform duration-300 ease-in-out',
      // Desktop layout — always visible, not positioned
      'md:relative md:translate-x-0 md:w-72 md:z-auto',
      // Mobile layout — fixed drawer
      'fixed top-0 left-0 h-full w-80 z-50',
      isMobileOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'
    ]"
  >
    <div class="p-4 border-b border-neutral-700 flex-shrink-0 flex items-center justify-between">
      <h2 class="font-semibold text-neutral-200">Arquivos .md</h2>
      <div class="flex items-center gap-1">
        <button
          @click="$emit('create-folder')"
          title="Criar nova pasta"
          class="w-7 h-7 flex items-center justify-center rounded-lg bg-blue-600/80 hover:bg-blue-500/80 active:scale-95 text-white transition"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10.5v6m3-3H9m4.06-7.19-2.12-2.12a1.5 1.5 0 00-1.061-.44H4.5A2.25 2.25 0 002.25 6v12a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9a2.25 2.25 0 00-2.25-2.25h-5.379a1.5 1.5 0 01-1.06-.44z"/>
          </svg>
        </button>
        <button
          @click="$emit('create-new')"
          title="Criar novo prompt"
          class="w-7 h-7 flex items-center justify-center rounded-lg bg-green-600/80 hover:bg-green-500/80 active:scale-95 text-white transition"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
          </svg>
        </button>

        <!-- Close button — only visible on mobile -->
        <button
          @click="$emit('close')"
          class="md:hidden ml-1 w-7 h-7 flex items-center justify-center rounded-lg text-neutral-400 hover:text-white hover:bg-neutral-700 active:scale-95 transition"
          title="Fechar"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
    
    <div 
      ref="dropRootEl"
      class="flex-1 overflow-y-auto p-2 space-y-1 relative"
      @dragenter.prevent
      @dragover.prevent
      @drop.stop="handleRootDrop"
    >
      <div v-if="loading || loadingFolders" class="p-4 text-center text-neutral-500">
        Carregando...
      </div>
      <div v-else-if="error || folderError" class="p-4 text-sm text-red-400 bg-red-900/20 rounded">
        {{ error || folderError }}
      </div>
      <template v-else>
        
        <!-- Pinned Prompts (base_prompt.md, router.md) -->
        <div class="pb-2 mb-2 border-b border-neutral-700/50 space-y-1">
          <button
            v-for="prompt in pinnedPrompts"
            :key="prompt.name"
            @click="selectAndClose(prompt.name)"
            :class="[
              'w-full text-left px-3 py-2.5 rounded-lg text-sm transition font-mono flex items-center gap-2 border active:scale-[0.98]',
              getPromptClass(prompt.name)
            ]"
          >
            <svg class="w-4 h-4 text-yellow-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
            <span class="truncate flex-1">{{ prompt.name }}</span>
            <span class="flex items-center gap-1 flex-shrink-0">
              <span
                v-for="viewer in (allPresence[prompt.name] || [])"
                :key="viewer.username"
                :title="viewer.status === 'editing' ? `${viewer.username} está editando` : `${viewer.username} está visualizando`"
                :class="[
                  'inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs font-medium leading-none',
                  viewer.status === 'editing'
                    ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40'
                    : 'bg-blue-500/15 text-blue-300 border border-blue-500/30'
                ]"
              >
                <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="viewer.status === 'editing' ? 'bg-amber-400' : 'bg-blue-400'"></span>
                {{ viewer.username }}
              </span>
            </span>
          </button>
        </div>

        <!-- Folders -->
        <PromptFolderItem
          v-for="folder in folders"
          :key="folder.id || folder.ID"
          :folder="folder"
          :prompts="getPromptsForFolder(folder.id || folder.ID)"
          :selectedPromptName="selectedPromptName"
          :allPresence="allPresence"
          @select-prompt="selectAndClose($event)"
          @update="$emit('update-folder', $event.id, $event.patch)"
          @delete="$emit('delete-folder', $event)"
          @drop-into-folder="handleDropIntoFolder"
        />

        <!-- Unassigned Prompts -->
        <div
          ref="unassignedDropEl"
          class="pt-2 mt-2 border-t border-neutral-700/50 space-y-1"
        >
          <button
            v-for="prompt in unassignedPrompts"
            :key="prompt.name"
            draggable="true"
            @dragstart="handleUnassignedPromptDragStart($event, prompt)"
            @touchstart.passive="handlePromptTouchStart($event, prompt)"
            @touchmove="handlePromptTouchMove"
            @touchend.passive="handlePromptTouchEnd"
            @click="selectAndClose(prompt.name)"
            :class="[
              'w-full text-left px-3 py-2.5 rounded-lg text-sm transition font-mono flex items-center gap-2 border active:scale-[0.98]',
              getPromptClass(prompt.name)
            ]"
          >
            <svg class="w-4 h-4 opacity-70 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="truncate flex-1">{{ prompt.name }}</span>
            <span class="flex items-center gap-1 flex-shrink-0">
              <span
                v-for="viewer in (allPresence[prompt.name] || [])"
                :key="viewer.username"
                :title="viewer.status === 'editing' ? `${viewer.username} está editando` : `${viewer.username} está visualizando`"
                :class="[
                  'inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-xs font-medium leading-none',
                  viewer.status === 'editing'
                    ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40'
                    : 'bg-blue-500/15 text-blue-300 border border-blue-500/30'
                ]"
              >
                <span class="w-1.5 h-1.5 rounded-full flex-shrink-0" :class="viewer.status === 'editing' ? 'bg-amber-400' : 'bg-blue-400'"></span>
                {{ viewer.username }}
              </span>
            </span>
          </button>
          
          <!-- Drop Area inside unassigned -->
          <div 
            v-if="unassignedPrompts.length === 0" 
            class="p-4 text-center border-2 border-dashed border-neutral-700/50 rounded-lg text-neutral-500 text-xs italic"
          >
            Arraste prompts para cá para removê-los de pastas
          </div>
        </div>
      </template>
    </div>
  </aside>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import PromptFolderItem from './PromptFolderItem.vue'
import { useTouchDragSource, useTouchDropTarget } from '@/composables/useTouchDrag'

const props = defineProps({
  prompts: {
    type: Array,
    default: () => []
  },
  loading: Boolean,
  error: String,
  folders: {
    type: Array,
    default: () => []
  },
  loadingFolders: Boolean,
  folderError: String,
  selectedPromptName: String,
  allPresence: {
    type: Object,
    default: () => ({})
  },
  isMobileOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits([
  'create-new',
  'create-folder',
  'select-prompt',
  'update-folder',
  'delete-folder',
  'assign-prompt',
  'reorder-folder',
  'close'
])

const dropRootEl = ref(null)
const unassignedDropEl = ref(null)

// Closes drawer on mobile after selecting a prompt
function selectAndClose(promptName) {
  emit('select-prompt', promptName)
  emit('close')
}

const pinnedPrompts = computed(() => {
  return props.prompts.filter(p => p.name === 'base_prompt.md' || p.name === 'router.md')
})

const unassignedPrompts = computed(() => {
  return props.prompts.filter(p => 
    p.name !== 'base_prompt.md' &&
    p.name !== 'router.md' &&
    !p.folder_id && !p.FolderID
  )
})

function getPromptsForFolder(folderId) {
  return props.prompts.filter(p => p.folder_id === folderId || p.FolderID === folderId)
}

function getPromptClass(promptName) {
  if (props.selectedPromptName === promptName) {
    return 'bg-blue-600 text-white border-transparent'
  }
  if (promptName === 'base_prompt.md' || promptName === 'router.md') {
    return 'bg-yellow-900/20 text-yellow-500 hover:bg-yellow-900/40 border-yellow-700/50'
  }
  return 'text-neutral-300 hover:bg-neutral-700 border-transparent'
}

// ── Desktop drag handlers ──────────────────────────────────────────────────
function handleUnassignedPromptDragStart(e, prompt) {
  e.dataTransfer.setData('type', 'prompt')
  e.dataTransfer.setData('promptName', prompt.name)
}

function handleRootDrop(e) {
  const type = e.dataTransfer.getData('type')
  if (type === 'prompt') {
    const promptName = e.dataTransfer.getData('promptName')
    if (promptName) {
      emit('assign-prompt', promptName, null)
    }
  }
}

function handleDropIntoFolder({ folderId, targetFolderId, promptName, draggedFolderId }) {
  if (promptName) {
    emit('assign-prompt', promptName, folderId || targetFolderId)
  } else if (draggedFolderId) {
    emit('reorder-folder', draggedFolderId, folderId || targetFolderId)
  }
}

// ── Touch drag handlers for unassigned prompts ─────────────────────────────
let activeTouchPrompt = null

const { handleTouchStart, handleTouchMove, handleTouchEnd } = useTouchDragSource(
  () => { activeTouchPrompt = null },
  () => { activeTouchPrompt = null }
)

function handlePromptTouchStart(event, prompt) {
  activeTouchPrompt = prompt
  handleTouchStart(event, event.currentTarget, 'prompt', { promptName: prompt.name })
}

function handlePromptTouchMove(event) {
  handleTouchMove(event)
}

function handlePromptTouchEnd(event) {
  handleTouchEnd(event)
}

// Touch drop zone for the root (unassigning a prompt)
const { bindDropZone, unbindDropZone } = useTouchDropTarget(
  (detail) => {
    if (detail.type === 'prompt' && detail.promptName) {
      emit('assign-prompt', detail.promptName, null)
    }
  }
)

onMounted(() => {
  if (unassignedDropEl.value) bindDropZone(unassignedDropEl.value)
})

onUnmounted(() => {
  if (unassignedDropEl.value) unbindDropZone(unassignedDropEl.value)
})
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
