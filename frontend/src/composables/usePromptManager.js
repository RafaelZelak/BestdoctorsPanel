import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { logout } from '@/api/auth'
import { createPromptApi } from '@/api/prompts'
import { usePresence } from '@/composables/usePresence'

export function usePromptManager(systemId, apiBase) {
  const router = useRouter()
  const { allPresence, publishPresence } = usePresence(systemId)
  const api = createPromptApi(apiBase)

  // Viewer State
  const prompts = ref([])
  const loadingList = ref(false)
  const listError = ref('')

  const selectedPromptName = ref(null)
  const selectedPromptData = ref(null)
  const loadingDetail = ref(false)
  const detailError = ref('')

  // Folder State
  const folders = ref([])
  const loadingFolders = ref(false)
  const folderError = ref('')

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
  const creating = ref(false)
  const createError = ref('')

  // History State
  const showHistoryPanel = ref(false)
  const historyList = ref([])
  const loadingHistory = ref(false)
  const historyError = ref('')

  const selectedHistoryVersion = ref(null)
  const selectedVersionDiff = ref(null)
  const loadingDiff = ref(false)
  const diffError = ref('')
  const reverting = ref(false)

  onMounted(async () => {
    await Promise.all([loadFolders(), loadPrompts()])
  })

  // Auth logic
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

  // Load folders
  async function loadFolders() {
    loadingFolders.value = true
    folderError.value = ''
    try {
      const data = await api.fetchFolders()
      folders.value = data.sort((a, b) => {
        const idxA = a.folder_idx !== undefined ? a.folder_idx : a.FolderIdx
        const idxB = b.folder_idx !== undefined ? b.folder_idx : b.FolderIdx
        if (idxA === null || idxA === undefined) return 1;
        if (idxB === null || idxB === undefined) return -1;
        return idxA - idxB;
      });
    } catch (error) {
      console.error(error)
      folderError.value = 'Falha ao carregar as pastas'
    } finally {
      loadingFolders.value = false
    }
  }

  // Folder actions
  async function handleCreateFolder(name, idx = null) {
    try {
      await api.createFolder(name, idx)
      await loadFolders()
    } catch (error) {
      console.error(error)
      alert("Erro ao criar pasta.")
    }
  }

  async function handleUpdateFolder(id, patch) {
    try {
      await api.updateFolder(id, patch)
      await loadFolders()
    } catch (error) {
      console.error(error)
      alert("Erro ao atualizar pasta.")
    }
  }

  async function handleDeleteFolder(id) {
    if (!confirm("Tem certeza de que deseja excluir esta pasta? Os prompts não serão excluídos.")) {
      return
    }
    try {
      await api.deleteFolder(id)
      await Promise.all([loadFolders(), loadPrompts()])
    } catch (error) {
      console.error(error)
      alert("Erro ao excluir pasta.")
    }
  }

  async function handleAssignPromptToFolder(promptName, folderId) {
    try {
      await api.assignPromptToFolder(promptName, folderId)
      await loadPrompts()
    } catch (error) {
      console.error(error)
      alert("Erro ao mover prompt.")
    }
  }

  // Load list
  async function loadPrompts() {
    loadingList.value = true
    listError.value = ''
    try {
      const data = await api.fetchPromptsList()

      prompts.value = data.sort((a, b) => {
        if (a.name === 'base_prompt.md') return -1;
        if (b.name === 'base_prompt.md') return 1;
        if (a.name === 'router.md') return -1;
        if (b.name === 'router.md') return 1;
        return a.name.localeCompare(b.name)
      })
    } catch (error) {
      console.error(error)
      listError.value = 'Falha ao carregar a lista de prompts'
    } finally {
      loadingList.value = false
    }
  }

  // Details
  async function selectPrompt(name) {
    if (isEditing.value) {
      if (!confirm("Você tem modificações não salvas. Deseja sair sem salvar?")) {
        return
      }
    }

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
    showHistoryPanel.value = false
    selectedHistoryVersion.value = null
    selectedVersionDiff.value = null

    publishPresence(name, 'viewing')

    try {
      const data = await api.fetchPromptDetails(name)
      selectedPromptData.value = data
    } catch (error) {
      console.error(error)
      detailError.value = `Falha ao carregar o conteúdo de ${name}`
    } finally {
      loadingDetail.value = false
    }
  }

  // Editing logic
  function startEditing() {
    editContent.value = selectedPromptData.value?.prompt || '';
    isEditing.value = true;
    publishPresence(selectedPromptName.value, 'editing')
  }

  function cancelEditing() {
    isEditing.value = false;
    editContent.value = '';
    publishPresence(selectedPromptName.value, 'viewing')
  }

  async function savePrompt() {
    saving.value = true;
    detailError.value = '';

    try {
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const fullName = userData.full_name || userData.username || 'System';
      const data = await api.updatePrompt(selectedPromptName.value, editContent.value, fullName);
      if (selectedPromptData.value) {
        selectedPromptData.value.prompt = data.prompt || editContent.value;
      }
      isEditing.value = false;
      publishPresence(selectedPromptName.value, 'viewing')
    } catch (error) {
      console.error(error);
      detailError.value = "Falha ao salvar as modificações. Tente novamente.";
    } finally {
      saving.value = false;
    }
  }

  // Deleting and syncing
  async function handleSyncAll() {
    if (!confirm("Isso irá resetar TODOS os prompts para a versão original. Tem certeza?")) {
      return
    }

    syncingAll.value = true
    try {
      await api.syncAllPrompts()
      await loadPrompts()
      if (selectedPromptName.value) {
        await selectPrompt(selectedPromptName.value)
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
    if (!confirm(`Deseja resetar "${selectedPromptName.value}" para a versão original?`)) {
      return
    }

    syncingSingle.value = true
    try {
      await api.syncSinglePrompt(selectedPromptName.value)
      await selectPrompt(selectedPromptName.value)
      alert("Prompt resetado com sucesso!")
    } catch (error) {
      console.error(error)
      alert("Erro ao resetar prompt.")
    } finally {
      syncingSingle.value = false
    }
  }

  // Create modal logic
  function closeCreateModal() {
    showCreateModal.value = false
    createError.value = ''
  }

  async function handleCreatePrompt({ name, content }) {
    const trimmedName = name.trim()
    if (!trimmedName) {
      createError.value = 'O nome do arquivo é obrigatório.'
      return
    }

    creating.value = true
    createError.value = ''
    try {
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const fullName = userData.full_name || userData.username || 'System';
      await api.createPrompt(trimmedName, content, fullName)
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
      await api.deletePrompt(selectedPromptName.value)
      selectedPromptName.value = null
      selectedPromptData.value = null
      showHistoryPanel.value = false
      await loadPrompts()
    } catch (error) {
      console.error(error)
      alert('Erro ao apagar o prompt.')
    } finally {
      deleting.value = false
    }
  }

  // History Logic
  function toggleHistory() {
    showHistoryPanel.value = !showHistoryPanel.value
    if (showHistoryPanel.value && selectedPromptName.value) {
      loadHistory(selectedPromptName.value)
    } else {
      selectedHistoryVersion.value = null
      selectedVersionDiff.value = null
    }
  }

  function clearHistoryVersion() {
    selectedHistoryVersion.value = null
    selectedVersionDiff.value = null
  }

  async function loadHistory(name) {
    loadingHistory.value = true
    historyError.value = ''
    try {
      const data = await api.fetchPromptHistory(name)
      historyList.value = data || []
    } catch (error) {
      console.error(error)
      historyError.value = 'Falha ao carregar o histórico'
    } finally {
      loadingHistory.value = false
    }
  }

  async function viewVersionDiff(versionItem) {
    selectedHistoryVersion.value = versionItem
    loadingDiff.value = true
    diffError.value = ''
    selectedVersionDiff.value = null
    try {
      const data = await api.fetchPromptDiff(selectedPromptName.value, versionItem.version)
      selectedVersionDiff.value = data
    } catch (error) {
      console.error(error)
      diffError.value = 'Falha ao carregar as mudanças'
    } finally {
      loadingDiff.value = false
    }
  }

  async function revertToVersion(versionItem) {
    if (!confirm(`Deseja mesmo reverter este prompt para a versão ${versionItem.version}?`)) {
      return
    }

    // We use the 'snapshot' from the loaded diff details, or from the version item if available
    const promptContent = selectedVersionDiff.value?.snapshot || versionItem.snapshot

    if (!promptContent) {
      alert("Não foi possível carregar o conteúdo dessa versão para reverter (certifique-se de selecionar a versão primeiro).")
      return
    }

    reverting.value = true
    try {
      const userData = JSON.parse(localStorage.getItem('user') || '{}');
      const fullName = userData.full_name || userData.username || 'System';
      const data = await api.updatePrompt(selectedPromptName.value, promptContent, fullName)
      if (selectedPromptData.value) {
        selectedPromptData.value.prompt = data.prompt || promptContent
      }

      showHistoryPanel.value = false
      selectedHistoryVersion.value = null
      selectedVersionDiff.value = null

      alert("Prompt revertido com sucesso!")
    } catch (error) {
      console.error(error)
      alert("Erro ao reverter o prompt.")
    } finally {
      reverting.value = false
    }
  }

  return {
    prompts,
    loadingList,
    listError,
    folders,
    loadingFolders,
    folderError,
    selectedPromptName,
    selectedPromptData,
    loadingDetail,
    detailError,
    isEditing,
    editContent,
    saving,
    syncingAll,
    syncingSingle,
    deleting,
    modalInfo,
    showCreateModal,
    creating,
    createError,

    showHistoryPanel,
    historyList,
    loadingHistory,
    historyError,
    selectedHistoryVersion,
    selectedVersionDiff,
    loadingDiff,
    diffError,
    reverting,

    allPresence,
    handleLogout,
    selectPrompt,
    startEditing,
    cancelEditing,
    savePrompt,
    handleSyncAll,
    handleSyncSingle,
    closeCreateModal,
    handleCreatePrompt,
    handleDeletePrompt,
    loadPrompts,
    loadFolders,
    handleCreateFolder,
    handleUpdateFolder,
    handleDeleteFolder,
    handleAssignPromptToFolder,
    toggleHistory,
    clearHistoryVersion,
    loadHistory,
    viewVersionDiff,
    revertToVersion
  }
}
