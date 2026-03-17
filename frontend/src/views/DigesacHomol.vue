<template>
  <div class="h-screen flex flex-col bg-neutral-900 text-neutral-100 overflow-hidden relative">
    
    <!-- Info Modal Overlay -->
    <div v-if="modalInfo" class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
      <div class="bg-neutral-800 rounded-xl p-6 max-w-md w-full border border-neutral-700 shadow-2xl transform transition-all">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-yellow-500/20 flex items-center justify-center text-yellow-500">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-xl font-bold text-white">{{ modalInfo.title }}</h3>
        </div>
        <p class="text-neutral-300 mb-6 leading-relaxed">{{ modalInfo.message }}</p>
        <div class="flex justify-end">
          <button @click="modalInfo = null" class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition font-medium">
            Entendido
          </button>
        </div>
      </div>
    </div>

    <!-- Create Prompt Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
      <div class="bg-neutral-800 rounded-xl p-6 max-w-lg w-full border border-neutral-700 shadow-2xl">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-full bg-green-500/20 flex items-center justify-center text-green-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <h3 class="text-xl font-bold text-white">Novo Prompt</h3>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-neutral-300 mb-1.5">Nome do arquivo</label>
            <input
              v-model="newPromptName"
              type="text"
              placeholder="ex: meu_prompt.md"
              class="w-full bg-neutral-900 text-neutral-100 px-3 py-2 rounded-lg border border-neutral-700 focus:outline-none focus:border-green-500 font-mono text-sm transition"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-neutral-300 mb-1.5">Conteúdo</label>
            <textarea
              v-model="newPromptContent"
              rows="8"
              placeholder="Digite o conteúdo markdown do prompt..."
              class="w-full bg-neutral-900 text-neutral-100 px-3 py-2 rounded-lg border border-neutral-700 focus:outline-none focus:border-green-500 font-mono text-sm resize-none transition"
            ></textarea>
          </div>
          <p v-if="createError" class="text-sm text-red-400">{{ createError }}</p>
        </div>
        <div class="flex justify-end gap-3 mt-6">
          <button @click="closeCreateModal" class="px-4 py-2 bg-neutral-700 hover:bg-neutral-600 text-white rounded-lg text-sm font-medium transition">
            Cancelar
          </button>
          <button @click="handleCreatePrompt" :disabled="creating" class="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition flex items-center gap-2">
            <svg v-if="creating" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ creating ? 'Criando...' : 'Criar Prompt' }}</span>
          </button>
        </div>
      </div>
    </div>
    <!-- Header -->
    <header class="flex-shrink-0 p-4 bg-neutral-800 border-b border-neutral-700 flex items-center justify-between">
      <div>
        <h1 class="text-xl font-bold tracking-wider text-blue-500">DIGESAC HOMOL</h1>
        <p class="text-xs text-neutral-400">Gerenciador de Prompts</p>
      </div>
      <div class="flex items-center gap-3">
        <button @click="handleSyncAll" :disabled="syncingAll" class="px-4 py-2 bg-yellow-600 hover:bg-yellow-700 disabled:opacity-50 rounded-lg text-sm transition font-medium flex items-center gap-2">
          <svg v-if="syncingAll" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span v-if="syncingAll">Sincronizando...</span>
          <span v-else>Sincronizar Tudo</span>
        </button>
        <button @click="handleLogout" class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-sm transition font-medium">
          Sair do Sistema
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Sidebar / List -->
      <aside class="w-72 bg-neutral-800 border-r border-neutral-700 flex flex-col overflow-hidden">
        <div class="p-4 border-b border-neutral-700 flex-shrink-0 flex items-center justify-between">
          <h2 class="font-semibold text-neutral-200">Arquivos .md</h2>
          <button
            @click="showCreateModal = true"
            title="Criar novo prompt"
            class="w-7 h-7 flex items-center justify-center rounded-lg bg-green-600 hover:bg-green-500 text-white transition"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
        
        <div class="flex-1 overflow-y-auto p-2 space-y-1">
          <div v-if="loadingList" class="p-4 text-center text-neutral-500">
            Carregando lista...
          </div>
          <div v-else-if="listError" class="p-4 text-sm text-red-400 bg-red-900/20 rounded">
            {{ listError }}
          </div>
          <button
            v-for="prompt in prompts"
            :key="prompt.name"
            @click="selectPrompt(prompt.name)"
            :class="[
              'w-full text-left px-3 py-2.5 rounded-lg text-sm transition font-mono flex items-center gap-2 border',
              getPromptClass(prompt.name)
            ]"
          >
            <svg v-if="prompt.name === 'base_prompt.md' || prompt.name === 'router.md'" class="w-4 h-4 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
            <svg v-else class="w-4 h-4 opacity-70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="truncate">{{ prompt.name }}</span>
          </button>
        </div>
      </aside>

      <!-- Viewer / Content -->
      <main class="flex-1 flex flex-col overflow-hidden bg-neutral-900">
        <div v-if="!selectedPromptName" class="flex-1 flex flex-col items-center justify-center text-neutral-500">
          <svg class="w-16 h-16 mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          <p class="text-lg">Selecione um prompt na lateral para visualizar</p>
        </div>
        
        <div v-else-if="loadingDetail" class="flex-1 flex items-center justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>

        <div v-else-if="detailError" class="p-6">
          <div class="p-4 text-sm text-red-400 bg-red-900/20 border border-red-500 rounded-lg">
            {{ detailError }}
          </div>
        </div>

        <div v-else-if="selectedPromptData" class="flex-1 flex flex-col overflow-hidden">
          <div class="p-4 border-b border-neutral-800 bg-neutral-800/50 flex-shrink-0 flex items-center justify-between">
            <h3 class="font-mono text-lg font-bold text-white">{{ selectedPromptData.name }}</h3>
            
            <!-- Edit Controls -->
            <div v-if="!isEditing" class="flex gap-2">
              <button @click="handleSyncSingle" :disabled="syncingSingle" class="px-4 py-1.5 bg-neutral-700 hover:bg-neutral-600 disabled:opacity-50 rounded-lg text-sm text-white font-medium transition cursor-pointer flex items-center gap-2">
                <svg v-if="syncingSingle" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ syncingSingle ? 'Resetando...' : 'Resetar p/ Original' }}</span>
              </button>
              <button @click="startEditing" class="px-4 py-1.5 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm text-white font-medium transition cursor-pointer">
                Editar
              </button>
              <button @click="handleDeletePrompt" :disabled="deleting" class="px-3 py-1.5 bg-red-700 hover:bg-red-600 disabled:opacity-50 rounded-lg text-sm text-white font-medium transition cursor-pointer flex items-center gap-1.5">
                <svg v-if="deleting" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                <span>{{ deleting ? 'Apagando...' : 'Apagar' }}</span>
              </button>
            </div>
            <div v-else class="flex gap-2">
              <button @click="cancelEditing" class="px-4 py-1.5 bg-neutral-600 hover:bg-neutral-500 rounded-lg text-sm text-white font-medium transition cursor-pointer">
                Cancelar
              </button>
              <button @click="savePrompt" :disabled="saving" class="px-4 py-1.5 bg-green-600 hover:bg-green-700 disabled:opacity-50 rounded-lg text-sm text-white font-medium transition flex items-center gap-2 cursor-pointer">
                <svg v-if="saving" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span v-if="saving">Salvando...</span>
                <span v-else>Salvar Alterações</span>
              </button>
            </div>
          </div>
          
          <div class="flex-1 overflow-hidden flex flex-col p-6 md:p-8">
            <!-- Edit Mode Textarea -->
            <textarea 
              v-if="isEditing" 
              v-model="editContent" 
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
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { logout } from '@/api/auth'
import { fetchPromptsList, fetchPromptDetails, updatePrompt, syncAllPrompts, syncSinglePrompt, createPrompt, deletePrompt } from '@/api/digesac'
import { marked } from 'marked'

const router = useRouter()

// Standard auth logic
async function handleLogout() {
  try {
    await logout()
    localStorage.removeItem("user")
  } catch (err) {
    console.error(err)
  } finally {
    router.push('/login')
  }
}

// Viewer State
const prompts = ref([])
const loadingList = ref(false)
const listError = ref('')

const selectedPromptName = ref(null)
const selectedPromptData = ref(null)
const loadingDetail = ref(false)
const detailError = ref('')

// Editing State
const isEditing = ref(false)
const editContent = ref('')
const saving = ref(false)
const syncingAll = ref(false)
const syncingSingle = ref(false)
const deleting = ref(false)

// Modal State
const modalInfo = ref(null)

// Create Prompt Modal State
const showCreateModal = ref(false)
const newPromptName = ref('')
const newPromptContent = ref('')
const creating = ref(false)
const createError = ref('')

function getPromptClass(promptName) {
  if (selectedPromptName.value === promptName) {
    return 'bg-blue-600 text-white border-transparent'
  }
  if (promptName === 'base_prompt.md' || promptName === 'router.md') {
    return 'bg-yellow-900/20 text-yellow-500 hover:bg-yellow-900/40 border-yellow-700/50'
  }
  return 'text-neutral-300 hover:bg-neutral-700 border-transparent'
}

// Initialize Markdown parser security configs natively
marked.setOptions({
  gfm: true,
  breaks: true,
})

// Computed for rendering markdown safely and reactively
const renderedMarkdown = computed(() => {
  if (!selectedPromptData.value || !selectedPromptData.value.prompt) return ''
  return marked.parse(selectedPromptData.value.prompt)
})

onMounted(async () => {
  await loadPrompts()
})

async function loadPrompts() {
  loadingList.value = true
  listError.value = ''
  try {
    const data = await fetchPromptsList()
    
    // Sort logic with special rules
    prompts.value = data.sort((a, b) => {
      // base_prompt.md goes to absolute top
      if (a.name === 'base_prompt.md') return -1;
      if (b.name === 'base_prompt.md') return 1;
      // router.md goes directly below base_prompt
      if (a.name === 'router.md') return -1;
      if (b.name === 'router.md') return 1;
      
      // otherwise, alphabetical
      return a.name.localeCompare(b.name)
    })
  } catch (error) {
    console.error(error)
    listError.value = 'Falha ao carregar a lista de prompts'
  } finally {
    loadingList.value = false
  }
}

async function selectPrompt(name) {
  // Prevent changing if currently editing to avoid losing work
  if (isEditing.value) {
    if(!confirm("Você tem modificações não salvas. Deseja sair sem salvar?")) {
      return
    }
  }

  // Trigger modal explanations for special files
  if (name === 'base_prompt.md') {
    modalInfo.value = {
      title: "Prompts do Agente",
      message: "O base_prompt.md sempre vai servir como a \"personalidade\" do Agente, tudo o que estiver escrito aqui será incluído em qualquer outro prompt selecionado pelo router.md"
    }
  } else if (name === 'router.md') {
    modalInfo.value = {
      title: "Roteamento",
      message: "O router.md serve para decidir o prompt que será concatenado ao base_prompt.md então aqui ele recebe apenas regras de QUAL prompt ele DEVE redirecionar (respostas geradas aqui NÃO são incluídas na resposta final)"
    }
  }

  isEditing.value = false
  editContent.value = ''
  selectedPromptName.value = name
  loadingDetail.value = true
  detailError.value = ''
  selectedPromptData.value = null

  try {
    const data = await fetchPromptDetails(name)
    selectedPromptData.value = data
  } catch (error) {
    console.error(error)
    detailError.value = `Falha ao carregar o conteúdo de ${name}`
  } finally {
    loadingDetail.value = false
  }
}

function startEditing() {
  editContent.value = selectedPromptData.value?.prompt || '';
  isEditing.value = true;
}

function cancelEditing() {
  isEditing.value = false;
  editContent.value = '';
}

async function savePrompt() {
  saving.value = true;
  detailError.value = '';
  
  try {
    const data = await updatePrompt(selectedPromptName.value, editContent.value);
    
    // Update local state with the saved data
    if(selectedPromptData.value) {
      selectedPromptData.value.prompt = data.prompt || editContent.value;
    }
    
    // Exit edit mode smoothly
    isEditing.value = false;
  } catch (error) {
    console.error(error);
    detailError.value = "Falha ao salvar as modificações. Tente novamente.";
  } finally {
    saving.value = false;
  }
}

async function handleSyncAll() {
  if (!confirm("Isso irá resetar TODOS os prompts para a versão de produção. Tem certeza?")) {
    return
  }

  syncingAll.value = true
  try {
    await syncAllPrompts()
    await loadPrompts() // Resync list
    if (selectedPromptName.value) {
      await selectPrompt(selectedPromptName.value) // Resync current detail
    }
    alert("Sincronização completa!")
  } catch (error) {
    console.error(error)
    alert("Erro ao sincronizar prompts.")
  } finally {
    syncingAll.value = false
  }
}

async function handleSyncSingle() {
  if (!confirm(`Deseja resetar "${selectedPromptName.value}" para a versão original de produção?`)) {
    return
  }

  syncingSingle.value = true
  try {
    await syncSinglePrompt(selectedPromptName.value)
    await selectPrompt(selectedPromptName.value) // Resync details
    alert("Prompt resetado com sucesso!")
  } catch (error) {
    console.error(error)
    alert("Erro ao resetar prompt.")
  } finally {
    syncingSingle.value = false
  }
}

function closeCreateModal() {
  showCreateModal.value = false
  newPromptName.value = ''
  newPromptContent.value = ''
  createError.value = ''
}

async function handleCreatePrompt() {
  const trimmedName = newPromptName.value.trim()
  if (!trimmedName) {
    createError.value = 'O nome do arquivo é obrigatório.'
    return
  }

  creating.value = true
  createError.value = ''
  try {
    await createPrompt(trimmedName, newPromptContent.value)
    closeCreateModal()
    await loadPrompts()
    await selectPrompt(trimmedName)
  } catch (error) {
    console.error(error)
    createError.value = 'Falha ao criar o prompt. Tente novamente.'
  } finally {
    creating.value = false
  }
}

async function handleDeletePrompt() {
  if (!confirm(`Tem certeza que deseja apagar "${selectedPromptName.value}" permanentemente?`)) {
    return
  }

  deleting.value = true
  try {
    await deletePrompt(selectedPromptName.value)
    selectedPromptName.value = null
    selectedPromptData.value = null
    await loadPrompts()
  } catch (error) {
    console.error(error)
    alert('Erro ao apagar o prompt.')
  } finally {
    deleting.value = false
  }
}
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
