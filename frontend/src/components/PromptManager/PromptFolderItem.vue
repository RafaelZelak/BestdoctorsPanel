<template>
  <div
    class="flex flex-col mb-1 border border-transparent rounded-lg transition-colors duration-200"
    :class="{'bg-neutral-800 border-neutral-700': isDragOver}"
    @dragenter.prevent="handleDragEnter"
    @dragleave.prevent="handleDragLeave"
    @dragover.prevent
    @drop.stop="handleDrop"
  >
    <!-- Folder Row -->
    <div
      class="flex items-center gap-2 px-2 py-2 rounded-lg hover:bg-neutral-700/50 group cursor-pointer"
      :class="{'opacity-50': isDragging}"
      draggable="true"
      @dragstart="handleDragStart"
      @dragend="handleDragEnd"
      @click="toggleOpen"
      @dblclick="startEditing"
    >
      <!-- Expand/Collapse Icon -->
      <svg
        class="w-4 h-4 text-neutral-400 transition-transform duration-200 flex-shrink-0"
        :class="{ 'rotate-90': isOpen }"
        fill="none" stroke="currentColor" viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
      
      <!-- Folder Icon -->
      <svg class="w-4 h-4 text-neutral-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      
      <!-- Title / Input -->
      <div class="flex-1 overflow-hidden" @click.stop="isEditing ? null : toggleOpen()">
        <input
          v-if="isEditing"
          ref="editInput"
          v-model="editName"
          @blur="saveEdit"
          @keyup.enter="saveEdit"
          @keyup.esc="cancelEdit"
          class="w-full bg-neutral-900 text-sm text-neutral-100 px-1 py-0.5 rounded border border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 font-mono"
          type="text"
        />
        <span v-else class="text-sm font-semibold text-neutral-300 truncate block font-mono">
          {{ folder.folder_name || folder.FolderName }}
        </span>
      </div>
      
      <!-- Actions -->
      <div class="flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          @click.stop="startEditing"
          class="p-1 text-neutral-400 hover:text-blue-400 rounded"
          title="Renomear"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>
        </button>
        <button
          @click.stop="$emit('delete', folder.id || folder.ID)"
          class="p-1 text-neutral-400 hover:text-red-400 rounded"
          title="Excluir Pasta"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
        </button>
      </div>
    </div>
    
    <!-- Prompts List -->
    <div
      v-show="isOpen"
      class="pl-6 pr-1 py-1 space-y-1"
    >
      <div v-if="prompts.length === 0" class="text-xs text-neutral-500 py-1 italic opacity-70">
        Pasta vazia...
      </div>
      <button
        v-for="prompt in prompts"
        :key="prompt.name"
        draggable="true"
        @dragstart="handlePromptDragStart($event, prompt)"
        @dragend="handlePromptDragEnd"
        @click="$emit('select-prompt', prompt.name)"
        :class="[
          'w-full text-left px-2 py-1.5 rounded-md text-xs transition font-mono flex items-center gap-2 border',
          getPromptClass(prompt.name)
        ]"
      >
        <svg class="w-3.5 h-3.5 opacity-70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <span class="truncate flex-1">{{ prompt.name }}</span>
        
        <!-- Presence -->
        <span class="flex items-center gap-0.5 flex-shrink-0">
          <span
            v-for="viewer in (allPresence[prompt.name] || [])"
            :key="viewer.username"
            :title="viewer.status === 'editing' ? `${viewer.username} está editando` : `${viewer.username} está visualizando`"
            :class="[
              'w-2 h-2 rounded-full border',
              viewer.status === 'editing' ? 'bg-amber-400 border-amber-600' : 'bg-blue-400 border-blue-600'
            ]"
          ></span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'

const props = defineProps({
  folder: {
    type: Object,
    required: true
  },
  prompts: {
    type: Array,
    default: () => []
  },
  selectedPromptName: {
    type: String,
    default: null
  },
  allPresence: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits([
  'select-prompt',
  'update',
  'delete',
  'drop-into-folder',
  'drag-folder-start',
  'drag-folder-end',
  'drag-prompt-start',
  'drag-prompt-end'
])

const isOpen = ref(false)
const isEditing = ref(false)
const editName = ref('')
const editInput = ref(null)

const isDragOver = ref(false)
const isDragging = ref(false)

function toggleOpen() {
  isOpen.value = !isOpen.value
}

async function startEditing() {
  isEditing.value = true
  editName.value = props.folder.folder_name || props.folder.FolderName
  await nextTick()
  editInput.value?.focus()
}

function saveEdit() {
  if (isEditing.value) {
    const newName = editName.value.trim()
    const folderName = props.folder.folder_name || props.folder.FolderName
    const folderId = props.folder.id || props.folder.ID
    if (newName && newName !== folderName) {
      emit('update', { id: folderId, patch: { folder_name: newName } })
    }
    isEditing.value = false
  }
}

function cancelEdit() {
  isEditing.value = false
}

function getPromptClass(promptName) {
  if (props.selectedPromptName === promptName) {
    return 'bg-blue-600 text-white border-blue-500'
  }
  return 'text-neutral-300 hover:bg-neutral-700/60 border-transparent'
}

// -- Drag & Drop Handlers --

let dragEnterCount = 0;

function handleDragEnter(e) {
  dragEnterCount++
  isDragOver.value = true
}

function handleDragLeave(e) {
  dragEnterCount--
  if (dragEnterCount === 0) {
    isDragOver.value = false
  }
}

function handleDrop(e) {
  dragEnterCount = 0
  isDragOver.value = false
  const type = e.dataTransfer.getData('type') // "folder" or "prompt"
  
  if (type === 'prompt') {
    const promptName = e.dataTransfer.getData('promptName')
    const folderId = props.folder.id || props.folder.ID
    if (promptName) {
      emit('drop-into-folder', {
        folderId: folderId,
        promptName
      })
    }
  } else if (type === 'folder') {
    // If we're dropping a folder on another folder, reorder them
    const draggedFolderId = e.dataTransfer.getData('folderId')
    const currentFolderId = props.folder.id || props.folder.ID
    if (draggedFolderId && parseInt(draggedFolderId) !== currentFolderId) {
       emit('drop-into-folder', {
         targetFolderId: currentFolderId,
         draggedFolderId: parseInt(draggedFolderId)
       })
    }
  }
}

function handleDragStart(e) {
  const folderId = props.folder.id || props.folder.ID
  isDragging.value = true
  e.dataTransfer.setData('type', 'folder')
  e.dataTransfer.setData('folderId', folderId)
  emit('drag-folder-start', folderId)
}

function handleDragEnd(e) {
  isDragging.value = false
  emit('drag-folder-end')
}

// Handle dragging a prompt from within this folder
function handlePromptDragStart(e, prompt) {
  e.dataTransfer.setData('type', 'prompt')
  e.dataTransfer.setData('promptName', prompt.name)
  emit('drag-prompt-start', prompt.name)
}

function handlePromptDragEnd(e) {
  emit('drag-prompt-end')
}
</script>
