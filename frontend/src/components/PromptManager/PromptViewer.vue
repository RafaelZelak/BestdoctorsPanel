<template>
  <main class="flex-1 flex flex-col overflow-hidden bg-neutral-900">
    <div v-if="!selectedPromptName" class="flex-1 flex flex-col items-center justify-center text-neutral-500">
      <svg class="w-16 h-16 mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
      <p class="text-lg">Selecione um prompt na lateral para visualizar</p>
    </div>
    
    <div v-else-if="loading" class="flex-1 flex items-center justify-center">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
    </div>

    <div v-else-if="error" class="p-6">
      <div class="p-4 text-sm text-red-400 bg-red-900/20 border border-red-500 rounded-lg">
        {{ error }}
      </div>
    </div>

    <div v-else-if="selectedPromptData" class="flex-1 flex flex-col overflow-hidden">
      <div class="p-4 border-b border-neutral-800 bg-neutral-800/50 flex-shrink-0 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <h3 class="font-mono text-lg font-bold text-white">{{ selectedPromptData.name }}</h3>
          <div class="flex items-center gap-2">
            <span v-if="selectedPromptData.version" class="text-xs bg-neutral-700 text-neutral-300 px-2 py-0.5 rounded border border-neutral-600 font-mono">
              v{{ selectedPromptData.version }}
            </span>
            <span v-if="selectedPromptData.edited_by" class="text-xs text-neutral-500 flex items-center gap-1">
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path></svg>
              {{ selectedPromptData.edited_by }}
            </span>
          </div>
        </div>
        
        <div v-if="!isEditing" class="flex gap-1">
          <!-- History Toggle -->
          <button @click="$emit('toggle-history')" 
            :class="['p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 group h-10', 
                    showHistoryPanel ? 'text-indigo-500 bg-indigo-500/10' : 'text-neutral-400 hover:text-indigo-500 hover:bg-indigo-500/10']"
            title="Histórico">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
              Histórico
            </span>
          </button>
          
          <template v-if="!showHistoryPanel">
            <button v-if="showSync" @click="$emit('sync-single')" :disabled="syncingSingle" 
              class="p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 text-neutral-400 hover:text-yellow-500 hover:bg-yellow-500/10 group h-10"
              title="Resetar p/ Original">
              <svg v-if="syncingSingle" class="animate-spin h-5 w-5 text-yellow-500" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
                {{ syncingSingle ? 'Resetando...' : 'Resetar p/ Original' }}
              </span>
            </button>

            <button @click="$emit('start-editing')" 
              class="p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 text-neutral-400 hover:text-blue-500 hover:bg-blue-500/10 group h-10"
              title="Editar">
              <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
              <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
                Editar
              </span>
            </button>

            <button @click="$emit('delete-prompt')" :disabled="deleting" 
              class="p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 text-neutral-400 hover:text-red-500 hover:bg-red-500/10 group h-10"
              title="Apagar">
              <svg v-if="deleting" class="animate-spin h-5 w-5 text-red-500" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
                {{ deleting ? 'Apagando...' : 'Apagar' }}
              </span>
            </button>
          </template>
        </div>
        <div v-else class="flex gap-2">
          <button @click="$emit('cancel-editing')" 
            class="px-4 py-1.5 transition-all duration-300 rounded-lg text-sm text-neutral-400 hover:text-white hover:bg-neutral-800 font-medium">
            Cancelar
          </button>
          <button @click="$emit('save-prompt')" :disabled="saving" 
            class="px-4 py-1.5 transition-all duration-300 rounded-lg text-sm font-medium flex items-center gap-2 text-green-500 hover:bg-green-500/10 disabled:opacity-50">
            <svg v-if="saving" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span v-if="saving">Salvando...</span>
            <span v-else>Salvar Alterações</span>
          </button>
        </div>
      </div>
      
      <div v-if="showHistoryPanel" class="flex-1 overflow-auto flex flex-col p-6 bg-neutral-900 border-t border-neutral-800 relative">
        <div v-if="loadingDiff" class="absolute inset-0 flex items-center justify-center bg-neutral-900/60 z-10 backdrop-blur-sm">
          <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-if="diffError" class="p-4 mb-4 text-sm text-red-400 bg-red-900/20 border border-red-500/50 rounded-lg">
          {{ diffError }}
        </div>
        
        <div v-else-if="!selectedVersionDiff" class="flex-1 flex flex-col items-center justify-center text-neutral-500 text-lg space-y-4">
          <svg class="w-16 h-16 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg>
          <p>Selecione uma versão na lateral para comparar o "Antes e Depois"</p>
        </div>
        
        <div v-else class="flex flex-col h-full gap-4 max-w-5xl mx-auto w-full">
          <!-- Action Revert -->
          <div class="flex justify-end shrink-0">
            <button 
              @click="$emit('revert-to-version')" 
              :disabled="reverting"
              class="bg-green-600 hover:bg-green-500 text-white px-5 py-2.5 rounded-lg font-bold shadow-lg disabled:opacity-50 flex items-center gap-2 cursor-pointer transition"
            >
              <svg v-if="reverting" class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" /></svg>
              <span>{{ reverting ? 'Restaurando...' : 'Restaurar esta versão' }}</span>
            </button>
          </div>
          
          <div class="flex-1 overflow-auto font-mono text-sm leading-relaxed p-6 bg-[#0d1117] text-[#c9d1d9] rounded-xl border border-neutral-700 shadow-inner whitespace-pre-wrap">
            <div class="break-all sm:break-normal">
              <template v-if="selectedVersionDiff.diff_from_previous && selectedVersionDiff.diff_from_previous.length > 0">
                <template v-for="(segment, i) in selectedVersionDiff.diff_from_previous" :key="i">
                  <span v-if="segment.op === 'equal' || !segment.op" class="text-neutral-400">{{ segment.text || segment.content }}</span>
                  <span v-else-if="segment.op === 'insert' || segment.op === 'added'" class="bg-[#2ea043]/30 text-[#fff] border-b border-[#2ea043] rounded-sm">{{ segment.text || segment.content }}</span>
                  <span v-else-if="segment.op === 'delete' || segment.op === 'removed'" class="bg-[#da3633]/30 text-red-300 line-through rounded-sm opacity-80">{{ segment.text || segment.content }}</span>
                </template>
              </template>
              <template v-else-if="selectedVersionDiff.snapshot">
                <span class="text-neutral-400">{{ selectedVersionDiff.snapshot }}</span>
              </template>
              <template v-else>
                <span class="text-neutral-500 italic">Nenhum snapshot ou diff disponível</span>
              </template>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="flex-1 overflow-hidden flex flex-col p-6 md:p-8">
        <!-- Edit Mode Textarea -->
        <textarea 
          v-if="isEditing" 
          :value="modelValue"
          @input="$emit('update:modelValue', $event.target.value)"
          class="flex-1 w-full bg-neutral-900 text-neutral-100 p-4 rounded-lg border border-neutral-700 focus:outline-none focus:border-blue-500 font-mono text-sm resize-none"
          placeholder="Digite o Markdown aqui..."
        ></textarea>

        <!-- Markdown Content Render -->
        <div 
          v-else
          class="flex-1 overflow-y-auto prose prose-invert prose-blue max-w-none 
                 prose-pre:bg-neutral-800 prose-pre:border prose-pre:border-neutral-700
                 prose-a:text-blue-400 hover:prose-a:text-blue-300"
          v-html="renderedMarkdown"
        ></div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { computed } from 'vue'
import { marked } from 'marked'

// Initialize Markdown parser security configs natively
marked.setOptions({
  gfm: true,
  breaks: true,
})

const props = defineProps({
  selectedPromptName: String,
  selectedPromptData: Object,
  loading: Boolean,
  error: String,
  isEditing: Boolean,
  modelValue: String,
  syncingSingle: Boolean,
  deleting: Boolean,
  saving: Boolean,
  
  showHistoryPanel: Boolean,
  selectedVersionDiff: Array,
  loadingDiff: Boolean,
  diffError: String,
  reverting: Boolean,
  showSync: {
    type: Boolean,
    default: true
  }
})

defineEmits([
  'sync-single',
  'start-editing',
  'cancel-editing',
  'delete-prompt',
  'save-prompt',
  'update:modelValue',
  'toggle-history',
  'revert-to-version'
])

// Computed for rendering markdown safely and reactively
const renderedMarkdown = computed(() => {
  if (!props.selectedPromptData || !props.selectedPromptData.prompt) return ''
  return marked.parse(props.selectedPromptData.prompt)
})
</script>

<style>
/* Scoped overrides to make typography gorgeous out-of-the-box natively alongside tailwind-prose standard constraints */
.prose h1, .prose h2, .prose h3 {
  color: #fff !important;
}
.prose code {
  color: #60A5FA !important;
  background-color: #1F2937;
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
  font-weight: 500;
}
.prose pre code {
  background-color: transparent;
  padding: 0;
  color: #D1D5DB !important;
}
</style>
